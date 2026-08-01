package handlers

import (
	"net/http"
	"univents/internal/features/products"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *products.Operations
}

func NewHandlers(ops *products.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwt func(http.Handler) http.Handler,
) {
	r.Get("/editions/{edition_id}/products", h.ListByEdition)
	r.Get("/editions/{edition_id}/products/{vendor_code}:by-code", h.GetByVendorCode)
	r.Get("/editions/{edition_id}/variants/{vendor_code}:by-code", h.GetVariantByVendorCode)
	r.Get("/products/{product_id}", h.GetByID)
	r.Get("/products/{product_id}/variants", h.ListVariants)
	r.With(jwt).Post("/editions/{edition_id}/products", h.CreateInitial)
	r.With(jwt).Post("/products/{product_id}/variants", h.CreateVariant)
	r.With(jwt).Patch("/products/{product_id}", h.PatchProduct)
	r.With(jwt).Patch("/variants/{variant_id}", h.PatchVariant)
	r.With(jwt).Delete("/products/{product_id}", h.DeleteProduct)
	r.With(jwt).Delete("/variants/{variant_id}", h.DeleteVariant)
}
