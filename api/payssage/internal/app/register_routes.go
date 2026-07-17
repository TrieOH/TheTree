package app

//
//func registerIntentsRoutes(
//	r *chi.Mux,
//	h *intents.Handler,
//	authMW *AuthMiddleware,
//) {
//	r.With(authMW.AnyAuth()).Get("/intents", h.List)
//	r.With(authMW.AnyAuth()).Get("/workspaces/{name}/intents", h.ListByWorkspace)
//	r.Group(func(r chi.Router) {
//		r.Use(authMW.APIKey())
//		r.Post("/intents", h.InitiateCheckout)
//		r.Get("/intents/{intent_id}", h.GetByID)
//		r.Post("/intents/{intent_id}/cancel", h.CancelIntent)
//		r.Post("/intents/{intent_id}/cancel-pix", h.CancelPix)
//		r.Post("/intents/{intent_id}/charge", h.Charge)
//	})
//}
//
//func registerWebhookRoutes(
//	r *chi.Mux,
//	h *webhooks.Handler,
//	authMW *AuthMiddleware,
//) {
//	// inbound from providers — no auth, verified by signature
//	r.Post("/webhooks/{provider}", h.HandleProviderWebhook)
//	r.Group(func(r chi.Router) {
//		r.Use(authMW.Auth())
//		r.Post("/workspaces/{name}/webhooks", h.RegisterWebhookEndpoint)
//		r.Get("/workspaces/{name}/webhooks", h.ListWebhookEndpoints)
//		r.Get("/workspaces/{name}/webhook-events", h.ListWebhookEvents)
//		r.Delete("/workspaces/{name}/webhooks/{endpoint_id}", h.DeleteWebhookEndpoint)
//		r.Get("/workspaces/{name}/webhooks/{endpoint_id}/deliveries", h.ListWebhookDeliveries)
//	})
//}
//
//func registerOAuthRoutes(
//	r *chi.Mux,
//	h *oauth.Handler,
//	authMW *AuthMiddleware,
//) {
//	// callback from provider — no auth, browser redirect
//	r.Get("/oauth/{provider}/callback", h.CompleteOAuth)
//	r.With(authMW.APIKey()).
//		Post("/workspaces/{name}/providers/{provider}/connect", h.ConnectSeller)
//
//	r.Group(func(r chi.Router) {
//		r.Use(authMW.Auth())
//		r.Post("/workspaces/{name}/providers/{provider}/setup", h.SetupProvider)
//		r.Get("/workspaces/{name}/marketplaces", h.ListMarketplaceConfigs)
//		r.Put("/workspaces/{name}/marketplaces", h.SetMarketplaceConfig)
//		r.Delete("/workspaces/{name}/marketplaces/{credential_id}", h.DeleteMarketplaceConfig)
//		r.Delete("/workspaces/{name}/providers/{credential_id}", h.RevokeProvider)
//	})
//
//	r.Group(func(r chi.Router) {
//		r.Use(authMW.APIKey())
//		r.Delete("/workspaces/{name}/providers/{credential_id}/disconnect", h.DisconnectProvider)
//	})
//}
