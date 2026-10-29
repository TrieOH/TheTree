package main

import (
	"go/format"
	"testing"
)

// The generated file is never formatted after it is written, so the source
// must already be gofmt'd for the Backends' lint to pass.
func TestSourceIsGofmted(t *testing.T) {
	formatted, err := format.Source([]byte(source))
	if err != nil {
		t.Fatalf("source does not parse: %v", err)
	}
	if string(formatted) != source {
		t.Fatalf("source is not gofmt'd; gofmt gives:\n%s", formatted)
	}
}
