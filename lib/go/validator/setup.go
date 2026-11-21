package validator

import (
	"lib/errx"
	"reflect"
	"strings"
	"sync"
	"unicode"
	"uuid"

	"github.com/go-playground/validator/v10"
)

// globalValidator is the shared instance used by Validate. The harness
// builds it at boot (SetupValidator); anything that validates first, like
// a test calling an operation directly, builds it on first use.
var (
	globalValidator *validator.Validate
	setupOnce       sync.Once
)

func getValidator() *validator.Validate {
	setupOnce.Do(func() {
		if globalValidator == nil {
			SetupValidator()
		}
	})
	return globalValidator
}

func SetupValidator() *validator.Validate {
	var v = validator.New()
	err := v.RegisterValidation("uuid7", func(fl validator.FieldLevel) bool {
		vv := fl.Field().String()

		u, err := uuid.Parse(vv)
		if err != nil {
			return false
		}

		// stdlib uuid has no Version(): the version is the high nibble of byte 6 (RFC 9562 §4.2).
		return u[6]>>4 == 7
	})
	if err != nil {
		errx.Exit(err, "failed to register uuid7 validator")
	}

	// Custom password validation - requires uppercase, number, and symbol
	err = v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
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
	})
	if err != nil {
		errx.Exit(err, "failed to register passwd validator")
	}

	// Use JSON field names for better API responses
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		if name == "" {
			return fld.Name
		}
		return name
	})

	globalValidator = v
	return v
}
