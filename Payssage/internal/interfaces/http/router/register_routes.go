package router

import (
	apiKeys "TriePayments/internal/core/interfaces/http/api_keys_handler"
	intents "TriePayments/internal/core/interfaces/http/intent_handler"
	"TriePayments/internal/core/interfaces/http/oauth_handler"
	webhooks "TriePayments/internal/core/interfaces/http/webhooks_handler"
	workspaces "TriePayments/internal/core/interfaces/http/workspaces_handler"
	"TriePayments/internal/interfaces/http/middleware"
	"TriePayments/internal/interfaces/http/system"

	"github.com/go-chi/chi/v5"
)

func registerRoutes(r *chi.Mux, deps *HTTPDeps) {
	registerSystemRoutes(r, deps.SystemHandler, deps.AuthMiddleware)
	registerIntentsRoutes(r, deps.IntentsHandler, deps.AuthMiddleware)
	registerWorkspacesRoutes(r, deps.WorkspacesHandler, deps.AuthMiddleware)
	registerApiKeysRoutes(r, deps.ApiKeysHandler, deps.AuthMiddleware)
	registerWebhookRoutes(r, deps.WebhooksHandler, deps.AuthMiddleware)
	registerOAuthRoutes(r, deps.OauthHandler, deps.AuthMiddleware)
}

func registerSystemRoutes(
	r *chi.Mux,
	h *system.SystemHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Get("/health", h.Health)
		r.With(authMW.Auth()).
			Get("/protected/health", h.ProtectedHealth)
	})
}

func registerWorkspacesRoutes(
	r *chi.Mux,
	h *workspaces.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Post("/workspaces", h.Create)
		r.Get("/workspaces", h.List)
		r.Post("/workspaces/{name}/sandbox/enable", h.EnableSandbox)
		r.Post("/workspaces/{name}/sandbox/disable", h.DisableSandbox)
	})
}

func registerApiKeysRoutes(
	r *chi.Mux,
	h *apiKeys.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Post("/workspaces/{name}/keys", h.Create)
		r.Get("/workspaces/{name}/keys", h.ListAPIKeys)
		r.Delete("/workspaces/{name}/keys/{id}", h.RevokeAPIKey)
	})
}

func registerIntentsRoutes(
	r *chi.Mux,
	h *intents.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.With(authMW.AnyAuth()).Get("/intents", h.List)
	r.Group(func(r chi.Router) {
		r.Use(authMW.APIKey())
		r.Post("/intents", h.CreateIntent)
		r.Get("/intents/{intent_id}", h.GetByID)
		r.Post("/intents/{intent_id}/cancel", h.CancelIntent)
		r.Post("/intents/{intent_id}/pay", h.PayIntent)
	})
}

func registerWebhookRoutes(
	r *chi.Mux,
	h *webhooks.Handler,
	authMW *middleware.AuthMiddleware,
) {
	// inbound from providers — no auth, verified by signature
	r.Post("/webhooks/{provider}", h.HandleProviderWebhook)
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Post("/workspaces/{name}/webhooks", h.RegisterWebhookEndpoint)
		r.Get("/workspaces/{name}/webhooks", h.ListWebhookEndpoints)
		r.Delete("/workspaces/{name}/webhooks/{endpoint_id}", h.DeleteWebhookEndpoint)
	})
}

func registerOAuthRoutes(
	r *chi.Mux,
	h *oauth_handler.Handler,
	authMW *middleware.AuthMiddleware,
) {
	// callback from provider — no auth, browser redirect
	r.Get("/oauth/{provider}/callback", h.CompleteOAuth)

	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Post("/workspaces/{name}/providers/{provider}/setup", h.SetupProvider)
		r.Post("/workspaces/{name}/providers/{provider}/connect", h.ConnectSeller)
		r.Put("/workspaces/{name}/marketplace", h.SetMarketplaceConfig)
		r.Delete("/workspaces/{name}/marketplace", h.DeleteMarketplaceConfig)
	})
}
