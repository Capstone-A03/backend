package validator

import "github.com/go-playground/validator/v10"

func RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error {
	return validate.RegisterValidation(tag, fn, callValidationEvenIfNull...)
}
