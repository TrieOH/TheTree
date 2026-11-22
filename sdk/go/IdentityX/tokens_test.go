package idx

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

// verifier serves one Ed25519 key as the JWKS and signs tokens with it.
type verifier struct {
	client *Client
	priv   ed25519.PrivateKey
	kid    string
}

func newVerifier(t *testing.T) *verifier {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	kid := uuid.New().String()
	jwks := JWKS{Keys: []JWK{{Kty: "OKP", Crv: "Ed25519", X: base64.RawURLEncoding.EncodeToString(pub), Kid: kid}}}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwks)
	}))
	t.Cleanup(srv.Close)

	client, err := NewClient(Config{BaseURL: srv.URL, APIKey: "test-api-key", ProjectID: uuid.New()})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return &verifier{client: client, priv: priv, kid: kid}
}

func (v *verifier) sign(t *testing.T, claims jwt.Claims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = v.kid
	signed, err := token.SignedString(v.priv)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return signed
}

func registered() jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Issuer:    "https://identityx.com.br",
		ID:        uuid.New().String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
}

func TestVerifyAccessTokenAcceptsAccessToken(t *testing.T) {
	v := newVerifier(t)
	id := uuid.New()
	token := v.sign(t, &AccessClaims{RegisteredClaims: registered(), Sub: AccessSub{ID: id, Type: "human"}})

	claims, err := v.client.Tokens.VerifyAccessToken(context.Background(), token)
	if err != nil {
		t.Fatalf("VerifyAccessToken: %v", err)
	}
	if claims.Sub.ID != id {
		t.Fatalf("subject id = %v, want %v", claims.Sub.ID, id)
	}
}

func TestVerifyAccessTokenRejectsRefreshToken(t *testing.T) {
	v := newVerifier(t)
	token := v.sign(t, &RefreshClaims{RegisteredClaims: registered(), Sub: RefreshSub{ID: uuid.New(), AccessJTI: uuid.New()}})

	_, err := v.client.Tokens.VerifyAccessToken(context.Background(), token)
	if _, ok := errors.AsType[*InvalidTokenError](err); !ok {
		t.Fatalf("want InvalidTokenError for a refresh token, got %v", err)
	}
}
