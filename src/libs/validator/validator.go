package validator

import "github.com/go-playground/validator/v10"

type FieldLevel = validator.FieldLevel
type Func = func(fl FieldLevel) bool

var validate = validator.New()
