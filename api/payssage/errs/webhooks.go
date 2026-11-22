package errs

// MissingParamError is a required header or query parameter the request
// lacks. In is "header" or "query". It stays by hand: it answers the shared
// validation problem with its own text.
type MissingParamError struct {
	In   string
	Name string
}

func (e *MissingParamError) Error() string { return "missing required " + e.In + ": " + e.Name }
