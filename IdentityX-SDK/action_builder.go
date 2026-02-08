package goauth

import (
	"regexp"
	"strings"

	"github.com/MintzyG/fail/v3"
)

type ActionBuilder interface {
	Sub(action string) ActionBuilder
	Any() ActionBuilder
	All() FinalizedBuilder
	FinalizedBuilder
}

type actionBuilder struct {
	segments []string
}

func Action(action string) ActionBuilder {
	return &actionBuilder{segments: []string{action}}
}

func ActionWildcard() FinalizedBuilder {
	return &actionBuilder{segments: []string{"*"}}
}

func (b *actionBuilder) Sub(action string) ActionBuilder {
	b.segments = append(b.segments, action)
	return b
}

func (b *actionBuilder) Any() ActionBuilder {
	b.segments = append(b.segments, "*")
	return b
}

func (b *actionBuilder) All() FinalizedBuilder {
	b.segments = append(b.segments, "**")
	return b
}

func (b *actionBuilder) String() string {
	return strings.Join(b.segments, ":")
}

func (b *actionBuilder) Build() (string, error) {
	s := b.String()
	if err := validateAction(s); err != nil {
		return "", err
	}
	return s, nil
}

var actionRegex = regexp.MustCompile(`^(?:\*|[a-zA-Z0-9_]+(?::(?:[a-zA-Z0-9_]+|\*))*(?::\*\*)?)$`)

func validateAction(act string) error {
	if !actionRegex.MatchString(act) {
		return fail.New(SDKInvalidActionFormatID).WithArgs(act)
	}
	return nil
}
