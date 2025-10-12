package handler

import (
	"net/http"

	resp "github.com/MintzyG/GoResponse/response"
)

func (h *AuthHandler) Hello(w http.ResponseWriter, r *http.Request) {
	resp.OK("Hello").Send(w)
}
