package httpserver

import (
	"context"
	"net/http"
	"reflect"
	"slices"

	"lib/validator"

	"github.com/MintzyG/fun"
)

// StrictHandlerFunc mirrors the strict-server handler shape every backend's
// oapi-codegen output generates (method = openapi.StrictHandlerFunc). The
// generated types share this underlying shape, so backends convert at the
// seam:
//
//	openapi.StrictMiddlewareFunc(httpserver.AuthDispatch(chains))
//
// instead of re-implementing the dispatch loop per backend.
type StrictHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error)

// StrictMiddlewareFunc mirrors the strict-server middleware shape from the
// same generated output (method = openapi.StrictMiddlewareFunc).
type StrictMiddlewareFunc func(f StrictHandlerFunc, operationID string) StrictHandlerFunc

// AuthDispatch returns a strict middleware that resolves each operation's
// auth chain by generated-form operationID — the exact value the generated
// server passes, so the lookup needs no string surgery — and runs it around
// the handler. Public operations (empty chain) pass through untouched; an
// operationID with no chain is spec/codegen drift and fails closed with a
// 500 instead of being treated as public.
func AuthDispatch(chains map[string][]func(http.Handler) http.Handler) StrictMiddlewareFunc {
	return func(f StrictHandlerFunc, operationID string) StrictHandlerFunc {
		chain, ok := chains[operationID]
		if !ok {
			return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
				fun.InternalServerError("operation not registered with auth resolver").Send(w)
				return nil, nil
			}
		}
		if len(chain) == 0 {
			return f
		}
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
			var resp any
			var ferr error
			var called bool
			var next http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				resp, ferr = f(r.Context(), w, r, request)
			})
			for i := range slices.Backward(chain) {
				next = chain[i](next)
			}
			next.ServeHTTP(w, r)
			if !called {
				// An auth middleware rejected the request and already
				// wrote the response; nothing to return.
				return nil, nil
			}
			return resp, ferr
		}
	}
}

// ValidateMiddleware returns a strict middleware that validates every
// request body before the handler runs. The generated request object's
// Body field is validated against its `validate` struct tags; operations
// without a body pass through. Validation errors are returned and flow
// through the strict server's ResponseErrorHandlerFunc.
func ValidateMiddleware() StrictMiddlewareFunc {
	return func(f StrictHandlerFunc, _ string) StrictHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
			if body := bodyOf(request); body != nil {
				err := validator.Validate(body)
				if err != nil {
					return nil, err
				}
			}
			return f(ctx, w, r, request)
		}
	}
}

// bodyOf extracts the generated request object's Body field, if any.
func bodyOf(request any) any {
	v := reflect.ValueOf(request)
	if v.Kind() != reflect.Struct {
		return nil
	}
	b := v.FieldByName("Body")
	if !b.IsValid() || b.IsNil() {
		return nil
	}
	return b.Interface()
}

// StrictRequestErrorHandler returns the fun-envelope handler for request
// binding failures, for StrictHTTPServerOptions.RequestErrorHandlerFunc.
func StrictRequestErrorHandler() func(w http.ResponseWriter, _ *http.Request, err error) {
	return func(w http.ResponseWriter, _ *http.Request, err error) {
		fun.Error(fun.Err("invalid request body").WithFields(&fun.FieldError{Field: "body", Message: err.Error()}).BadRequest()).Send(w)
	}
}

// StrictResponseErrorHandler returns the fun-envelope handler for errors
// returned by strict handlers, for
// StrictHTTPServerOptions.ResponseErrorHandlerFunc.
func StrictResponseErrorHandler() func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		fun.Error(err).SendWithCtx(r.Context(), w)
	}
}
