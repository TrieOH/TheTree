package utils

import (
	"errors"
	"strings"
)

var (
	ErrDuplicate  = errors.New("duplicate value")
	ErrForeignKey = errors.New("foreign key constraint failed")
	ErrNotFound   = errors.New("record not found")
	ErrUnknown    = errors.New("unknown database error")
)

var FieldMessages = map[string]string{
	"users_email_key": "email is already in use",
}

func ParseDBError(err error) error {
	if err == nil {
		return nil
	}

	msg := err.Error()

	switch {
	case strings.Contains(msg, "duplicate key value"):
		field := extractConstraint(msg)
		if custom, ok := FieldMessages[field]; ok {
			return errors.New(custom)
		}
		// fallback generic message
		if fieldName := extractField(field); fieldName != "" {
			return errors.New(fieldName + " already exists")
		}
		return ErrDuplicate

	case strings.Contains(msg, "violates foreign key constraint"):
		field := extractConstraint(msg)
		if custom, ok := FieldMessages[field]; ok {
			return errors.New(custom)
		}
		return ErrForeignKey

	case strings.Contains(msg, "no rows in result set"):
		return ErrNotFound

	default:
		return ErrUnknown
	}
}

// extractConstraint extracts the actual constraint name from the error message
// e.g. `"users_email_key"` → `users_email_key`
func extractConstraint(msg string) string {
	start := strings.Index(msg, "\"")
	end := strings.LastIndex(msg, "\"")
	if start != -1 && end > start {
		return msg[start+1 : end]
	}
	return ""
}

// extractField tries to get a readable field name from a constraint name
// e.g. `users_email_key` → `email`
func extractField(constraint string) string {
	start := strings.Index(constraint, "_") + 1
	end := strings.LastIndex(constraint, "_key")
	if start > 0 && end > start {
		return constraint[start:end]
	}
	return ""
}
