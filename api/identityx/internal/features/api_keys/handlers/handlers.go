package handlers

import (
	"github.com/google/uuid"

	"IdentityX/internal/features/api_keys/service"
	"IdentityX/internal/openapi"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

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
