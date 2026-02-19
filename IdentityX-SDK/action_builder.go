package goauth

import (
	"regexp"

	"github.com/MintzyG/fail/v3"
)

var actionRegex = regexp.MustCompile(`^(\*|[a-zA-Z][a-zA-Z0-9_]*)$`)

func validateAction(act string) error {
	if !actionRegex.MatchString(act) {
		return fail.New(SDKInvalidActionFormatID).WithArgs(act)
	}
	return nil
}
