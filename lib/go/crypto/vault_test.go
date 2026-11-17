package crypto

import (
	"strings"
	"testing"
)

func TestNewVault_RejectsBadKeys(t *testing.T) {
	for name, key := range map[string]string{
		"empty":   "",
		"not hex": strings.Repeat("zz", 32),
		"short":   strings.Repeat("ab", 16),
	} {
		_, err := NewVault(key)
		if err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
}

func TestVault_EncryptDecryptRoundTrip(t *testing.T) {
	v, err := NewVault(strings.Repeat("ab", 32))
	if err != nil {
		t.Fatalf("NewVault: %v", err)
	}
	sealed, err := v.EncryptPrivateKey([]byte("secret"))
	if err != nil {
		t.Fatalf("EncryptPrivateKey: %v", err)
	}
	plain, err := v.DecryptPrivateKey(sealed)
	if err != nil {
		t.Fatalf("DecryptPrivateKey: %v", err)
	}
	if string(plain) != "secret" {
		t.Fatalf("got %q, want %q", plain, "secret")
	}

	other, _ := NewVault(strings.Repeat("cd", 32))
	_, err = other.DecryptPrivateKey(sealed)
	if err == nil {
		t.Fatal("a different master key must not open the sealed value")
	}
}
