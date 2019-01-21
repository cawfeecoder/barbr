package models

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"strings"
)

type ValidationErrors struct {
	V_errors validator.ValidationErrors
}

func (err *ValidationErrors) ToHumanReadable() (errors []HumanReadableError){
	for _, err := range err.V_errors {
		errors = append(errors, HumanReadableError{
			Key: strings.ToLower(err.Field()),
			Reason: GenerateReason(err.Tag(), strings.ToLower(err.Field()), err.Param()),
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