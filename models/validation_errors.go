package models

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

type ValidationErrors struct {
	Err validator.ValidationErrors
}

func (err *ValidationErrors) ToHumanReadable(param string) (errors []HumanReadableStatus){
	for _, err := range err.Err {
		errors = append(errors, HumanReadableStatus{
			Type: fmt.Sprintf("%s-is-invalid", strings.ToLower(err.Field())),
			Message: GenerateReason(err.Tag(), strings.ToLower(err.Field()), err.Param()),
			Param: strings.ToLower(err.Field()),
			Value: err.Value(),
			Source: param,
		})
		}
	return
}

func GenerateReason(err_type string, err_field string, err_param string) string{
	switch (err_type) {
	case "min":
		return fmt.Sprintf("The provided %s is invalid because it is less than the minimum length of %s", err_field, err_param)
	case "email":
		return fmt.Sprintf("The provided %s is invalid because it is not a valid email", err_field)
	case "required":
		return fmt.Sprintf("%s is a required field, but no value was provided", err_field)
	default:
		return ""
	}
}