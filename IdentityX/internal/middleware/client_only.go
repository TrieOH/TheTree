package middleware

import (
	"GoAuth/internal/utils"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

func ClientOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessTokenCookie, err := r.Cookie("access_token")
		if err != nil {
			resp.Unauthorized("missing access_token cookie").WithModule("AuthMW").Send(w)
			return
		}

		accessToken, rs := utils.ParseAccessToken(accessTokenCookie.Value, utils.GoAuthPublicKey)
		if rs != nil {
			rs.WithModule("AuthMW").Send(w)
			return
		}

		if accessToken.Sub.ProjectId != nil {
			resp.Unauthorized("only clients can access this endpoint").WithModule("ClientOnlyMW").Send(w)
			return
		}

		h.ServeHTTP(w, r.WithContext(r.Context()))
	}
}
