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
