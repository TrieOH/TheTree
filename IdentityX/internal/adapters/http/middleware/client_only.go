package middleware

import (
	"GoAuth/internal/utils"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

func ClientOnly() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessTokenCookie, err := r.Cookie("access_token")
			if err != nil {
				resp.Unauthorized("missing access_token cookie").WithModule("AuthMW").Send(w)
				return
			}

			accessToken, err := utils.ParseAccessToken(accessTokenCookie.Value, utils.GoAuthPublicKey)
			if err != nil {
				ErrToResp(err).WithModule("AuthMW").Send(w)
				return
			}

			if accessToken.Sub.ProjectID != nil {
				resp.Unauthorized("only clients can access this endpoint").WithModule("ClientOnlyMW").Send(w)
				return
			}

			next.ServeHTTP(w, r.WithContext(r.Context()))
		})
	}
}
