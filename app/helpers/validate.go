package helpers

import (
	"fmt"

	"github.com/go-playground/validator"
)

type ValidationError struct {
	ActualTag string `json:"tag"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Param     string `json:"param"`
}

func WrapValidationErrors(errs validator.ValidationErrors) []ValidationError {
	validationErrors := make([]ValidationError, 0, len(errs))
	for _, validationErr := range errs {
		validationErrors = append(validationErrors, ValidationError{
			ActualTag: validationErr.ActualTag(),
			Namespace: validationErr.Namespace(),
			Kind:      validationErr.Kind().String(),
			Type:      validationErr.Type().String(),
			Value:     fmt.Sprintf("%v", validationErr.Value()),
			Param:     validationErr.Param(),
		})
	}
	return validationErrors
}
