package permissions

import (
	"GoAuth/internal/errx"
	"context"
	"regexp"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
)

type Permission struct {
	ID        uuid.UUID
	ProjectID *uuid.UUID
	Object    string
	Action    string
	CreatedAt time.Time
}

var (
	// Plain object name only - no wildcards, no separators
	objectRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	// Plain action name only - no wildcards, no separators
	actionRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
)

func ValidateObject(ctx context.Context, object string) error {
	if object == "" {
		return fail.New(errx.PERMissionInvalidObject).WithArgs("(empty)").RecordCtx(ctx)
	}
	if !objectRegex.MatchString(object) {
		return fail.New(errx.PERMissionInvalidObject).WithArgs(object).RecordCtx(ctx)
	}
	return nil
}

func ValidateAction(ctx context.Context, action string) error {
	if action == "" {
		return fail.New(errx.PERMissionInvalidAction).WithArgs("(empty)").RecordCtx(ctx)
	}
	if !actionRegex.MatchString(action) {
		return fail.New(errx.PERMissionInvalidAction).WithArgs(action).RecordCtx(ctx)
	}
	return nil
}

func ValidatePermission(ctx context.Context, object, action string) error {
	if err := ValidateObject(ctx, object); err != nil {
		return err
	}
	if err := ValidateAction(ctx, action); err != nil {
		return err
	}
	return nil
}
