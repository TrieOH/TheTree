package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"univents/internal/apierr"

	"github.com/MintzyG/fail/v3"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// FIXME make these errors either go into spans or logs

var validate = validator.New()

func init() {
	// Register custom validators
	if err := validate.RegisterValidation("uuid7", func(fl validator.FieldLevel) bool {
		v := fl.Field().String()

		u, err := uuid.Parse(v)
		if err != nil {
			return false
		}

		return u.Version() == 7
	}); err != nil {
		panic("failed to register uuid7 validator: " + err.Error())
	}

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

	verr := fail.New(apierr.RequestValidationError)
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return fail.New(apierr.RequestValidationError).With(err)
	}
	for _, e := range validationErrors {
		msg := formatValidationMessage(e, fieldName, value)
		verr = verr.Trace(msg)
	}
	return verr.AddMeta("module", "validation")
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

	verr := fail.New(apierr.RequestValidationError)
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return fail.New(apierr.RequestValidationError).With(err)
	}
	for _, err := range validationErrors {
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

		msg = formatValidationMessage(err, fieldName, value)
		verr = verr.Trace(msg)
	}
	return verr
}

// ValidateInto validates JSON into a provided pointer
func ValidateInto[T any](r *http.Request, target *T) error {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return fail.New(apierr.RequestNotApplicationJSON)
	}

	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return fail.New(apierr.RequestInvalidJSONFormat).With(err)
	}

	if err := ValidateStruct(*target); err != nil {
		fe, ok := fail.As(err)
		if ok {
			return fe.AddMeta("module", "validation")
		}
		return err
	}

	return nil
}
