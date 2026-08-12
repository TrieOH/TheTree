package app

import (
	"errors"
	"net/http"

	"lib/errx"
	"lib/httpserver"
	spec "univents"
	"univents/internal/handlers"
	"univents/internal/handlers/webhooks"
	"univents/internal/openapi"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
	"riverqueue.com/riverui"
)

func (app *Univents) CreateRouter(middlewares middlewares, h *handlers.Server, riverUIHandler *riverui.Handler) http.Handler {
	chains, err := resolveAuthChains(middlewares)
	errx.Exit(err, "resolve auth chains")
	return httpserver.NewRouter(httpserver.Config{
		AppName:            app.cfg.AppName,
		CorsAllowedOrigins: app.cfg.CorsAllowedOrigins,
		CorsAllowedHeaders: app.cfg.CorsAllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
		SkipLogPrefixes:    []string{"/admin/asynq"},
		Routes: func(r *chi.Mux) {
			mountStrict(r, h, chains)

			// Raw realtime routes (split 6) — deliberately outside the strict
			// handler: the WS handshake cannot carry Authorization headers (the
			// one-time query token is the auth), and SSE must stream without the
			// fun/validate envelope machinery buffering the body. Both are
			// documented in the spec's description block, not as spec ops.
			r.Get("/editions/{edition_id}/store/stream", h.ServeStoreStream)
			r.Handle("/ws", http.HandlerFunc(h.ServeWS))

			r.Group(func(r chi.Router) {
				r.Mount("/riverui", riverUIHandler)
			})
		},
	})
}

// mountStrict registers the generated strict handler on r with the harness's
// validation + auth middleware stack and fun-envelope error handlers. Only
// the generated-type conversions and the param-binding error mapping stay
// here; the rest lives in lib/httpserver.
func mountStrict(r *chi.Mux, h *handlers.Server, chains map[string][]func(http.Handler) http.Handler) {
	strict := openapi.NewStrictHandlerWithOptions(h,
		[]openapi.StrictMiddlewareFunc{
			adapt(httpserver.ValidateMiddleware()),
			adapt(httpserver.AuthDispatch(chains)),
		},
		openapi.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  httpserver.StrictRequestErrorHandler(),
			ResponseErrorHandlerFunc: httpserver.StrictResponseErrorHandler(),
		})
	openapi.HandlerWithOptions(strict, openapi.ChiServerOptions{
		BaseRouter: r,
		// The raw-request capture runs first so the Payssage webhook can
		// verify the signature against the exact body bytes payssage POSTed
		// (the strict server decodes the body afterwards; D2). Path-scoped
		// to /webhooks/ inside the middleware.
		Middlewares: []openapi.MiddlewareFunc{webhooks.RawRequestMiddleware},
		ErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
			var required *openapi.RequiredParamError
			var invalid *openapi.InvalidParamFormatError
			switch {
			case errors.As(err, &required):
				fun.Error(fun.Err("invalid request parameter").WithFields(&fun.FieldError{Field: required.ParamName, Message: "parameter is required"}).Validation()).Send(w)
			case errors.As(err, &invalid):
				fun.Error(fun.Err("invalid request parameter").WithFields(&fun.FieldError{Field: invalid.ParamName, Message: "invalid format"}).Validation()).Send(w)
			default:
				fun.InternalServerError("internal error").Send(w)
			}
		},
	})
}
