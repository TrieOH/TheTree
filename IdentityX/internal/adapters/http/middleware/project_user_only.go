package middleware

import (
	"GoAuth/internal/utils"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

func ProjectUserOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if accessToken.Sub.ProjectID == nil {
			resp.Unauthorized("only project users can access this endpoint").WithModule("ProjectUserOnlyMW").Send(w)
			return
		}

		h.ServeHTTP(w, r.WithContext(r.Context()))
	}
}
