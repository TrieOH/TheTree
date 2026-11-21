package main

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
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

func TestServerSourceEmbedsEveryFeatureWithHandlers(t *testing.T) {
	dir := t.TempDir()
	write := func(rel, body string) {
		t.Helper()
		p := filepath.Join(dir, rel)
		err := os.MkdirAll(filepath.Dir(p), 0o750)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(p, []byte(body), 0o600)
		if err != nil {
			t.Fatal(err)
		}
	}
	write("api_keys/handlers.go", "package api_keys\n\ntype Handlers struct{}\n")
	write("store/raw.go", "package storefront\n\ntype Handlers struct{}\n")
	write("crypto/repo.go", "package crypto\n\ntype Repo struct{}\n")
	write("crypto/handlers_test.go", "package crypto\n\ntype Handlers struct{}\n")

	features, err := findFeatures("Svc", dir)
	if err != nil {
		t.Fatal(err)
	}
	src, err := serverSource("Svc", features)
	if err != nil {
		t.Fatal(err)
	}
	got := string(src)
	for _, want := range []string{
		`"Svc/internal/features/api_keys"`,
		`storefront "Svc/internal/features/store"`,
		"apiKeysHandlers = api_keys.Handlers",
		"storeHandlers   = storefront.Handlers",
		"\t*apiKeysHandlers\n\t*storeHandlers\n}",
		"var _ openapi.StrictServerInterface = (*server)(nil)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("generated source lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "crypto") {
		t.Errorf("a folder without a non-test Handlers type was embedded:\n%s", got)
	}
}
