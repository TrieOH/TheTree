package field

import "fmt"

type ErrFieldNotDefined struct {
	Key string
}

func (e ErrFieldNotDefined) Error() string {
	return "field not defined in schema: " + e.Key
}

type ErrInvalidFieldType struct {
	Key      string
	Expected string
	Got      any
}

func (e ErrInvalidFieldType) Error() string {
	return fmt.Sprintf("field %q expects type %s, got %T", e.Key, e.Expected, e.Got)
}

type ErrMissingRequiredFields struct {
	Key string
}

func (e ErrMissingRequiredFields) Error() string {
	return "Missing required field: " + e.Key
}

type ErrFieldsValidation struct {
	FieldErrors []error
}

func (ErrFieldsValidation) Error() string {
	return "error validating fields"
}
