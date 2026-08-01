package app

import (
	"errors"
	"net/http"

	"lib/httpserver"
	spec "univents"
	"univents/internal/handlers"
	"univents/internal/openapi"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
	"riverqueue.com/riverui"
)

func (app *Univents) CreateRouter(middlewares middlewares, h *handlers.Server, riverUIHandler *riverui.Handler) http.Handler {
	return httpserver.NewRouter(httpserver.Config{
		AppName:         app.cfg.AppName,
		OpenAPISpec:     spec.OpenAPISpec,
		SkipLogPrefixes: []string{"/admin/asynq"},
		Routes: func(r *chi.Mux) {
			mountStrict(r, h, middlewares)

			r.Group(func(r chi.Router) {
				r.Mount("/riverui", riverUIHandler)
			})
		},
	})
}

// mountStrict registers the generated strict handler on r with the
// validation + auth middleware stack and the fun-envelope error handlers.
func mountStrict(r *chi.Mux, h *handlers.Server, mw middlewares) {
	strict := openapi.NewStrictHandlerWithOptions(h,
		[]openapi.StrictMiddlewareFunc{handlers.ValidateMiddleware(), authDispatch(mw)},
		openapi.StrictHTTPServerOptions{
			RequestErrorHandlerFunc: func(w http.ResponseWriter, _ *http.Request, err error) {
				fun.Error(fun.Err("invalid request body").WithFields(&fun.FieldError{Field: "body", Message: err.Error()}).BadRequest()).Send(w)
			},
			ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
				fun.Error(err).SendWithCtx(r.Context(), w)
			},
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
