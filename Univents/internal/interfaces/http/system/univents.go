package system

import (
	"net/http"
	"strings"
	"time"
	domain2 "univents/internal/commerce/domain"
	"univents/internal/core/domain"
	"univents/internal/interfaces/http/system/dto"
	"univents/internal/shared/authz"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type UniventsHandler struct {
	gaClient *goauth.Client
}

func NewUniventsHandler(gaClient *goauth.Client) *UniventsHandler {
	return &UniventsHandler{
		gaClient: gaClient,
	}
}

// Health godoc
// @Summary Health check
// @Description Returns service health status
// @Tags system
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Router /health [get]
func (handler *UniventsHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := dto.HealthResponse{
		Status:  "ok",
		Service: "univents-api",
	}

	resp.OK("ok").WithData(response).Send(w)
}

// ProtectedHealth godoc
// @Summary Protected health check
// @Description Returns service health status and authenticated user id
// @Tags system
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Success 200 {object} dto.HealthResponse
// @Failure 401 {object} swag.ErrorResponse
// @Router /protected/health [get]
func (handler *UniventsHandler) ProtectedHealth(w http.ResponseWriter, r *http.Request) {
	sub, err := authz.RequireSubject(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	response := dto.HealthResponse{
		Status:  "ok",
		Service: "univents-api",
		UserID:  sub.ID,
	}

	resp.OK("ok").WithData(response).Send(w)
}

// Exchange godoc
// @Summary      Exchange global access token for service session
// @Description  Validates a global access token and creates a service session payload using a snapshot builder. Returns the session ID, TTL, and claims. Frontend sets the session cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer token, e.g., 'Bearer <token>'"
// @Success      200  {object}  goauth.SessionResult "Service session created"
// @Failure      400  {object}  object  "Bad request / empty token"
// @Failure      401  {object}  object  "Unauthorized / invalid token"
// @Failure      500  {object}  object  "Internal server error"
// @Router       /auth/exchange [post]
func (handler *UniventsHandler) Exchange(w http.ResponseWriter, r *http.Request) {
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

	sessionRes, err := handler.gaClient.ExchangeAndCreateSession(ctx, accessToken, domain.BuildServiceSnapshot)
	if err != nil {
		resp.Unauthorized("failed to create service session: " + err.Error()).Send(w)
		return
	}

	resp.OK("service session created").WithData(sessionRes).Send(w)
}

// WSAuth godoc
// @Summary Get WebSocket auth token
// @Description Returns a short-lived JWT (30s) to authenticate a WebSocket connection
// @Tags auth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Success 200 {object} object "Token generated"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /ws/token [get]
func (handler *UniventsHandler) WSAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	now := time.Now()
	claims := domain2.WSClaims{
		UserID: sub.ID,
		Email:  sub.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub.ID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	secret := viper.GetString("WS_JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		resp.InternalServerError("failed to sign token").Send(w)
		return
	}

	resp.OK("Token generated").WithData(map[string]string{
		"token": signed,
	}).Send(w)
}
