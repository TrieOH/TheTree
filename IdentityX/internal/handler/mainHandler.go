package handler

import (
	resp "github.com/MintzyG/GoResponse/response"
	"net/http"
	"GoAuth/internal/service"
)

type GoAuthHandler struct {
	GoAuthService *service.GoAuthService
}

func NewGoAuthHandler(service *service.GoAuthService) *GoAuthHandler {
	return &GoAuthHandler{GoAuthService: service}
}

func (h *GoAuthHandler) Hi(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	name := queryParams.Get("name")

	normalizedName := h.GoAuthService.Hi(name)

	resp.OK("Hi " + normalizedName).WithData(normalizedName).WithModule("Auth").Send(w)
}
