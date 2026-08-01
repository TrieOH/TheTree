package api_keys

import (
	"github.com/google/uuid"

	"IdentityX/internal/openapi"
	"IdentityX/internal/services"
)

type Handlers struct {
	ops *services.APIKeys
}

func New(ops *services.APIKeys) *Handlers { return &Handlers{ops: ops} }

const module = "IdentityX"

func deref(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func derefSlice(v *[]openapi.UUID) []uuid.UUID {
	if v == nil {
		return nil
	}
	return *v
}
