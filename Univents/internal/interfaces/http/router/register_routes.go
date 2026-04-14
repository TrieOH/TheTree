package router

import (
	"univents/internal/features/activities"
	"univents/internal/features/checkpoints"
	"univents/internal/features/editions"
	"univents/internal/features/events"
	"univents/internal/features/products"
	"univents/internal/features/purchases"
	"univents/internal/features/tickets"
	"univents/internal/interfaces/http/middleware"
	"univents/internal/interfaces/http/system"

	"github.com/go-chi/chi/v5"
)

func registerRoutes(r *chi.Mux, deps *HTTPDeps) {
	registerSystemRoutes(r, deps.System, deps.AuthMiddleware)
	registerEventsRoutes(r, deps.Events, deps.AuthMiddleware)
	registerEditionsRoutes(r, deps.Editions, deps.AuthMiddleware)
	registerTicketsRoutes(r, deps.Tickets, deps.AuthMiddleware)
	registerActivitiesRoutes(r, deps.Activities, deps.AuthMiddleware)
	registerCheckpointsRoutes(r, deps.Checkpoints, deps.AuthMiddleware)
	registerProductsRoutes(r, deps.Products, deps.AuthMiddleware)
	registerPurchasesRoutes(r, deps.Purchases, deps.AuthMiddleware)
}

func registerSystemRoutes(
	r *chi.Mux,
	h *system.UniventsHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Post("/auth/exchange", h.Exchange)
		r.Get("/health", h.Health)
		r.With(authMW.Auth()).
			Get("/protected/health", h.ProtectedHealth)
		r.With(authMW.Auth()).
			Get("/ws/token", h.WSAuth)
	})
}

func registerEventsRoutes(
	r *chi.Mux,
	h *events.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Get("/events", h.ListEvents)
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Get("/events/own", h.ListOwnEvents)
		r.Post("/events", h.CreateEvent)
		r.Patch("/events/{event_id}", h.PatchEvent)
		r.Post("/events/{event_id}/publish", h.PublishEvent)
		r.Post("/events/{event_id}/gallery", h.AddGalleryImage)
		r.Delete("/events/{event_id}/gallery", h.RemoveGalleryImage)
		r.Put("/events/{event_id}/logo", h.SetLogo)
		r.Delete("/events/{event_id}/logo", h.UnsetLogo)
		r.Put("/events/{event_id}/banner", h.SetBanner)
		r.Delete("/events/{event_id}/banner", h.UnsetBanner)
	})
}

func registerEditionsRoutes(
	r *chi.Mux,
	h *editions.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Get("/events/{event_id}/editions", h.List)
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Get("/events/{event_id}/editions/admin", h.ListAdmin)
		r.Post("/events/{event_id}/editions", h.Create)
		r.Post("/events/{event_id}/editions/{edition_id}/announce", h.Announce)
		r.Post("/events/{event_id}/editions/{edition_id}/payments/connect", h.ConnectPaymentAccount)
		r.Post("/events/{event_id}/editions/{edition_id}/payments/disconnect", h.DisconnectPaymentAccount)
	})
}

func registerTicketsRoutes(
	r *chi.Mux,
	h *tickets.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Get("/events/{event_id}/editions/{edition_id}/tickets", h.List)
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Post("/events/{event_id}/editions/{edition_id}/tickets", h.Create)
		r.Post("/events/{event_id}/editions/{edition_id}/tickets/{ticket_id}/permissions", h.AddPermission)
		r.Delete("/events/{event_id}/editions/{edition_id}/tickets/{ticket_id}/permissions/{permission_id}", h.RemovePermission)
	})
}

func registerActivitiesRoutes(
	r *chi.Mux,
	h *activities.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Get("/events/{event_id}/editions/{edition_id}/activities", h.List)
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Post("/events/{event_id}/editions/{edition_id}/activities", h.Create)
		r.Get("/events/{event_id}/editions/{edition_id}/activities/admin", h.ListAdmin)
		r.Post("/events/{event_id}/editions/{edition_id}/activities/{activity_id}/publish", h.Publish)
		r.Post("/events/{event_id}/editions/{edition_id}/activities/{activity_id}/register", h.Register)
		r.Post("/events/{event_id}/editions/{edition_id}/activities/{activity_id}/unregister", h.Unregister)
		r.Get("/events/{event_id}/editions/{edition_id}/activities/{activity_id}/records", h.ListRecords)
		r.Post("/events/{event_id}/editions/{edition_id}/activities/{activity_id}/records/{record_id}", h.MarkAttendance)
	})
}

func registerCheckpointsRoutes(
	r *chi.Mux,
	h *checkpoints.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Post("/events/{event_id}/editions/{edition_id}/checkpoints", h.Create)
		r.Get("/events/{event_id}/editions/{edition_id}/checkpoints", h.List)
	})
}

func registerProductsRoutes(
	r *chi.Mux,
	h *products.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Get("/events/{event_id}/editions/{edition_id}/products", h.List)
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Post("/events/{event_id}/editions/{edition_id}/products", h.Create)
		r.Post("/events/{event_id}/editions/{edition_id}/products/{product_id}/publish", h.Publish)
		r.Get("/events/{event_id}/editions/{edition_id}/products/admin", h.ListAdmin)
		r.Delete("/events/{event_id}/editions/{edition_id}/products/{product_id}", h.Delete)
		r.Post("/events/{event_id}/editions/{edition_id}/products/{product_id}/restore", h.Restore)
		r.Post("/events/{event_id}/editions/{edition_id}/products/{product_id}/gallery", h.AddGalleryImage)
		r.Delete("/events/{event_id}/editions/{edition_id}/products/{product_id}/gallery", h.RemoveGalleryImage)
		r.Put("/events/{event_id}/editions/{edition_id}/products/{product_id}/thumbnail", h.SetThumbnail)
		r.Delete("/events/{event_id}/editions/{edition_id}/products/{product_id}/thumbnail", h.UnsetThumbnail)
		r.Get("/events/{event_id}/editions/{edition_id}/products/inventory/stream", h.StreamInventory) // SSE upgrade
	})
}

func registerPurchasesRoutes(
	r *chi.Mux,
	h *purchases.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Post("/webhooks/payments", h.WebhookHandler)
	r.Get("/events/{event_id}/editions/{edition_id}/products/purchase", h.Purchase) // WS upgrade
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Get("/purchases", h.ListUserPurchases)
		r.Get("/purchases/{purchase_id}/items", h.ListPurchaseItems)
	})
}
