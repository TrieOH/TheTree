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

type ErrPasswordTooLong struct{}

func (ErrPasswordTooLong) Error() string {
	return "password length exceeds 72 bytes"
}
