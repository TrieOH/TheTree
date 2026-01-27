package permissions

import "fmt"

type ErrInvalidPermissionObject struct {
	Object  string
	IsEmpty bool
}

func (e ErrInvalidPermissionObject) Error() string {
	if e.IsEmpty {
		return fmt.Sprintf("invalid permission object: (empty)")
	}
	return fmt.Sprintf("invalid permission object: %s", e.Object)
}

type ErrInvalidPermissionAction struct {
	Action  string
	IsEmpty bool
}

func (e ErrInvalidPermissionAction) Error() string {
	if e.IsEmpty {
		return fmt.Sprintf("invalid permission action: (empty)")
	}
	return fmt.Sprintf("invalid permission action: %s", e.Action)
}
