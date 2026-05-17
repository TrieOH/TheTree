package security

import (
	"net/http"
	"time"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	"univents/internal/shared/sockets"

	"github.com/MintzyG/fun"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Handler struct {
	Registry *sockets.Registry
}

func NewHandler(
	registry *sockets.Registry,
) *Handler {
	return &Handler{
		Registry: registry,
	}
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
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /ws/token [get]
func (handler *Handler) WSAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sub, err := authz.RequireSubject(ctx)
	if fun.Bail(w, err) {
		return
	}

	now := time.Now()
	claims := contracts.WSClaims{
		UserID: sub.ID,
		Email:  sub.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub.ID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	secret := viper.GetString("WS_JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		fun.InternalServerError("failed to sign token").Send(w)
		return
	}

	fun.OK("Token generated").WithData(map[string]string{
		"token": signed,
	}).Send(w)
}
