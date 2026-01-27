package roles

type ErrEmptyRoleName struct{}

func (e ErrEmptyRoleName) Error() string {
	return "role name is empty"
}
