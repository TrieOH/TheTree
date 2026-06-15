package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/MintzyG/fun"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = fun.ErrInternal("argon2: invalid hash format")
	ErrIncompatibleVersion = fun.ErrInternal("argon2: incompatible version")
	ErrMismatch            = fun.ErrBadRequest("argon2: password does not match hash")
)

type Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// Default is a safe baseline for most use cases.
var Default = Params{
	Memory:      64 * 1024, // 64mb
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// Strong is for high-value credentials like super admin passwords.
var Strong = Params{
	Memory:      128 * 1024, // 128mb
	Iterations:  4,
	Parallelism: 4,
	SaltLength:  16,
	KeyLength:   32,
}

func Hash(password string, p Params) (string, error) {
	salt := make([]byte, p.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fun.ErrInternal("argon2: failed to generate salt: " + err.Error())
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		p.Iterations,
		p.Memory,
		p.Parallelism,
		p.KeyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.Memory,
		p.Iterations,
		p.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}

func Verify(password, encoded string) error {
	p, salt, hash, err := decode(encoded)
	if err != nil {
		return err
	}

	candidate := argon2.IDKey(
		[]byte(password),
		salt,
		p.Iterations,
		p.Memory,
		p.Parallelism,
		p.KeyLength,
	)

	if subtle.ConstantTimeCompare(hash, candidate) != 1 {
		return ErrMismatch
	}

	return nil
}

// NeedsRehash returns true if the hash was produced with different params,
// useful for transparently upgrading hashes on login.
func NeedsRehash(encoded string, p Params) (bool, error) {
	current, _, _, err := decode(encoded)
	if err != nil {
		return false, err
	}

	return current.Memory != p.Memory ||
		current.Iterations != p.Iterations ||
		current.Parallelism != p.Parallelism ||
		current.KeyLength != p.KeyLength, nil
}

func decode(encoded string) (Params, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return Params{}, nil, nil, ErrInvalidHash
	}

	if parts[1] != "argon2id" {
		return Params{}, nil, nil, ErrInvalidHash
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return Params{}, nil, nil, ErrInvalidHash
	}
	if version != argon2.Version {
		return Params{}, nil, nil, ErrIncompatibleVersion
	}

	var p Params
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism); err != nil {
		return Params{}, nil, nil, ErrInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return Params{}, nil, nil, ErrInvalidHash
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return Params{}, nil, nil, ErrInvalidHash
	}

	p.KeyLength = uint32(len(hash))
	p.SaltLength = uint32(len(salt))

	return p, salt, hash, nil
}
