package permissions

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID         uuid.UUID
	ProjectID  *uuid.UUID
	Object     string
	Action     string
	Conditions *json.RawMessage
	CreatedAt  time.Time
}

var (
	objectRegex = regexp.MustCompile(`^(?:\*|[a-zA-Z][a-zA-Z0-9_]*:(?:[a-zA-Z0-9_]+|\*)(?:/[a-zA-Z][a-zA-Z0-9_]*:(?:[a-zA-Z0-9_]+|\*))*(?:/(?:\*|\*\*))?)$`)
	actionRegex = regexp.MustCompile(`^(?:\*|[a-zA-Z0-9_]+(?::[a-zA-Z0-9_]+)*)$`)
)

func ValidateObject(object string) error {
	if object == "" {
		return ErrInvalidPermissionObject{IsEmpty: true}
	}
	if !objectRegex.MatchString(object) {
		return ErrInvalidPermissionObject{Object: object}
	}
	return nil
}

func ValidateAction(action string) error {
	if action == "" {
		return ErrInvalidPermissionAction{IsEmpty: true}
	}
	if !actionRegex.MatchString(action) {
		return ErrInvalidPermissionAction{Action: action}
	}
	return nil
}

func ValidatePermission(object, action string) error {
	if err := ValidateObject(object); err != nil {
		return err
	}
	if err := ValidateAction(action); err != nil {
		return err
	}
	return nil
}
