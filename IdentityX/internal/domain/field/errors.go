package field

type ErrFieldNotDefined struct {
	Key string
}

func (e ErrFieldNotDefined) Error() string {
	return "field not defined in schema: " + e.Key
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

type ErrInvalidFieldValue struct {
	Key   string
	Type  string
	Value any
}

func (e ErrInvalidFieldValue) Error() string {
	return "invalid value for field '" + e.Key + "' of type " + e.Type
}
