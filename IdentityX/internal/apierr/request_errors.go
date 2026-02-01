package apierr

type ErrMissingCustomFields struct{}

func (ErrMissingCustomFields) Error() string {
	return "schema custom fields are required on a schema register"
}

type ErrInvalidCustomFieldsJSON struct {
	Cause error
}

func (ErrInvalidCustomFieldsJSON) Error() string {
	return "invalid custom fields JSON"
}

type ErrParsingNumber struct {
	Cause error
}

func (e ErrParsingNumber) Error() string {
	return "error parsing number: " + e.Cause.Error()
}

type ErrMissingParam struct {
	Param string
}

func (e ErrMissingParam) Error() string {
	return "missing parameter: " + e.Param
}
