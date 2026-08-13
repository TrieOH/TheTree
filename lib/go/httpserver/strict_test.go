package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"lib/validator"
)

type adKey string

const adSentinelKey adKey = "auth-dispatch-sentinel"

var adSentinel = &struct{}{}

// TestAuthDispatchRunsChainForProtectedOperation — the generated server
// passes the operationID in generated form; the chain must run, and the
// handler must observe the middleware-modified request context (auth
// middlewares replace the request via r.WithContext).
func TestAuthDispatchRunsChainForProtectedOperation(t *testing.T) {
	var ran bool
	jwt := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ran = true
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), adSentinelKey, adSentinel)))
		})
	}
	dispatch := AuthDispatch(map[string][]func(http.Handler) http.Handler{
		"PostLogout": {jwt},
	})

	handler := dispatch(func(ctx context.Context, _ http.ResponseWriter, _ *http.Request, _ any) (any, error) {
		if got := ctx.Value(adSentinelKey); got != adSentinel {
			return nil, fmt.Errorf("handler saw stale context: sentinel %v, want %v", got, adSentinel)
		}
		return "ok", nil
	}, "PostLogout") // exactly as oapi-codegen passes it

	resp, err := handler(context.Background(), httptest.NewRecorder(), httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/auth/logout", nil), nil)
	if err != nil {
		t.Fatalf("handler must see the middleware-modified request context: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("want %q, got %v", "ok", resp)
	}
	if !ran {
		t.Fatal("auth middleware never ran: operationID mismatch")
	}
}

func TestAuthDispatchSkipsPublicOperation(t *testing.T) {
	var ran bool
	dispatch := AuthDispatch(map[string][]func(http.Handler) http.Handler{
		"GetJWKS": nil, // public operation — empty chain
	})
	handler := dispatch(func(_ context.Context, _ http.ResponseWriter, _ *http.Request, _ any) (any, error) {
		ran = true
		return "ok", nil
	}, "GetJWKS")
	_, err := handler(context.Background(), httptest.NewRecorder(), httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/.well-known/jwks.json", nil), nil)
	if err != nil {
		t.Fatalf("public operation must pass through: %v", err)
	}
	if !ran {
		t.Fatal("public operation handler must run")
	}
}

func TestAuthDispatchRejectionShortCircuits(t *testing.T) {
	dispatch := AuthDispatch(map[string][]func(http.Handler) http.Handler{
		"PostLogout": {func(_ http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
			})
		}},
	})
	var called bool
	handler := dispatch(func(_ context.Context, _ http.ResponseWriter, _ *http.Request, _ any) (any, error) {
		called = true
		return "ok", nil
	}, "PostLogout")
	rec := httptest.NewRecorder()
	resp, err := handler(context.Background(), rec, httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/auth/logout", nil), nil)
	if err != nil {
		t.Fatalf("rejected request must return nil error: %v", err)
	}
	if resp != nil {
		t.Fatalf("rejected request must return nil response, got %v", resp)
	}
	if called {
		t.Fatal("handler must not run after auth rejection")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d", rec.Code)
	}
}

// Fail-closed: an operationID the resolver does not know (spec/codegen
// drift) must 500, never be treated as public.
func TestAuthDispatchUnknownOperationFailsClosed(t *testing.T) {
	dispatch := AuthDispatch(map[string][]func(http.Handler) http.Handler{
		"PostLogout": {func(next http.Handler) http.Handler { return next }},
	})
	var called bool
	handler := dispatch(func(_ context.Context, _ http.ResponseWriter, _ *http.Request, _ any) (any, error) {
		called = true
		return "ok", nil
	}, "RenamedOperation") // in the generated server but not in the spec
	rec := httptest.NewRecorder()
	resp, err := handler(context.Background(), rec, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil), nil)
	if err != nil {
		t.Fatalf("unknown operation must write a response, not error out: %v", err)
	}
	if resp != nil {
		t.Fatalf("unknown operation must return nil response, got %v", resp)
	}
	if called {
		t.Fatal("handler must not run for an unknown operation")
	}
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("want 500, got %d", rec.Code)
	}
}

type bodyCarrier struct {
	Body *validatedBody
}

// mapBodyCarrier mirrors generated request objects whose body schema is
// `type: object, additionalProperties: true` (e.g. provider webhook payloads)
// — the body is a map, which has no `validate` tags to enforce.
type mapBodyCarrier struct {
	Body *map[string]interface{}
}

type validatedBody struct {
	Name string `validate:"required"`
}

func TestValidateMiddlewareValidatesBody(t *testing.T) {
	validator.SetupValidator()
	validate := ValidateMiddleware()

	handler := validate(func(_ context.Context, _ http.ResponseWriter, _ *http.Request, _ any) (any, error) {
		return "ok", nil
	}, "CreateThing")

	// valid body passes through (generated request objects are passed by value)
	_, err := handler(context.Background(), httptest.NewRecorder(), httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", nil),
		bodyCarrier{Body: &validatedBody{Name: "x"}})
	if err != nil {
		t.Fatalf("valid body must pass validation: %v", err)
	}

	// invalid body returns the validation error
	_, err = handler(context.Background(), httptest.NewRecorder(), httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", nil),
		bodyCarrier{Body: &validatedBody{}})
	if err == nil {
		t.Fatal("invalid body must fail validation")
	}
}

func TestValidateMiddlewareSkipsNonStructBodies(t *testing.T) {
	validator.SetupValidator()
	validate := ValidateMiddleware()

	for _, name := range []string{"map", "slice"} {
		t.Run(name, func(t *testing.T) {
			var called bool
			handler := validate(func(_ context.Context, _ http.ResponseWriter, _ *http.Request, _ any) (any, error) {
				called = true
				return "ok", nil
			}, "ReceiveWebhook")

			_, err := handler(context.Background(), httptest.NewRecorder(), httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/webhooks/x", nil),
				mapBodyCarrier{Body: &map[string]interface{}{}})
			if err != nil {
				t.Fatalf("%s body must pass through unvalidated: %v", name, err)
			}
			if !called {
				t.Fatal("handler must run for non-struct bodies")
			}
		})
	}

	// struct bodies still validate
	handler := validate(func(_ context.Context, _ http.ResponseWriter, _ *http.Request, _ any) (any, error) {
		return "ok", nil
	}, "CreateThing")
	_, err := handler(context.Background(), httptest.NewRecorder(), httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", nil),
		bodyCarrier{Body: &validatedBody{}}) // Name required but empty
	if err == nil {
		t.Fatal("struct body must still fail validation")
	}
}

func TestValidateMiddlewareSkipsBodylessRequest(t *testing.T) {
	validator.SetupValidator()
	validate := ValidateMiddleware()

	var called bool
	handler := validate(func(_ context.Context, _ http.ResponseWriter, _ *http.Request, _ any) (any, error) {
		called = true
		return "ok", nil
	}, "GetThing")

	// request object without a Body field passes through untouched
	_, err := handler(context.Background(), httptest.NewRecorder(), httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil),
		&struct{ Params any }{Params: nil})
	if err != nil {
		t.Fatalf("bodyless request must pass through: %v", err)
	}
	if !called {
		t.Fatal("handler must run for bodyless requests")
	}
}
