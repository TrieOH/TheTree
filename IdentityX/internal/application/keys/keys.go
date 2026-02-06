package keys

import (
	"GoAuth/internal/crypto"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"

	"github.com/google/uuid"
)

// FIXME make these accept context and record to span

type cachedPublic struct {
	pub       ed25519.PublicKey
	projectID *uuid.UUID
}

type UseCase struct {
	repo            outbounds.KeysRepository
	privateKeyCache outbounds.CacheService
	publicKeyCache  outbounds.CacheService
}

var _ inbounds.KeysService = (*UseCase)(nil)

func New(repo outbounds.KeysRepository, privateCache outbounds.CacheService, publicCache outbounds.CacheService) *UseCase {
	return &UseCase{
		repo:            repo,
		privateKeyCache: privateCache,
		publicKeyCache:  publicCache,
	}
}

func (uc *UseCase) SignGoAuth(ctx context.Context, payload []byte) ([]byte, error) {
	pair, err := uc.repo.GetActiveGoAuthSigningKey(ctx)
	if err != nil {
		return nil, err
	}

	val, ok := uc.privateKeyCache.Get(ctx, pair.KID)
	var priv ed25519.PrivateKey
	if ok {
		priv = val.(ed25519.PrivateKey)
	} else {
		decrypted, err := crypto.Decrypt(pair.PrivateKey)
		if err != nil {
			return nil, err
		}

		priv, err = parseEd25519Private(decrypted)
		if err != nil {
			return nil, err
		}
		uc.privateKeyCache.Set(ctx, pair.KID, priv, 0)
	}

	sig := ed25519.Sign(priv, payload)
	return sig, nil
}

func (uc *UseCase) SignProject(ctx context.Context, projectID uuid.UUID, payload []byte) ([]byte, error) {
	pair, err := uc.repo.GetActiveProjectSigningKey(ctx, projectID)
	if err != nil {
		return nil, err
	}

	val, ok := uc.privateKeyCache.Get(ctx, pair.KID)
	var priv ed25519.PrivateKey
	if ok {
		priv = val.(ed25519.PrivateKey)
	} else {
		decrypted, err := crypto.Decrypt(pair.PrivateKey)
		if err != nil {
			return nil, err
		}

		priv, err = parseEd25519Private(decrypted)
		if err != nil {
			return nil, err
		}
		uc.privateKeyCache.Set(ctx, pair.KID, priv, 0)
	}

	sig := ed25519.Sign(priv, payload)
	return sig, nil
}

func (uc *UseCase) VerifyGoAuth(ctx context.Context, kid string, payload, sig []byte) error {
	val, ok := uc.publicKeyCache.Get(ctx, kid)
	if ok {
		cached := val.(cachedPublic)
		if cached.projectID != nil {
			return inbounds.ErrKeyProjectMismatch{}
		}
		if !ed25519.Verify(cached.pub, payload, sig) {
			return inbounds.ErrInvalidSignature{}
		}
		return nil
	}

	pair, err := uc.repo.GetGoAuthKeyByKID(ctx, kid)
	if err != nil {
		return err
	}

	pub, err := parseEd25519Public(pair.PublicKey)
	if err != nil {
		return err
	}

	uc.publicKeyCache.Set(ctx, kid, cachedPublic{pub: pub, projectID: nil}, 0)

	if !ed25519.Verify(pub, payload, sig) {
		return inbounds.ErrInvalidSignature{}
	}

	return nil
}

func (uc *UseCase) VerifyProject(ctx context.Context, projectID uuid.UUID, kid string, payload, sig []byte) error {
	val, ok := uc.publicKeyCache.Get(ctx, kid)
	if ok {
		cached := val.(cachedPublic)
		if cached.projectID == nil || *cached.projectID != projectID {
			return inbounds.ErrKeyProjectMismatch{}
		}
		if !ed25519.Verify(cached.pub, payload, sig) {
			return inbounds.ErrInvalidSignature{}
		}
		return nil
	}

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

	uc.publicKeyCache.Set(ctx, kid, cachedPublic{pub: pub, projectID: pair.ProjectID}, 0)

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

func (uc *UseCase) RevokeKey(ctx context.Context, kid string) error {
	if err := uc.repo.RevokeKeyByKID(ctx, kid); err != nil {
		return err
	}

	uc.privateKeyCache.Delete(ctx, kid)
	uc.publicKeyCache.Delete(ctx, kid)

	return nil
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
