package authz

import (
	"encoding/json"

	"github.com/MintzyG/fail/v3"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/google/uuid"
)

// SnapshotPayload is the typed payload stored in the service session
type SnapshotPayload struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

// BuildServiceSnapshot builds a typed snapshot for the service session
func BuildServiceSnapshot(claims *goauth.AccessClaims) ([]byte, error) {
	payload := SnapshotPayload{
		UserID: claims.Sub.ID,
		Email:  claims.Sub.Email,
	}
	return json.Marshal(payload)
}

// UnmarshalSnapshot unmarshals the session bytes into a typed payload
func UnmarshalSnapshot(data []byte) (*SnapshotPayload, error) {
	if data == nil {
		return nil, fail.New(goauth.SDKUnknownErrorID).WithArgs("session data is nil")
	}

	var payload SnapshotPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fail.New(goauth.SDKUnknownErrorID).WithArgs("failed to unmarshal session: " + err.Error())
	}

	return &payload, nil
}
