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

type ErrObjectMismatch struct {
	Expected string
	Actual   string
}

func (e ErrObjectMismatch) Error() string {
	return fmt.Sprintf("object mismatch: permission=%s, request=%s", e.Expected, e.Actual)
}

type ErrActionMismatch struct {
	Expected string
	Actual   string
}

func (e ErrActionMismatch) Error() string {
	return fmt.Sprintf("action mismatch: permission=%s, request=%s", e.Expected, e.Actual)
}

type ErrInsufficientPermissions struct{}

func (e ErrInsufficientPermissions) Error() string {
	return "insufficient permissions"
}
