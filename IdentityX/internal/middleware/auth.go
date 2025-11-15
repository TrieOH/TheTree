package middleware

import (
	"GoAuth/internal/models"
	"GoAuth/internal/repository"
	"GoAuth/internal/utils"
	"context"
	"net/http"
	"strings"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type AuthMiddleware struct {
	queries *repository.Queries
}

func NewAuthMiddleware(queries *repository.Queries) *AuthMiddleware {
	return &AuthMiddleware{queries: queries}
}

func (mw *AuthMiddleware) Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessTokenCookie, err := r.Cookie("access_token")
		if err != nil {
			resp.Unauthorized("missing access_token cookie").WithModule("AuthMW").Send(w)
			return
		}

		refreshTokenCookie, err := r.Cookie("refresh_token")
		if err != nil {
			resp.Unauthorized("missing refresh_token cookie").WithModule("AuthMW").Send(w)
			return
		}

		accessToken, rs := utils.ParseAccessToken(accessTokenCookie.Value, viper.GetString("JWT_SECRET"))
		if rs != nil {
			rs.WithModule("AuthMW").Send(w)
			return
		}

		refreshToken, rs := utils.ParseRefreshToken(refreshTokenCookie.Value, viper.GetString("JWT_SECRET"))
		if rs != nil {
			rs.WithModule("AuthMW").Send(w)
			return
		}

		if accessToken.Issuer != "GoAuth" || refreshToken.Issuer != "GoAuth" {
			resp.Unauthorized("Invalid Issuer").WithModule("AuthMW").Send(w)
			return
		}

		refreshUUID, err := uuid.Parse(refreshToken.ID)
		if err != nil {
			resp.Unauthorized("couldn't parse refresh JTI").WithModule("AuthMW").Send(w)
			return
		}

		blacklisted, err := mw.queries.GetRefreshBlacklistById(r.Context(), refreshUUID)
		if err != nil && !strings.Contains(err.Error(), "no rows") {
			resp.Unauthorized("couldn't fetch refresh token").WithModule("AuthMW").WithTracePrefix("database-error").AddTrace(err).Send(w)
			return
		}

		if blacklisted.TokenID == refreshUUID {
			resp.Unauthorized("refresh token is invalidated").WithModule("AuthMW").Send(w)
			return
		}

		ctx := context.WithValue(r.Context(), models.AccessClaimsKey, accessToken)
		ctx = context.WithValue(ctx, models.RefreshClaimsKey, refreshToken)

		h.ServeHTTP(w, r.WithContext(ctx))
	}
}
