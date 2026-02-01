package permissions

import (
	"GoAuth/internal/apierr"
	"regexp"
	"time"

	"github.com/MintzyG/fail"
	"github.com/google/uuid"
)

type Permission struct {
	ID        uuid.UUID
	ProjectID *uuid.UUID
	Object    string
	Action    string
	CreatedAt time.Time

	Conditions *Condition
}

var (
	objectRegex = regexp.MustCompile(`^(?:\*|[a-zA-Z][a-zA-Z0-9_]*:(?:[a-zA-Z0-9_]+|\*)(?:/[a-zA-Z][a-zA-Z0-9_]*:(?:[a-zA-Z0-9_]+|\*))*(?:/(?:\*|\*\*))?)$`)
	actionRegex = regexp.MustCompile(`^(?:\*|[a-zA-Z0-9_]+(?::(?:[a-zA-Z0-9_]+|\*))*(?::\*\*)?)$`)
)

func ValidateObject(object string) error {
	if object == "" {
		return fail.New(apierr.PERMissionInvalidObject).WithArgs("(empty)")
	}
	if !objectRegex.MatchString(object) {
		return fail.New(apierr.PERMissionInvalidObject).WithArgs(object)
	}
	return nil
}

func ValidateAction(action string) error {
	if action == "" {
		return fail.New(apierr.PERMissionInvalidAction).WithArgs("(empty)")
	}
	if !actionRegex.MatchString(action) {
		return fail.New(apierr.PERMissionInvalidAction).WithArgs(action)
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
