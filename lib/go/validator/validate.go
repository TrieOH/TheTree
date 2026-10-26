package validator

import (
	"errors"
	"reflect"
	"strings"
	"unicode"

	"lib/fielderr"

	"github.com/go-playground/validator/v10"
)

// Error is a request that failed validation. It reads as every failed
// field, so clients that only show the message still learn which field
// failed and why, and lists the fields for the response (fielderr.Lister).
// lib/httpserver maps it to a 400 VALIDATION_ERROR for every Backend.
type Error struct{ Fields []fielderr.Field }

func (e *Error) Error() string { return validationMessage(e.Fields) }

// FieldList implements fielderr.Lister.
func (e *Error) FieldList() []fielderr.Field { return e.Fields }

// Validate runs the shared validator against v — typically a strict-server
// request body (or the whole RequestObject, since the validator dives into
// nested structs). It returns an *Error listing every failed field, so
// handlers can return it directly.
//
//	body, _ := ...
//	if err := validator.Validate(request.Body); err != nil {
//	    return nil, err
//	}
func Validate(v any) error {
	vd := getValidator()
	if vd == nil {
		panic("validator: no validator registered — call SetupValidator at startup")
	}
	err := vd.Struct(v)
	if err != nil {
		return &Error{Fields: validationErrsToFields(err)}
	}
	return nil
}

// validationMessage builds the top-level envelope message from the
// per-field errors, so clients that only surface `message` still learn
// which field failed and why.
func validationMessage(fields []fielderr.Field) string {
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		parts = append(parts, f.Field+": "+f.Message)
	}
	if len(parts) == 0 {
		return "invalid request body"
	}
	return "invalid request body: " + strings.Join(parts, "; ")
}

// validationErrsToFields converts validator.ValidationErrors into fields,
// with password masking and passwd expansion.
func validationErrsToFields(err error) []fielderr.Field {
	ve, ok := errors.AsType[validator.ValidationErrors](err)
	if !ok {
		return []fielderr.Field{{Field: "body", Message: err.Error()}}
	}

	out := make([]fielderr.Field, 0, len(ve))
	for _, fe := range ve {
		if fe.Tag() == "passwd" {
			out = append(out, passwdFieldErrors(fe)...)
			continue
		}

		var value any
		if !isPasswordField(fe.Field()) {
			value = fe.Value()
		}

		out = append(out, fielderr.Field{
			Field:   fe.Field(),
			Message: tagMessage(fe),
			Value:   value,
		})
	}
	return out
}

// isPasswordField reports whether a field name looks like a password field.
func isPasswordField(name string) bool {
	lower := strings.ToLower(name)
	for _, word := range []string{"password", "passwd", "pwd", "pass"} {
		if lower == word || strings.Contains(lower, word) {
			return true
		}
	}
	return false
}

// passwdFieldErrors expands a failed passwd tag into one field per
// missing requirement, without echoing the password value.
func passwdFieldErrors(fe validator.FieldError) []fielderr.Field {
	password, ok := fe.Value().(string)
	if !ok {
		return []fielderr.Field{{Field: fe.Field(), Message: "must be a valid string"}}
	}

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

	var out []fielderr.Field
	if len(password) < 8 {
		out = append(out, fielderr.Field{Field: fe.Field(), Message: "must be at least 8 characters long"})
	}
	if !hasUpper {
		out = append(out, fielderr.Field{Field: fe.Field(), Message: "must contain at least one uppercase letter"})
	}
	if !hasNumber {
		out = append(out, fielderr.Field{Field: fe.Field(), Message: "must contain at least one number"})
	}
	if !hasSymbol {
		out = append(out, fielderr.Field{Field: fe.Field(), Message: "must contain at least one symbol or punctuation"})
	}
	return out
}

// staticTagMessages holds the fixed human-readable messages for validator
// tags that don't depend on the field kind or bound.
var staticTagMessages = map[string]string{
	"required": "this field is required",
	"email":    "must be a valid email address",
	"url":      "must be a valid URL",
	"uuid":     "must be a valid UUID",
	"uuid4":    "must be a valid UUIDv4",
	"uuid7":    "must be a valid UUIDv7",
	"numeric":  "must be a numeric value",
	"alpha":    "must contain only letters",
	"alphanum": "must contain only letters and numbers",
	"oneof":    "must be one of the allowed values",
}

// tagMessage produces a human-readable message for a validation failure.
// Kind-aware for len/min/max/gt/gte/lt/lte so strings say "characters" and
// numbers say the actual bound.
func tagMessage(fe validator.FieldError) string {
	if msg, ok := staticTagMessages[fe.Tag()]; ok {
		if fe.Tag() == "oneof" {
			return "must be one of: " + strings.ReplaceAll(fe.Param(), " ", ", ")
		}
		return msg
	}

	param := fe.Param()
	isStr := fe.Kind() == reflect.String

	switch fe.Tag() {
	case "len":
		if isStr {
			return "must be exactly " + param + " characters long"
		}
		return "must have exactly " + param + " elements"
	case "min":
		if isStr {
			return "must be at least " + param + " characters long"
		}
		return "must be at least " + param
	case "max":
		if isStr {
			return "must be at most " + param + " characters long"
		}
		return "must be at most " + param
	case "gt":
		if isStr {
			return "must be longer than " + param + " characters"
		}
		return "must be greater than " + param
	case "gte":
		if isStr {
			return "must be at least " + param + " characters long"
		}
		return "must be greater than or equal to " + param
	case "lt":
		if isStr {
			return "must be shorter than " + param + " characters"
		}
		return "must be less than " + param
	case "lte":
		if isStr {
			return "must be at most " + param + " characters long"
		}
		return "must be less than or equal to " + param
	}

	return "failed validation: " + fe.Tag()
}
