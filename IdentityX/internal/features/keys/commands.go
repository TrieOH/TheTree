package keys

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/shared/crypto"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"context"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// FIXME make these accept context and record to span

type cachedPublic struct {
	pub       ed25519.PublicKey
	projectID *uuid.UUID
}

type CommandService struct {
	keys            ports.KeysRepository
	privateKeyCache ports.CacheService
	publicKeyCache  ports.CacheService
	logger          *zap.Logger
	tracer          trace.Tracer
	txRunner        database.TxRunner
}

func NewCommandService(
	repo ports.KeysRepository,
	privateCache ports.CacheService,
	publicCache ports.CacheService,
	logger *zap.Logger,
	tracer trace.Tracer,
	txRunner database.TxRunner,
) *CommandService {
	return &CommandService{
		keys:            repo,
		privateKeyCache: privateCache,
		publicKeyCache:  publicCache,
		logger:          logger,
		tracer:          tracer,
		txRunner:        txRunner,
	}
}

func (uc *CommandService) SignGoAuth(ctx context.Context, payload []byte) ([]byte, error) {
	pair, err := uc.keys.GetActiveGoAuthSigningKey(ctx)
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

func (uc *CommandService) SignProject(ctx context.Context, projectID uuid.UUID, payload []byte) ([]byte, error) {
	pair, err := uc.keys.GetActiveProjectSigningKey(ctx, projectID)
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

func (uc *CommandService) VerifyGoAuth(ctx context.Context, kid string, payload, sig []byte) error {
	val, ok := uc.publicKeyCache.Get(ctx, kid)
	if ok {
		cached := val.(cachedPublic)
		if cached.projectID != nil {
			return fail.New(errx.KeysProjectKeyMismatch).RecordCtx(ctx)
		}
		if !ed25519.Verify(cached.pub, payload, sig) {
			return fail.New(errx.KeysInvalidSignature).RecordCtx(ctx)
		}
		return nil
	}

	pair, err := uc.keys.GetGoAuthKeyByKID(ctx, kid)
	if err != nil {
		return err
	}

	pub, err := parseEd25519Public(pair.PublicKey)
	if err != nil {
		return err
	}

	uc.publicKeyCache.Set(ctx, kid, cachedPublic{pub: pub, projectID: nil}, 0)

	if !ed25519.Verify(pub, payload, sig) {
		return fail.New(errx.KeysInvalidSignature).RecordCtx(ctx)
	}

	return nil
}

func (uc *CommandService) VerifyProject(ctx context.Context, projectID uuid.UUID, kid string, payload, sig []byte) error {
	val, ok := uc.publicKeyCache.Get(ctx, kid)
	if ok {
		cached := val.(cachedPublic)
		if cached.projectID == nil || *cached.projectID != projectID {
			return fail.New(errx.KeysProjectKeyMismatch).RecordCtx(ctx)
		}
		if !ed25519.Verify(cached.pub, payload, sig) {
			return fail.New(errx.KeysInvalidSignature).RecordCtx(ctx)
		}
		return nil
	}

	pair, err := uc.keys.GetProjectKeyByKID(ctx, kid)
	if err != nil {
		return err
	}

	if pair.ProjectID == nil || *pair.ProjectID != projectID {
		return fail.New(errx.KeysProjectKeyMismatch).RecordCtx(ctx)
	}

	pub, err := parseEd25519Public(pair.PublicKey)
	if err != nil {
		return err
	}

	uc.publicKeyCache.Set(ctx, kid, cachedPublic{pub: pub, projectID: pair.ProjectID}, 0)

	if !ed25519.Verify(pub, payload, sig) {
		return fail.New(errx.KeysInvalidSignature).RecordCtx(ctx)
	}

	return nil
}

func (uc *CommandService) GetActiveGoAuthSigningKID(ctx context.Context) (string, error) {
	return uc.keys.GetActiveGoAuthSigningKID(ctx)
}

func (uc *CommandService) GetActiveProjectSigningKID(ctx context.Context, projectID uuid.UUID) (string, error) {
	return uc.keys.GetActiveProjectSigningKID(ctx, projectID)
}

func (uc *CommandService) RevokeKey(ctx context.Context, kid string) error {
	if err := uc.keys.RevokeKeyByKID(ctx, kid); err != nil {
		return err
	}

	uc.privateKeyCache.Delete(ctx, kid)
	uc.publicKeyCache.Delete(ctx, kid)

	return nil
}

func parseEd25519Private(pemBytes []byte) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid PEM Private Key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing PKCS8 Public key: %w", err)
	}

	priv, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not an ED25519 private key")
	}

	return priv, nil
}

func parseEd25519Public(pemStr string) (ed25519.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("invalid PEM Public Key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing PKIX Public key: %w", err)
	}

	pub, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("public key is not an ED25519 public key")
	}

	return pub, nil
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
