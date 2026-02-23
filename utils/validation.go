package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatError(err error) []string {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errors = append(errors, fmt.Sprintf("%s is required", e.Field()))
			case "email":
				errors = append(errors, fmt.Sprintf("%s must be a valid email", e.Field()))
			case "min":
				errors = append(errors, fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param()))
			case "gt":
				errors = append(errors, fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param()))
			case "unique_email", "unique_username", "unique_member_code", "unique_inventory_code":
				errors = append(errors, fmt.Sprintf("%s must be unique", e.Field()))
			case "oneof":
				errors = append(errors, fmt.Sprintf("%s must be one of the following: %s", e.Field(), e.Param()))
			default:
				errors = append(errors, fmt.Sprintf("%s is invalid (%s)", e.Field(), e.Tag()))
			}
		}
	}
	return errors
}
