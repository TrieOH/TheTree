package keys

import (
	"context"
	"strings"
	"testing"
	"uuid"

	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"
)

// sealedKeys answers GetActive with one key sealed by a given vault; every
// other method is left to the embedded nil interface.
type sealedKeys struct {
	ports.CryptoKeysRepo

	key models.CryptoKey
}

func (s sealedKeys) GetActive(context.Context, models.CryptoKeyType, *uuid.UUID) (*models.CryptoKey, error) {
	return &s.key, nil
}

func newVault(t *testing.T, hex string) *crypto.Vault {
	t.Helper()
	v, err := crypto.NewVault(hex)
	if err != nil {
		t.Fatalf("NewVault: %v", err)
	}
	return v
}

func sealedWith(t *testing.T, v *crypto.Vault) sealedKeys {
	t.Helper()
	kp, err := v.GenerateKeyPair(string(models.SigningCryptoKeyType))
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	return sealedKeys{key: models.CryptoKey{EncryptedPrivateKey: kp.EncryptedPrivate}}
}

func TestCheckMasterKeyOpensKeysSealedWithIt(t *testing.T) {
	v := newVault(t, strings.Repeat("ab", 32))
	m := NewManager(sealedWith(t, v), nil, Config{Vault: v})

	err := m.CheckMasterKey(context.Background())
	if err != nil {
		t.Fatalf("want the master key to open its own keys, got %v", err)
	}
}

func TestCheckMasterKeyFailsWithAnotherKey(t *testing.T) {
	sealer := newVault(t, strings.Repeat("ab", 32))
	m := NewManager(sealedWith(t, sealer), nil, Config{Vault: newVault(t, strings.Repeat("cd", 32))})

	err := m.CheckMasterKey(context.Background())
	if err == nil || !strings.Contains(err.Error(), "ENCRYPTION_KEY") {
		t.Fatalf("want an ENCRYPTION_KEY error for a different master key, got %v", err)
	}
}
