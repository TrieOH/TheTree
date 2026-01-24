package inbounds

type ErrParseProjectKey struct {
	KeyType string
	Cause   error
}

func (e ErrParseProjectKey) Error() string {
	return "failed to parse project " + e.KeyType + " key"
}

type ErrInvalidPEMPubKey struct{}

func (e ErrInvalidPEMPubKey) Error() string {
	return "invalid PEM public key"
}

type ErrInvalidPEMPrivKey struct{}

func (e ErrInvalidPEMPrivKey) Error() string {
	return "invalid PEM private key"
}

type ErrParsingPKIXPubKey struct {
	Cause error
}

func (e ErrParsingPKIXPubKey) Error() string {
	return "failed to parse PKIX public key"
}

type ErrParsingPKCS8PrivKey struct {
	Cause error
}

func (e ErrParsingPKCS8PrivKey) Error() string {
	return "failed to parse PKCS8 private key"
}

type ErrNotED25519PubKey struct{}

func (e ErrNotED25519PubKey) Error() string {
	return "not an ED25519 public key"
}

type ErrNotED25519PrivKey struct{}

func (e ErrNotED25519PrivKey) Error() string {
	return "not an ED25519 private key"
}

type ErrInvalidSignature struct{}

func (e ErrInvalidSignature) Error() string {
	return "invalid signature"
}

type ErrKeyProjectMismatch struct{}

func (e ErrKeyProjectMismatch) Error() string {
	return "key project mismatch"
}
