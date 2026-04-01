package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	fun "github.com/MintzyG/FastUtilitiesNet/response"
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

	var errs []error
	var validationErrors validator.ValidationErrors
	ok := errors.As(err, &validationErrors)
	if !ok {
		return fun.NewErrorf("error validating request: %s", err.Error()).BadRequest()
	}
	for _, e := range validationErrors {
		ve := formatValidationMessage(e, fieldName, value)
		errs = append(errs, ve)
	}
	return fun.NewError("validation errors").WithErrors(errs...).Validation()
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
func formatValidationMessage(err validator.FieldError, fieldName string, value any) error {
	field := fieldName
	val := formatValue(value)

	var message string
	switch err.Tag() {
	case "required":
		message = fmt.Sprintf("%s is required", field)
	case "email":
		message = fmt.Sprintf("%s must be a valid email address (got %s)", field, val)
	case "uuid":
		message = fmt.Sprintf("%s must be a valid UUID (got %s)", field, val)
	case "uuid7":
		message = fmt.Sprintf("%s must be a valid UUIDv7 (got %s)", field, val)
	case "min":
		if err.Kind() == reflect.String {
			message = fmt.Sprintf("%s must be at least %s characters long", field, err.Param())
		} else {
			message = fmt.Sprintf("%s must be at least %s", field, err.Param())
		}
	case "max":
		if err.Kind() == reflect.String {
			message = fmt.Sprintf("%s must be at most %s characters long", field, err.Param())
		} else {
			message = fmt.Sprintf("%s must be at most %s", field, err.Param())
		}
	case "oneof":
		message = fmt.Sprintf("%s must be one of: %s (got %s)", field, strings.ReplaceAll(err.Param(), " ", ", "), val)
	default:
		message = fmt.Sprintf("%s failed %s validation (got %s)", field, err.Tag(), val)
	}

	return &fun.FieldError{
		Field:   fieldName,
		Message: message,
		Value:   value,
	}
}

// ValidateStruct validates any struct and returns a Response with validation errors
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errs []error
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return fun.NewErrorf("error validating request: %s", err.Error()).BadRequest()
	}
	for _, err := range validationErrors {
		structFieldName := err.StructField()
		fieldName := structFieldName
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

		errs = append(errs, fun.FieldError{
			Field:   fieldName,
			Message: err.Error(),
			Value:   value,
		})
	}
	return fun.NewError("validation errors").WithErrors(errs...).Validation()
}

// ValidateInto validates JSON into a provided pointer
func ValidateInto[T any](r *http.Request, target *T) error {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return fun.NewError("request content type is not application/json").BadRequest()
	}

	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return fun.NewErrorf("request body is invalid json: %s", err.Error()).BadRequest()
	}

	if err := ValidateStruct(*target); err != nil {
		return err
	}

	return nil
}

type Validator func() error

func Run(validators ...Validator) error {
	var errs []error

	for _, v := range validators {
		if v == nil {
			continue
		}
		if err := v(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return fun.NewError("validation errors").WithErrors(errs...).Validation()
}

func RequireUUID(domain, field string, value uuid.UUID) Validator {
	return func() error {
		if value == uuid.Nil {
			return fun.FieldError{Field: domain + "." + field, Message: "missing required field"}
		}
		return nil
	}
}

func RequireString(domain, field, value string) Validator {
	return func() error {
		if value == "" {
			return fun.FieldError{Field: domain + "." + field, Message: "missing required field"}
		}
		return nil
	}
}

func RequireTime(domain, field string, value time.Time) Validator {
	return func() error {
		if value.IsZero() {
			return fun.FieldError{Field: domain + "." + field, Message: "missing required field"}
		}
		return nil
	}
}

func Assert(domain string, condition bool, message string) Validator {
	return func() error {
		if !condition {
			return fun.FieldError{Field: domain, Message: message}
		}
		return nil
	}
}

func AssertIf(domain string, guard func() bool, condition func() bool, message string) Validator {
	return func() error {
		if !guard() {
			return nil
		}
		if !condition() {
			return fun.FieldError{Field: domain, Message: message}
		}
		return nil
	}
}
