package goauth

import (
	"regexp"

	"github.com/MintzyG/fail/v3"
)

var objectRegex = regexp.MustCompile(`^(\*|[a-zA-Z][a-zA-Z0-9_]*)$`)

func validateObject(obj string) error {
	if !objectRegex.MatchString(obj) {
		return fail.New(SDKInvalidObjectFormatID).WithArgs(obj)
	}
	return nil
}
