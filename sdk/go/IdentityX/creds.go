package idx

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
)

// CredentialHandler encrypts and persists API credentials to a local file.
// Create one with NewCredentialHandler and attach it to Client.Creds. Call
// SaveCreds to persist an API key + project ID, and GetCreds to recover them.
//
// The encryption key is derived from password via SHA-256. Data is encrypted
// with AES-256-GCM using a random 12-byte nonce per write. Writes are atomic
// (temp file + rename) and the backing file is created with 0600 permissions.
type CredentialHandler struct {
	store *credentialStore
}

// credentialStore is the unexported engine that does the actual crypto I/O.
type credentialStore struct {
	path string
	key  [32]byte
}

type credsData struct {
	APIKey    string    `json:"api_key"`
	ProjectID uuid.UUID `json:"project_id"`
}

// NewCredentialHandler returns a CredentialHandler that stores credentials at
// the given file path, encrypted with a key derived from password.
func NewCredentialHandler(path string, password []byte) *CredentialHandler {
	key := sha256.Sum256(password)
	return &CredentialHandler{
		store: &credentialStore{path: path, key: key},
	}
}

// SaveCreds encrypts apiKey and projectID and writes them atomically to the
// backing file.
func (h *CredentialHandler) SaveCreds(apiKey string, projectID uuid.UUID) error {
	return h.store.save(apiKey, projectID)
}

// GetCreds reads and decrypts credentials from the backing file. Returns the
// API key and project ID. If the file does not exist, it returns zero values
// without an error.
func (h *CredentialHandler) GetCreds() (string, uuid.UUID, error) {
	return h.store.load()
}

// ---------------------------------------------------------------------------
// internal engine
// ---------------------------------------------------------------------------

func (s *credentialStore) save(apiKey string, projectID uuid.UUID) error {
	data := credsData{APIKey: apiKey, ProjectID: projectID}
	plaintext, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("idx: creds: marshal: %w", err)
	}

	block, err := aes.NewCipher(s.key[:])
	if err != nil {
		return fmt.Errorf("idx: creds: cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("idx: creds: gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("idx: creds: nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, ciphertext, 0600); err != nil {
		return fmt.Errorf("idx: creds: write: %w", err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return fmt.Errorf("idx: creds: rename: %w", err)
	}
	return nil
}

func (s *credentialStore) load() (string, uuid.UUID, error) {
	ciphertext, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", uuid.Nil, nil
		}
		return "", uuid.Nil, fmt.Errorf("idx: creds: read: %w", err)
	}

	block, err := aes.NewCipher(s.key[:])
	if err != nil {
		return "", uuid.Nil, fmt.Errorf("idx: creds: cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", uuid.Nil, fmt.Errorf("idx: creds: gcm: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", uuid.Nil, fmt.Errorf("idx: creds: file too short (%d bytes)", len(ciphertext))
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", uuid.Nil, fmt.Errorf("idx: creds: decrypt: wrong password or corrupt file")
	}

	var data credsData
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return "", uuid.Nil, fmt.Errorf("idx: creds: unmarshal: %w", err)
	}
	return data.APIKey, data.ProjectID, nil
}
