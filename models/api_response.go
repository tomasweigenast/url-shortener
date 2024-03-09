package models

import (
	"github.com/go-playground/validator/v10"
)

type ApiResponse[T any] struct {
	Success bool     `json:"success"`
	Error   ApiError `json:"error,omitempty"`
	Data    *T       `json:"data,omitempty"`
}

type ApiError interface {
	apiError()
}

type ValidationError struct {
	Violations []Violation `json:"violations"`
}

type Violation struct {
	FieldName  string `json:"field"`
	Constraint string `json:"constraint"`
	Argument   string `json:"argument,omitempty"`
}

func (ValidationError) apiError() {}

func NewValidationError(errors *validator.ValidationErrors) *ValidationError {
	violations := make([]Violation, len(*errors))

	for i, err := range *errors {
		violations[i] = Violation{
			FieldName:  err.Field(),
			Constraint: err.Tag(),
			Argument:   err.Param(),
		}
	}

	return &ValidationError{violations}
}

type StringError struct {
	Reason string `json:"reason"`
}

func (StringError) apiError() {}

func ParseError(err error) ApiError {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		return NewValidationError(&validationErrors)
	}

	return StringError{
		Reason: err.Error(),
	}
}
