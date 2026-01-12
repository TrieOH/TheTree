package project_users

type ErrEncodingProjectUserMetadata struct {
	Cause error
}

func (ErrEncodingProjectUserMetadata) Error() string {
	return "error encoding project user metadata"
}
