package apierr

func (e Error) WithMsg(msg string) *Error {
	e.Message = msg
	return &e
}

func (e Error) WithCause(err error) *Error {
	e.Cause = err
	return &e
}

func (e Error) WithCode(code Code) *Error {
	e.Code = code
	return &e
}

func (e Error) WithID(id ID) *Error {
	e.ID = id
	return &e
}

func (e Error) WithField(key string, value any) *Error {
	if e.Fields == nil {
		e.Fields = map[string]any{}
	}
	e.Fields[key] = value
	return &e
}
