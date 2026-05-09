package validation

import (
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func GetUUID(r *http.Request, fieldName string) (uuid.UUID, *fun.Response) {
	IDStr := chi.URLParam(r, fieldName)
	if err := ValidateRule(IDStr, "required,uuid7", fieldName); err != nil {
		return uuid.Nil, fun.Error(err)
	}

	ID, err := ParseUUID(IDStr, fieldName)
	if err != nil {
		return uuid.Nil, fun.Error(err)
	}
	return ID, nil
}
