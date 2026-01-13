package validation

import (
	"GoAuth/internal/apierr"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

// PasswordValidationResult holds detailed password validation info
type PasswordValidationResult struct {
	HasUpper  bool
	HasNumber bool
	HasSymbol bool
	IsValid   bool
}

// Custom password validation - requires uppercase, number, and symbol
func passwordValidator(fl validator.FieldLevel) bool {
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
}

// Helper function to check if a field is a password field
func isPasswordField(fieldName, structFieldName string) bool {
	fieldLower := strings.ToLower(fieldName)
	structFieldLower := strings.ToLower(structFieldName)
	passwordWords := []string{"password", "passwd", "pwd", "pass"}

	for _, word := range passwordWords {
		if fieldLower == word || structFieldLower == word ||
			strings.Contains(fieldLower, word) || strings.Contains(structFieldLower, word) {
			return true
		}
	}
	return false
}

func init() {
	// Register custom validators
	_ = validate.RegisterValidation("passwd", passwordValidator)

	_ = validate.RegisterValidation("uuid7", func(fl validator.FieldLevel) bool {
		v := fl.Field().String()

		u, err := uuid.Parse(v)
		if err != nil {
			return false
		}

		return u.Version() == 7
	})

	// Use JSON field names for better API responses
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return fld.Name
		}
		return name
	})
}

// ValidateRule validates a single value and returns a user-friendly error
func ValidateRule(value interface{}, rule, fieldName string) error {
	err := validate.Var(value, rule)
	if err == nil {
		return nil
	}

	verr := apierr.ErrInvalidInput.WithMsg("Validation failed")
	for _, e := range err.(validator.ValidationErrors) {
		msg := formatValidationMessage(e, fieldName, value)
		verr = verr.WithCause(errors.New(msg))
	}

	return verr.WithID(apierr.RequestValidationError)
}

func formatValue(v interface{}) string {
	if v == nil {
		return "null"
	}

	switch val := v.(type) {
	case string:
		if len(val) > 32 {
			return fmt.Sprintf("%q…", val[:32])
		}
		return fmt.Sprintf("%q", val)

	case fmt.Stringer:
		return fmt.Sprintf("%q", val.String())

	default:
		return fmt.Sprintf("%v", val)
	}
}

// formatValidationMessage creates user-friendly error messages that include the field name
func formatValidationMessage(err validator.FieldError, fieldName string, value interface{}) string {
	field := fieldName
	val := formatValue(value)

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address (got %s)", field, val)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID (got %s)", field, val)
	case "uuid7":
		return fmt.Sprintf("%s must be a valid UUIDv7 (got %s)", field, val)
	case "min":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at least %s characters long", field, err.Param())
		}
		return fmt.Sprintf("%s must be at least %s", field, err.Param())
	case "max":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at most %s characters long", field, err.Param())
		}
		return fmt.Sprintf("%s must be at most %s", field, err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s (got %s)", field, strings.ReplaceAll(err.Param(), " ", ", "), val)
	default:
		return fmt.Sprintf("%s failed %s validation (got %s)", field, err.Tag(), val)
	}
}

// ValidateStruct validates any struct and returns a Response with validation errors
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	verr := apierr.ErrInvalidInput.WithMsg("Validation failed")
	for _, err := range err.(validator.ValidationErrors) {
		structFieldName := err.StructField()
		fieldName := structFieldName
		var msg string
		var value interface{}

		// Get JSON field name if available
		structType := reflect.TypeOf(s)
		if structType.Kind() == reflect.Ptr {
			structType = structType.Elem()
		}
		if structField, ok := structType.FieldByName(structFieldName); ok {
			if jsonTag := structField.Tag.Get("json"); jsonTag != "" {
				jsonFieldName := strings.SplitN(jsonTag, ",", 2)[0]
				if jsonFieldName != "-" && jsonFieldName != "" {
					fieldName = jsonFieldName
				}
			}
		}

		isPassword := isPasswordField(fieldName, structFieldName)

		// Don't expose password values
		if err.Value() != nil && !isPassword {
			value = err.Value()
		}

		// Handle password validation specially
		if err.Tag() == "passwd" {
			password := err.Value().(string)
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

			if !hasUpper {
				verr = verr.WithCause(errors.New(fieldName + " must contain at least one uppercase letter"))
			}
			if !hasNumber {
				verr = verr.WithCause(errors.New(fieldName + " must contain at least one number"))
			}
			if !hasSymbol {
				verr = verr.WithCause(errors.New(fieldName + " must contain at least one symbol or punctuation"))
			}
			continue
		}

		msg = formatValidationMessage(err, fieldName, value)
		verr = verr.WithCause(errors.New(msg))
	}

	return verr.WithID(apierr.RequestValidationError)
}

// ValidateInto validates JSON into a provided pointer
func ValidateInto[T any](r *http.Request, target *T) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return apierr.ErrBadRequest.WithMsg("Content-Type must be application/json").WithID(apierr.RequestNotApplicationJSON)
	}

	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return apierr.ErrBadRequest.WithMsg("Invalid JSON format").WithID(apierr.RequestInvalidJSONFormat).WithCause(err)
	}

	if err := ValidateStruct(*target); err != nil {
		return err
	}

	return nil
}
