package system

import (
	"TrieForms/internal/shared/authz"
	"net/http"
	"strings"

	_ "TrieForms/internal/shared/contracts"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

// Exchange godoc
// @Summary      Exchange global access token for service session
// @Description  Validates a global access token and creates a service session payload using a snapshot builder. Returns the session ID, TTL, and claims. Frontend sets the session cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer token, e.g., 'Bearer <token>'"
// @Success      200  {object} goauth.SessionResult "Service session created"
// @Failure      400  {object} contracts.ErrorResponse
// @Failure      401  {object} contracts.ErrorResponse
// @Failure      500  {object} contracts.ErrorResponse
// @Router       /auth/exchange [post]
func (handler *SystemHandler) Exchange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		resp.Unauthorized().WithMsg("missing bearer token").Send(w)
		return
	}

	accessToken := strings.TrimSpace(authHeader[7:])
	if accessToken == "" {
		resp.Unauthorized().WithMsg("empty bearer token").Send(w)
		return
	}

	sessionRes, err := handler.gaClient.ExchangeAndCreateSession(ctx, accessToken, authz.BuildServiceSnapshot)
	if err != nil {
		resp.Unauthorized("failed to create service session: " + err.Error()).Send(w)
		return
	}

	resp.OK("service session created").WithData(sessionRes).Send(w)
}
