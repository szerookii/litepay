package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var validate = validator.New()

type ValidationError struct {
	FailedField string
	Value       interface{}
	Tag         string
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("'%s' ('%v') needs to implement '%s'", v.FailedField, v.Value, v.Tag)
}

func FormatValidationErrors(errs []ValidationError) string {
	var errors []string
	for _, err := range errs {
		errors = append(errors, err.Error())
	}

	return strings.Join(errors, ", ")
}

func Validate(data any) []ValidationError {
	errs := validate.Struct(data)
	if errs == nil {
		return nil
	}

	var validationErrors []ValidationError
	for _, err := range errs.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, ValidationError{
			FailedField: err.Field(),
			Value:       err.Value(),
			Tag:         err.Tag(),
		})
	}

	return validationErrors
}
