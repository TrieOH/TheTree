package app

import (
	"errors"
	"net/http"

	spec "IdentityX"
	"IdentityX/internal/handlers"
	"IdentityX/internal/openapi"
	"lib/errx"
	"lib/httpserver"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
)

func (app *IdentityX) CreateRouter(middlewares middlewares, h *handlers.Server) http.Handler {
	chains, err := resolveAuthChains(middlewares)
	errx.Exit(err, "resolve auth chains")
	return httpserver.NewRouter(httpserver.Config{
		AppName:     app.cfg.AppName,
		OpenAPISpec: spec.OpenAPISpec,
		Routes: func(r *chi.Mux) {
			mountStrict(r, h, chains)
		},
	})
}

// mountStrict registers the generated strict handler on r with the
// validation + auth middleware stack and the fun-envelope error handlers.
func mountStrict(r *chi.Mux, h *handlers.Server, chains map[string][]func(http.Handler) http.Handler) {
	strict := openapi.NewStrictHandlerWithOptions(h,
		[]openapi.StrictMiddlewareFunc{handlers.ValidateMiddleware(), authDispatch(chains)},
		openapi.StrictHTTPServerOptions{
			RequestErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
				// body decode failures — known client error
				fun.Error(fun.Err("invalid request body").WithFields(&fun.FieldError{Field: "body", Message: err.Error()}).BadRequest()).Send(w)
			},
			ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				// *fun.AppError carries its status (code -> HTTP);
				// anything unknown resolves to 500.
				fun.Error(err).SendWithCtx(r.Context(), w)
			},
		})
	openapi.HandlerWithOptions(strict, openapi.ChiServerOptions{
		BaseRouter: r,
		// param binding failures from the generated wrapper
		ErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
			var required *openapi.RequiredParamError
			var invalid *openapi.InvalidParamFormatError
			var requiredHeader *openapi.RequiredHeaderError
			switch {
			case errors.As(err, &required):
				fun.Error(fun.Err("invalid request parameter").WithFields(&fun.FieldError{Field: required.ParamName, Message: "parameter is required"}).Validation()).Send(w)
			case errors.As(err, &invalid):
				fun.Error(fun.Err("invalid request parameter").WithFields(&fun.FieldError{Field: invalid.ParamName, Message: "invalid format"}).Validation()).Send(w)
			case errors.As(err, &requiredHeader):
				fun.Error(fun.Err("invalid request header").WithFields(&fun.FieldError{Field: requiredHeader.ParamName, Message: "header is required"}).Validation()).Send(w)
			default:
				fun.InternalServerError("internal error").Send(w)
			}
		},
	})
}
