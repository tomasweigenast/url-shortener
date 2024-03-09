package models

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var ModelValidator = validator.New(validator.WithRequiredStructEnabled())

func init() {
	ModelValidator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
}
