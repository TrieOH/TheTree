package validator

import (
	"lib/errx"
	"reflect"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func SetupValidator() *validator.Validate {
	var v = validator.New()
	if err := v.RegisterValidation("uuid7", func(fl validator.FieldLevel) bool {
		vv := fl.Field().String()

		u, err := uuid.Parse(vv)
		if err != nil {
			return false
		}

		return u.Version() == 7
	}); err != nil {
		errx.Exit(err, "failed to register uuid7 validator")
	}

	// Custom password validation - requires uppercase, number, and symbol
	if err := v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		var hasUpper, hasNumber, hasSymbol bool

		for _, c := range password {
			switch {
			case unicode.IsUpper(c):
				hasUpper = true
			case unicode.IsNumber(c):
				hasNumber = true
			case unicode.IsPunct(c) || unicode.IsSymbol(c):
				hasSymbol = true
			}
		}

		return hasUpper && hasNumber && hasSymbol
	}); err != nil {
		errx.Exit(err, "failed to register passwd validator")
	}

	// Use JSON field names for better API responses
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return fld.Name
		}
		return name
	})

	return v
}
