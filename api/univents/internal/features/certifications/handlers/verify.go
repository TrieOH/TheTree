package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
)

func (h *Handlers) Verify(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	hash := req.Path("hash").String()

	cert, err := h.queries.GetByHash(r.Context(), hash)
	if fun.Bail(w, err) {
		return
	}

	fun.Respond(w, models.VerifyCertificationResponse{
		IsVerified:  true,
		ID:          cert.ID,
		UserID:      cert.UserID,
		TargetID:    cert.TargetID,
		TargetType:  cert.TargetType,
		CertifiedAt: cert.CertifiedAt,
	})
}
