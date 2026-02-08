package goauth

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/MintzyG/fail/v3"
)

type FinalizedBuilder interface {
	String() string
	Build() (string, error)
}

type ObjectBuilder interface {
	Child(namespace, specifier string) ObjectBuilder
	NamespaceAny(namespace string) ObjectBuilder
	Any() FinalizedBuilder
	All() FinalizedBuilder
	FinalizedBuilder
}

type objectBuilder struct {
	segments []string
}

func Object(namespace, specifier string) ObjectBuilder {
	return &objectBuilder{segments: []string{fmt.Sprintf("%s:%s", namespace, specifier)}}
}

func ObjectWildcard() FinalizedBuilder {
	return &objectBuilder{segments: []string{"*"}}
}

func (b *objectBuilder) Child(namespace, specifier string) ObjectBuilder {
	b.segments = append(b.segments, fmt.Sprintf("%s:%s", namespace, specifier))
	return b
}

func (b *objectBuilder) NamespaceAny(namespace string) ObjectBuilder {
	b.segments = append(b.segments, fmt.Sprintf("%s:*", namespace))
	return b
}

func (b *objectBuilder) Any() FinalizedBuilder {
	b.segments = append(b.segments, "*")
	return b
}

func (b *objectBuilder) All() FinalizedBuilder {
	b.segments = append(b.segments, "**")
	return b
}

func (b *objectBuilder) String() string {
	return strings.Join(b.segments, "/")
}

func (b *objectBuilder) Build() (string, error) {
	s := b.String()
	if err := validateObject(s); err != nil {
		return "", err
	}
	return s, nil
}

var objectRegex = regexp.MustCompile(`^(?:\*|[a-zA-Z][a-zA-Z0-9_]*:(?:[a-zA-Z0-9_]+|\*)(?:/[a-zA-Z][a-zA-Z0-9_]*:(?:[a-zA-Z0-9_]+|\*))*(?:/(?:\*|\*\*))?)$`)

func validateObject(obj string) error {
	if !objectRegex.MatchString(obj) {
		return fail.New(SDKInvalidObjectFormatID).WithArgs(obj)
	}
	return nil
}
