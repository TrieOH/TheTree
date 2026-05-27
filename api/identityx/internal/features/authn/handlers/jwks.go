package handlers

import (
	"encoding/json"
	"lib/globals"
	"lib/telemetry"
	"net/http"

	"github.com/MintzyG/fun"
	"go.uber.org/zap"
)

// JWKS godoc
// @Summary Get a JWKS set
// @Description Provides the JSON Web Key Set (JWKS) for verifying JWTs issued by the authentication service.
// @Tags authn
// @ID authn_jwks
// @Accept json
// @Produce json
// @Success 200 {object} object "JSON Web Key Set (JWKS)"
// @Failure 500 {object} fun.Response "Internal Server Error"
// @Failure 503 {object} fun.Response "Internal Server Error"
// @Router /.well-known/jwks.json [get]
func (h *Handlers) JWKS(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	jwks, err := h.queries.JWKS(r.Context())
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
	if _, err = w.Write(data); err != nil {
		telemetry.Log().Error("failed to write JWKS response", zap.Error(err))
	}
}
