package app

import (
	"net/http"

	"lib/httpserver"
	"univents/internal/handlers"
	"univents/internal/handlers/webhooks"
	"univents/internal/openapi"

	"github.com/go-chi/chi/v5"
)

// mountRaw mounts the realtime routes (split 6) — deliberately outside the
// strict handler: the WS handshake cannot carry Authorization headers (the
// one-time query token is the auth), and SSE must stream without the
// fun/validate envelope machinery buffering the body. Both are documented in
// the spec's description block, not as spec ops.
func mountRaw(r chi.Router, h *handlers.Server) {
	r.Get("/editions/{edition_id}/store/stream", h.ServeStoreStream)
	r.Handle("/ws", http.HandlerFunc(h.ServeWS))
}

// mountStrict is the backend's one strict-server mount point: the generated
// package's constructors cross the seam as values, and every piece of
// policy — middleware order, fail-closed dispatch, unified param-binding
// error mapping — lives in httpserver.MountStrict. The raw-request capture
// middleware runs first so the Payssage webhook can verify the signature
// against the exact body bytes payssage POSTed (the strict server decodes
// the body afterwards; D2). Path-scoped to /webhooks/ inside the middleware.
func mountStrict(r chi.Router, h openapi.StrictServerInterface, chains map[string][]func(http.Handler) http.Handler) {
	httpserver.MountStrict[openapi.StrictHandlerFunc](
		h,
		chains,
		openapi.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  httpserver.StrictRequestErrorHandler(),
			ResponseErrorHandlerFunc: httpserver.StrictResponseErrorHandler(),
		},
		openapi.NewStrictHandlerWithOptions,
		openapi.ChiServerOptions{
			BaseRouter:       r,
			Middlewares:      []openapi.MiddlewareFunc{webhooks.RawRequestMiddleware},
			ErrorHandlerFunc: httpserver.ParamBindingErrorHandler(),
		},
		openapi.HandlerWithOptions,
	)
}
