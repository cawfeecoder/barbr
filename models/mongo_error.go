package models

import (
	"fmt"
	"regexp"
	"strings"
)

type MongoError struct {
	Key string `json:"key"`
	Value string `json:"value"`
}
var bad_value_regex = regexp.MustCompile(`[(]BadValue[)]\s([$]\w+)\s`)
var unknown_op_regex = regexp.MustCompile(`.+:(?P<Operator>.+)`)
var duplicate_regex = regexp.MustCompile(`\D\d+\s+duplicate`)
var extract_dup_values = regexp.MustCompile(`collection: \D+[.](?P<Collection>\D+) index: (?P<Field>\D+)_.+key:.+"(?P<Value>\D+)"`)

func GetErrorFromMongo(err error, param string) []HumanReadableStatus {
	switch {
	case strings.Contains(err.Error(), "(BadValue)") == true:
		match := bad_value_regex.FindStringSubmatch(err.Error())
		fmt.Printf("%v", match)
		return []HumanReadableStatus{HumanReadableStatus{Type:"bad-query-value", Message: "The provided query value is invalid", Param: "comparators", Value: match[1]}}
	case strings.Contains(err.Error(), "unknown operator") == true:
		match := unknown_op_regex.FindStringSubmatch(err.Error())
		return []HumanReadableStatus{HumanReadableStatus{Type:"unknown-query-operator", Message: "The provided operator is unknown or an otherwise bad value", Param: "comparators", Value: match[1]}}
	case err.Error() == "incorrect password provided":
		return []HumanReadableStatus{HumanReadableStatus{Type: "user-incorrect-password", Message: "Incorrect password provided", Param: param}}
	case err.Error() == "mongo: no documents in result":
		return []HumanReadableStatus{HumanReadableStatus{Type: "id-not-found", Message: "ID does not reference any users, either because they don't exist or do not have an active status", Param: "id", Value: param, Source: param}}
	case err.Error() == "the provided hex string is not a valid ObjectID" || err.Error() == "encoding/hex: odd length hex string":
		return []HumanReadableStatus{HumanReadableStatus{Type: "id-not-valid", Message: "Provided ID is invalid because it is not a valid ObjectID", Param: "id", Value: param}}
	case duplicate_regex.MatchString(err.Error()[1:len(err.Error())]):
		split := strings.Split(err.Error()[1:len(err.Error())], ",")
		for _, val := range split {
			var dup_errors []HumanReadableStatus
			match := extract_dup_values.FindStringSubmatch(val)
			dup_errors = append(dup_errors, HumanReadableStatus{Type: fmt.Sprintf("dup-value-%s", match[2]), Message: fmt.Sprintf("Duplicate %s found when creating user", match[2]), Param: match[2], Value: match[3]})
			return dup_errors
		}
	default:
		return []HumanReadableStatus{HumanReadableStatus{Type: "unknown-error", Message: err.Error()}}
	}
	return nil
}
