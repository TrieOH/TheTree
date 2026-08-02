package app

import (
	"errors"
	"net/http"

	spec "Informd"
	"Informd/internal/handlers"
	"Informd/internal/openapi"
	"lib/errx"
	"lib/httpserver"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
)

func (app *Informd) CreateRouter(handlers *handlers.Server, middlewares middlewares) http.Handler {
	chains, err := resolveAuthChains(middlewares)
	errx.Exit(err, "resolve auth chains")
	return httpserver.NewRouter(httpserver.Config{
		AppName:            app.cfg.AppName,
		CorsAllowedOrigins: app.cfg.CorsAllowedOrigins,
		CorsAllowedHeaders: app.cfg.CorsAllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
		Routes: func(r *chi.Mux) {
			mountStrict(r, handlers, chains)
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
