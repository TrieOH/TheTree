package handlers

import (
	"encoding/json"
	"lib/globals"
	"lib/telemetry"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *Handlers) JWKS(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}

	projectID, ok := fun.From(r).Query("project_id").UUIDOpt()
	if !ok && r.URL.Query().Get("project_id") != "" {
		fun.BadRequest("invalid project_id").Send(w)
		return
	}

	var pid *uuid.UUID
	if ok {
		pid = &projectID
	}

	jwks, err := h.queries.JWKS(r.Context(), pid)
	if fun.Bail(w, err) {
		return
	}
	data, err := json.Marshal(jwks)
	if fun.Bail(w, err) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=7200")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		telemetry.Log().Error("failed to write JWKS response", zap.Error(err))
	}
}
