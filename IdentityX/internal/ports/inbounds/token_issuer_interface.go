package inbounds

import (
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/user"
	"crypto/ed25519"
	"time"

	"github.com/google/uuid"
)

type TokenIssuer interface {
	NewAccessToken(user user.User, key ed25519.PrivateKey, ip, agent, accessJTI, keyID string, sessionID uuid.UUID, expiresAt time.Time) (string, error)
	NewRefreshToken(keyID string, privKey ed25519.PrivateKey, accessJTI, refreshJTI uuid.UUID, expiresAt time.Time) (string, error)
	NewProjectAccessToken(user project_users.ProjectUser, ip, agent, accessJTI, keyID string, sessionID uuid.UUID, expiresAt time.Time, privKey ed25519.PrivateKey) (string, error)
}
