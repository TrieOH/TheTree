package contracts

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAPIKey_HappyPath(t *testing.T) {
	ownerID := uuid.New()

	ak, err := NewAPIKey(ownerID, "my-key", "hashvalue", "ik_")

	require.NoError(t, err)
	require.NotNil(t, ak)
	assert.Equal(t, ownerID, ak.OwnerID)
	assert.Equal(t, "my-key", ak.Name)
	assert.Equal(t, "hashvalue", ak.KeyHash)
	assert.Equal(t, "ik_", ak.KeyPrefix)
	assert.Equal(t, uuid.Nil, ak.ID)
	assert.True(t, ak.CreatedAt.IsZero())
	assert.Nil(t, ak.RevokedAt)
}

func TestNewAPIKey_KeyHash_NotExposedInJSON(t *testing.T) {
	ak, err := NewAPIKey(uuid.New(), "key", "supersecretvalue", "pk_")
	require.NoError(t, err)

	data, err := json.Marshal(ak)
	require.NoError(t, err)

	serialised := string(data)
	assert.NotContains(t, serialised, "supersecretvalue")
	assert.NotContains(t, serialised, "key_hash")
	assert.NotContains(t, serialised, "KeyHash")
}

func TestNewAPIKey_ValidationErrors(t *testing.T) {
	validOwner := uuid.New()

	tests := []struct {
		name      string
		ownerID   uuid.UUID
		keyName   string
		keyHash   string
		keyPrefix string
		wantErr   bool
	}{
		{
			name:      "empty name",
			ownerID:   validOwner,
			keyName:   "",
			keyHash:   "abc123hash",
			keyPrefix: "pk_",
			wantErr:   true,
		},
		{
			name:      "empty keyHash",
			ownerID:   validOwner,
			keyName:   "my-key",
			keyHash:   "",
			keyPrefix: "pk_",
			wantErr:   true,
		},
		{
			name:      "empty keyPrefix",
			ownerID:   validOwner,
			keyName:   "my-key",
			keyHash:   "abc123hash",
			keyPrefix: "",
			wantErr:   true,
		},
		{
			name:    "all string fields empty",
			ownerID: validOwner,
			wantErr: true,
		},
		{
			name:      "zero UUID owner — rejected by required tag",
			ownerID:   uuid.Nil,
			keyName:   "my-key",
			keyHash:   "abc123hash",
			keyPrefix: "pk_",
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ak, err := NewAPIKey(tc.ownerID, tc.keyName, tc.keyHash, tc.keyPrefix)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, ak)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, ak)
			}
		})
	}
}
