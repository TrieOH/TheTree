package secrets_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"lib/secrets"
)

type static map[string]string

func (s static) Load(context.Context) (map[string]string, error) { return s, nil }

func TestApplyKeepsExplicitEnv(t *testing.T) {
	t.Setenv("SECRETS_TEST_SET", "explicit")
	t.Setenv("SECRETS_TEST_UNSET", "")
	_ = os.Unsetenv("SECRETS_TEST_UNSET")

	applied, err := secrets.Apply(context.Background(), static{
		"SECRETS_TEST_SET":   "from-secret",
		"SECRETS_TEST_UNSET": "from-secret",
	})
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("SECRETS_TEST_SET") != "explicit" {
		t.Fatalf("explicit env was overwritten")
	}
	if os.Getenv("SECRETS_TEST_UNSET") != "from-secret" {
		t.Fatalf("missing env was not filled")
	}
	if len(applied) != 1 || applied[0] != "SECRETS_TEST_UNSET" {
		t.Fatalf("applied = %v", applied)
	}
}

func TestParseDotenv(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	err := os.WriteFile(path, []byte("# comment\n\nA=1\nexport B='two'\nC=\"x=y\"\nnot a pair\nEMPTY=\n"), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	got, err := secrets.ParseDotenv(path)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"A": "1", "B": "two", "C": "x=y", "EMPTY": ""}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("%s = %q, want %q", k, got[k], v)
		}
	}
}
