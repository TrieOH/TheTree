package router

import (
	"univents/internal/commerce/interfaces/http/products"
	"univents/internal/commerce/interfaces/http/tickets"
	eventhttp "univents/internal/core/interfaces/http"
	activityhttp "univents/internal/core/interfaces/http/activities"
	"univents/internal/core/interfaces/http/checkpoints"
	editionhttp "univents/internal/core/interfaces/http/editions"
	"univents/internal/interfaces/http/middleware"
	systemhttp "univents/internal/interfaces/http/system"

	"github.com/go-chi/chi/v5"
)

func registerRoutes(r *chi.Mux, deps *HTTPDeps) {
	registerSystemRoutes(r, deps.SystemHandler, deps.AuthMiddleware)
	registerEventsRoutes(r, deps.EventsHandler, deps.AuthMiddleware)
	registerEditionsRoutes(r, deps.EditionsHandler, deps.AuthMiddleware)
	registerTicketsRoutes(r, deps.TicketsHandler, deps.AuthMiddleware)
	registerActivitiesRoutes(r, deps.ActivitiesHandler, deps.AuthMiddleware)
	registerCheckpointsRoutes(r, deps.CheckpointsHandler, deps.AuthMiddleware)
	registerProductsRoutes(r, deps.ProductsHandler, deps.AuthMiddleware)
}

func registerSystemRoutes(
	r *chi.Mux,
	h *systemhttp.UniventsHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Get("/health", h.Health)
		r.With(authMW.Auth()).
			Get("/protected/health", h.ProtectedHealth)
	})
}

func registerEventsRoutes(
	r *chi.Mux,
	h *eventhttp.EventsHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Get("/events", h.ListEvents)
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Get("/events/own", h.ListOwnEvents)
		r.Post("/events", h.CreateEvent)
		r.Patch("/events/{event_id}", h.PatchEvent)
		r.Post("/events/{event_id}/publish", h.PublishEvent)
	})
}

func registerEditionsRoutes(
	r *chi.Mux,
	h *editionhttp.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Get("/events/{event_id}/editions", h.List)
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Get("/events/{event_id}/editions/admin", h.ListAdmin)
		r.Post("/events/{event_id}/editions", h.Create)
		r.Post("/events/{event_id}/editions/{edition_id}/announce", h.Announce)
	})
}

func registerTicketsRoutes(
	r *chi.Mux,
	h *tickets.TicketsHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Get("/events/{event_id}/editions/{edition_id}/tickets", h.List)
		r.Post("/events/{event_id}/editions/{edition_id}/tickets", h.Create)
		r.Post("/events/{event_id}/editions/{edition_id}/tickets/{ticket_id}/permissions", h.AddPermission)
		r.Delete("/events/{event_id}/editions/{edition_id}/tickets/{ticket_id}/permissions/{permission_id}", h.RemovePermission)
	})
}

func registerActivitiesRoutes(
	r *chi.Mux,
	h *activityhttp.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Post("/events/{event_id}/editions/{edition_id}/activities", h.Create)
		r.Get("/events/{event_id}/editions/{edition_id}/activities", h.List)
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
	r.Post("/webhooks/payments", h.WebhookHandler)
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Post("/events/{event_id}/editions/{edition_id}/products", h.Create)
		r.Post("/events/{event_id}/editions/{edition_id}/products/{product_id}/publish", h.Publish)
		r.Get("/events/{event_id}/editions/{edition_id}/products", h.List)
		r.Get("/events/{event_id}/editions/{edition_id}/products/admin", h.ListAdmin)
		r.Get("/events/{event_id}/editions/{edition_id}/products/purchase", h.Purchase) // WS upgrade
		r.Get("/purchases", h.ListUserPurchases)
		r.Get("/purchases/{purchase_id}/items", h.ListPurchaseItems)
	})
}
