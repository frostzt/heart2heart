package utils

import (
	"github.com/go-playground/validator/v10"
)

func GenerateValidationErrors(validationErrors error) []ErrorMessage {
	errorsEncountered := []ErrorMessage{}

	for _, e := range validationErrors.(validator.ValidationErrors) {
		errorMessage := ErrorMessage{
			Message: e.Error(),
		}

		errorsEncountered = append(errorsEncountered, errorMessage)
	}

	return errorsEncountered
}

func GenerateGenericErrors(validationErrors error) []ErrorMessage {
	errorsEncountered := []ErrorMessage{}

	for _, e := range validationErrors.(validator.ValidationErrors) {
		errorMessage := ErrorMessage{
			Message: e.Error(),
		}

		errorsEncountered = append(errorsEncountered, errorMessage)
	}

	return errorsEncountered
}
