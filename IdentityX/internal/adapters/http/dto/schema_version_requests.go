package dto

import "github.com/google/uuid"

type GetVersionVerboseRequest struct {
	VersionID *uuid.UUID `json:"version_id"`
}
