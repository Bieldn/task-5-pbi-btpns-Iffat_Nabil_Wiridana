package helpers

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) []string {
	var ve validator.ValidationErrors
	var errorMessage []string

	if errors.As(err, &ve) {
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage = append(errorMessage, e.Error())
		}
	} else {
		errorMessage = append(errorMessage, err.Error())
	}

	return errorMessage
}
