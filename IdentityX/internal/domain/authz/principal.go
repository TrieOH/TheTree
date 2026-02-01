package authz

import (
	"GoAuth/internal/domain/auth"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Principal struct {
	// ===== Identity =====
	UserID     uuid.UUID
	Email      string
	UserType   string
	ProjectID  *uuid.UUID
	Metadata   *json.RawMessage
	IsVerified bool
	VerifiedAt *time.Time

	// ===== Session =====
	SessionID uuid.UUID
	UserAgent string
	UserIP    string

	// ===== Token linkage =====
	AccessJTI  uuid.UUID
	RefreshJTI uuid.UUID

	// ===== Raw claims (escape hatch) =====
	AccessClaims  *auth.AccessClaims
	RefreshClaims *auth.RefreshClaims
}
