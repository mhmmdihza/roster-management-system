package util

import "github.com/go-playground/validator/v10"

func NewValidator() *validator.Validate {
	return validator.New()
}
