package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"slices"

	"lib/fielderr"
	"lib/telemetry"
	"lib/validator"

	"go.opentelemetry.io/otel/codes"
)

// The request binding's rejections, each naming the offending field
// (lib/fielderr). The Backend's mapper answers them (ADR-0019).
var (
	ErrInvalidParameter = errors.New("a path or query parameter is missing or malformed; errors lists which")
	ErrInvalidHeader    = errors.New("a required header is missing or malformed; errors lists which")
	// errOperationNotRegistered is spec/codegen drift: an operation Chains
	// resolved no chain for. It is a 500.
	errOperationNotRegistered = errors.New("operation has no auth chain")
)

// InvalidBodyError is a request body that does not decode. Its message
// carries the decode error, so a client reading only the detail sees why.
type InvalidBodyError struct{ Err error }

func (e *InvalidBodyError) Error() string { return "invalid request body: " + e.Err.Error() }

func (e *InvalidBodyError) Unwrap() error { return e.Err }

// StrictHandlerFunc mirrors the strict-server handler shape every backend's
// oapi-codegen output generates (openapi.StrictHandlerFunc). Each generated
// Mount adapts StrictStack to its own package's types (tools/strictmount).
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
			return func(_ context.Context, w http.ResponseWriter, r *http.Request, _ any) (any, error) {
				WriteError(w, r, fmt.Errorf("%s: %w", operationID, errOperationNotRegistered))
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

// TraceMiddleware returns a strict middleware that runs each operation
// inside one span named after its operationID, recording a returned error
// on it. Only an error answered with a 5xx marks the span as failed: a 4xx
// is the caller's mistake, not the service's (ADR-0016). It is the
// outermost strict middleware, so validation and auth show up inside the
// span. Repos need no spans of their own: the pgx tracer already records
// every query (lib/database).
func TraceMiddleware() StrictMiddlewareFunc {
	return func(f StrictHandlerFunc, operationID string) StrictHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
			ctx, span := telemetry.StartSpan(ctx, operationID)
			defer span.End()
			resp, err := f(ctx, w, r.WithContext(ctx), request)
			if err != nil {
				span.RecordError(err)
				if errorStatus(err) >= http.StatusInternalServerError {
					span.SetStatus(codes.Error, err.Error())
				}
			}
			return resp, err
		}
	}
}

// StrictStack is the harness's strict middleware stack, in the order the
// generated server takes it (it wraps in order, so the last runs first):
// one span per operation, around the operationID record, around the
// spec-derived fail-closed auth dispatch, around request validation. Every
// Backend's generated Mount (tools/strictmount) serves its strict server
// through this list, so the policy lives here once and the four Backends
// cannot drift.
func StrictStack(chains map[string][]func(http.Handler) http.Handler) []StrictMiddlewareFunc {
	return []StrictMiddlewareFunc{
		ValidateMiddleware(),
		AuthDispatch(chains),
		TraceMiddleware(),
		operationMiddleware(),
	}
}

// WriteParamError writes the generated strict server's param-binding
// failures (ChiServerOptions.ErrorHandlerFunc), unified across the backends
// to the superset: required params, invalid param formats, and required
// headers. The generated error classes are oapi-codegen's and identical
// across the four packages except for their package path, so the harness
// classifies by type name and reads ParamName by reflection — the generated
// surface is machine-owned; binding to it by name is the harness's one
// concession to the codegen.
func WriteParamError(w http.ResponseWriter, r *http.Request, err error) {
	bindErr := paramBindingError(err)
	if bindErr == nil {
		bindErr = err // an unknown class: a 500, its cause logged
	}
	WriteError(w, r, bindErr)
}

// paramBindingError classifies one param-binding error into its named
// error, naming the parameter. Only the three generated classes match;
// anything else is nil, and the caller's default (500).
func paramBindingError(err error) error {
	t := reflect.TypeOf(err)
	if t == nil || t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		return nil
	}
	s := t.Elem()
	f, ok := s.FieldByName("ParamName")
	if !ok || f.Type.Kind() != reflect.String {
		return nil
	}
	name := reflect.ValueOf(err).Elem().FieldByName("ParamName").String()
	switch s.Name() {
	case "RequiredParamError":
		return fielderr.On(ErrInvalidParameter, name, "parameter is required")
	case "InvalidParamFormatError":
		return fielderr.On(ErrInvalidParameter, name, "invalid format")
	case "RequiredHeaderError":
		return fielderr.On(ErrInvalidHeader, name, "header is required")
	}
	return nil
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

// WriteBodyError writes a request body that does not decode
// (StrictHTTPServerOptions.RequestErrorHandlerFunc), naming the body field.
// Errors returned by strict handlers go straight to WriteError
// (ResponseErrorHandlerFunc).
func WriteBodyError(w http.ResponseWriter, r *http.Request, err error) {
	WriteError(w, r, fielderr.On(&InvalidBodyError{Err: err}, "body", err.Error()))
}
