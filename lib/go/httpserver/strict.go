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
			return func(_ context.Context, w http.ResponseWriter, _ *http.Request, _ any) (any, error) {
				fun.InternalServerError("operation not registered with auth resolver").Send(w)
				return nil, nil
			}
		}
		if len(chain) == 0 {
			return f
		}
		return func(_ context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
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

// StrictHandlerFuncT and StrictMiddlewareFuncT capture the generated
// packages' strict types structurally: every backend's oapi-codegen output
// defines StrictHandlerFunc and StrictMiddlewareFunc with these exact
// underlying shapes (distinct types, one per package), so the harness
// converts to them with plain casts instead of per-backend adapt() funcs.
type StrictHandlerFuncT interface {
	~func(context.Context, http.ResponseWriter, *http.Request, any) (any, error)
}

type StrictMiddlewareFuncT[F StrictHandlerFuncT] interface {
	~func(F, string) F
}

// adaptTo converts the harness strict middleware into the generated
// package's strict-middleware type MF. The generated types share the
// harness's underlying shape, so the leaf conversion is a plain cast.
func adaptTo[HF StrictHandlerFuncT, MF StrictMiddlewareFuncT[HF]](mw StrictMiddlewareFunc) MF {
	return MF(func(f HF, operationID string) HF {
		return HF(mw(StrictHandlerFunc(f), operationID))
	})
}

// MountStrict mounts the backend's generated strict server with the
// harness's standard strict stack — request validation, then the
// spec-derived fail-closed auth dispatch — in that order, on every
// backend. The generated package's surface crosses the seam as function
// values and type parameters; every piece of policy (middleware order,
// fail-closed dispatch, unified param-binding error mapping) lives here,
// so the four backends cannot drift.
//
// The caller supplies the generated package's constructors as values
// (newStrict = openapi.NewStrictHandlerWithOptions, handlerWith =
// openapi.HandlerWithOptions) and builds the two generated options
// structs — strictOpts with the harness's StrictRequestErrorHandler and
// StrictResponseErrorHandler, chiOpts with BaseRouter, any extra
// ChiServerOptions middlewares, and ParamBindingErrorHandler. All type
// parameters except HF are inferred from the arguments.
func MountStrict[HF StrictHandlerFuncT, ServerT any, MF StrictMiddlewareFuncT[HF], SO, SI, CO any](
	server ServerT,
	chains map[string][]func(http.Handler) http.Handler,
	strictOpts SO,
	newStrict func(ServerT, []MF, SO) SI,
	chiOpts CO,
	handlerWith func(SI, CO) http.Handler,
) http.Handler {
	strict := newStrict(server, []MF{
		adaptTo[HF, MF](ValidateMiddleware()),
		adaptTo[HF, MF](AuthDispatch(chains)),
	}, strictOpts)
	return handlerWith(strict, chiOpts)
}

// ParamBindingErrorHandler returns the fun-envelope handler for the
// generated strict server's param-binding failures, unified across the
// backends to the superset: required params, invalid param formats, and
// required headers. The generated error classes are oapi-codegen's and
// identical across the four packages except for their package path, so the
// harness classifies by type name and reads ParamName by reflection — the
// generated surface is machine-owned; binding to it by name is the
// harness's one concession to the codegen.
func ParamBindingErrorHandler() func(w http.ResponseWriter, _ *http.Request, err error) {
	return func(w http.ResponseWriter, _ *http.Request, err error) {
		if resp, ok := paramBindingResponse(err); ok {
			resp.Send(w)
			return
		}
		fun.InternalServerError("internal error").Send(w)
	}
}

// paramBindingResponse classifies one param-binding error into its fun
// validation envelope. Only the three generated classes match; anything
// else falls through to the caller's default (500).
func paramBindingResponse(err error) (*fun.Response, bool) {
	t := reflect.TypeOf(err)
	if t == nil || t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		return nil, false
	}
	s := t.Elem()
	f, ok := s.FieldByName("ParamName")
	if !ok || f.Type.Kind() != reflect.String {
		return nil, false
	}
	name := reflect.ValueOf(err).Elem().FieldByName("ParamName").String()
	switch s.Name() {
	case "RequiredParamError":
		return fun.Error(fun.Err("invalid request parameter").WithFields(&fun.FieldError{Field: name, Message: "parameter is required"}).Validation()), true
	case "InvalidParamFormatError":
		return fun.Error(fun.Err("invalid request parameter").WithFields(&fun.FieldError{Field: name, Message: "invalid format"}).Validation()), true
	case "RequiredHeaderError":
		return fun.Error(fun.Err("invalid request header").WithFields(&fun.FieldError{Field: name, Message: "header is required"}).Validation()), true
	}
	return nil, false
}

// ValidateMiddleware returns a strict middleware that validates every
// request body before the handler runs. The generated request object's
// Body field is validated against its `validate` struct tags; operations
// without a body pass through. Validation errors are returned and flow
// through the strict server's ResponseErrorHandlerFunc.
//
// Non-struct bodies (maps, slices, json.RawMessage, scalars — generated
// from `type: object` with additionalProperties, `type: array`, or
// free-form schemas) carry no `validate` tags, so there is nothing to
// enforce: they pass through unvalidated. go-playground's Struct rejects
// non-struct values outright, so validating them would 400 every such
// request (e.g. provider webhook payloads like payssage's ReceiveWebhook,
// or informd's bulk-edit endpoints).
func ValidateMiddleware() StrictMiddlewareFunc {
	return func(f StrictHandlerFunc, _ string) StrictHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
			if body := bodyOf(request); body != nil && isStructish(body) {
				err := validator.Validate(body)
				if err != nil {
					return nil, err
				}
			}
			return f(ctx, w, r, request)
		}
	}
}

// isStructish reports whether v is a struct or a pointer to a struct — the
// only shapes go-playground's Struct validator can validate. Maps, slices,
// raw JSON and scalars have no fields with `validate` tags and are skipped.
func isStructish(v any) bool {
	t := reflect.TypeOf(v)
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t != nil && t.Kind() == reflect.Struct
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
// The underlying decode error is surfaced both in the message (so clients
// that only read the top-level message see why) and in the body field.
func StrictRequestErrorHandler() func(w http.ResponseWriter, _ *http.Request, err error) {
	return func(w http.ResponseWriter, _ *http.Request, err error) {
		fun.Error(fun.Err("invalid request body: " + err.Error()).WithFields(&fun.FieldError{Field: "body", Message: err.Error()}).BadRequest()).Send(w)
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
