package errs

// AccessLevelError is a ticket whose access level is below the program's.
type AccessLevelError struct{ Kind string }

func (e *AccessLevelError) Error() string {
	return "your ticket level does not grant access to this " + e.Kind
}
