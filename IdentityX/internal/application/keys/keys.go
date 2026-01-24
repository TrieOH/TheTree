package keys

import (
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"

	"github.com/google/uuid"
)

type UseCase struct {
	repo outbounds.KeysRepository
}

var _ inbounds.KeysService = (*UseCase)(nil)

func New(repo outbounds.KeysRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) SignGoAuth(ctx context.Context, payload []byte) ([]byte, error) {
	pair, err := uc.repo.GetActiveGoAuthSigningKey(ctx)
	if err != nil {
		return nil, err
	}

	priv, err := parseEd25519Private(pair.PrivateKey)
	if err != nil {
		return nil, err
	}
	defer zero(priv)

	sig := ed25519.Sign(priv, payload)
	return sig, nil
}

func (uc *UseCase) SignProject(ctx context.Context, projectID uuid.UUID, payload []byte) ([]byte, error) {
	pair, err := uc.repo.GetActiveProjectSigningKey(ctx, projectID)
	if err != nil {
		return nil, err
	}

	priv, err := parseEd25519Private(pair.PrivateKey)
	if err != nil {
		return nil, err
	}
	defer zero(priv)

	sig := ed25519.Sign(priv, payload)
	return sig, nil
}

func (uc *UseCase) VerifyGoAuth(ctx context.Context, kid string, payload, sig []byte) error {
	pair, err := uc.repo.GetGoAuthKeyByKID(ctx, kid)
	if err != nil {
		return err
	}

	pub, err := parseEd25519Public(pair.PublicKey)
	if err != nil {
		return err
	}

	if !ed25519.Verify(pub, payload, sig) {
		return inbounds.ErrInvalidSignature{}
	}

	return nil
}

func (uc *UseCase) VerifyProject(ctx context.Context, projectID uuid.UUID, kid string, payload, sig []byte) error {
	pair, err := uc.repo.GetProjectKeyByKID(ctx, kid)
	if err != nil {
		return err
	}

	if pair.ProjectID == nil || *pair.ProjectID != projectID {
		return inbounds.ErrKeyProjectMismatch{}
	}

	pub, err := parseEd25519Public(pair.PublicKey)
	if err != nil {
		return err
	}

	if !ed25519.Verify(pub, payload, sig) {
		return inbounds.ErrInvalidSignature{}
	}

	return nil
}

func (uc *UseCase) GetActiveGoAuthSigningKID(ctx context.Context) (string, error) {
	return uc.repo.GetActiveGoAuthSigningKID(ctx)
}

func (uc *UseCase) GetActiveProjectSigningKID(ctx context.Context, projectID uuid.UUID) (string, error) {
	return uc.repo.GetActiveProjectSigningKID(ctx, projectID)
}

func parseEd25519Private(pemBytes []byte) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, inbounds.ErrInvalidPEMPrivKey{}
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, inbounds.ErrParsingPKCS8PrivKey{Cause: err}
	}

	priv, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, inbounds.ErrNotED25519PrivKey{}
	}

	return priv, nil
}

func parseEd25519Public(pemStr string) (ed25519.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, inbounds.ErrInvalidPEMPubKey{}
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, inbounds.ErrParsingPKIXPubKey{Cause: err}
	}

	pub, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, inbounds.ErrNotED25519PubKey{}
	}

	return pub, nil
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
