package utils

import (
	"fmt"
	"toncap-backend/types"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct[T any](payload T) []*types.ValidateErrorResponse {
	var errors []*types.ValidateErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element types.ValidateErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func ValidateErrorString(errors []*types.ValidateErrorResponse) string {
	errorMessage := ""

	for _, err := range errors {
		errorMessage += fmt.Sprintf("Field: %s, Tag: %s, Value: %s\n", err.Field, err.Tag, err.Value)
	}

	return errorMessage
}
