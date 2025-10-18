package middleware

import (
	"strings"
	"net/http"
	"context"
	"GoAuth/internal/utils"
	"GoAuth/internal/repository"
	"GoAuth/internal/models"
  "github.com/google/uuid"
	"github.com/spf13/viper"
	resp "github.com/MintzyG/GoResponse/response"
)

type AuthMiddleware struct {
	queries *repository.Queries
}

func NewAuthMiddleware(queries *repository.Queries) *AuthMiddleware {
	return &AuthMiddleware{queries: queries}
}

func (mw *AuthMiddleware) Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		access_token_cookie, err := r.Cookie("access_token")
		if err != nil {
			resp.Unauthorized("missing access_token cookie").WithModule("AuthMW").Send(w)
			return
		}

		refresh_token_cookie, err := r.Cookie("refresh_token")
		if err != nil {
			resp.Unauthorized("missing refresh_token cookie").WithModule("AuthMW").Send(w)
			return
		}

    access_token, rs := utils.ParseAccessToken(access_token_cookie.Value, viper.GetString("JWT_SECRET"))
		if rs != nil {
			rs.WithModule("AuthMW").Send(w)
			return
		}

		refresh_token, rs := utils.ParseRefreshToken(refresh_token_cookie.Value, viper.GetString("JWT_SECRET"))
		if rs != nil {
			rs.WithModule("AuthMW").Send(w)
			return
		}

    if access_token.Issuer != "GoAuth" || refresh_token.Issuer != "GoAuth" {
			resp.Unauthorized("Invalid Issuer").WithModule("AuthMW").Send(w)
			return
		}

		refreshUUID, err := uuid.Parse(refresh_token.ID)
		if err != nil {
			resp.Unauthorized("couln't parse refresh JTI").WithModule("AuthWM").Send(w)
			return
		}

		blacklisted, err := mw.queries.GetRefreshBlacklistById(r.Context(), refreshUUID)
		if err != nil && !strings.Contains(err.Error(), "no rows") {
			resp.Unauthorized("couldn't fetch refresh token").WithModule("AuthWM").WithTracePrefix("database-error").AddTrace(err).Send(w)
			return
		}

		if blacklisted.TokenID == refreshUUID {
			resp.Unauthorized("refresh token is invalidated").WithModule("AuthWM").Send(w)
			return
		}

		ctx := context.WithValue(r.Context(), models.AccessClaimsKey, access_token)
		ctx = context.WithValue(ctx, models.RefreshClaimsKey, refresh_token)

		h.ServeHTTP(w, r)
	}
}
