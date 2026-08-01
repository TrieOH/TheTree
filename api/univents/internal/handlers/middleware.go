package handlers

import (
	"context"
	"net/http"
	"reflect"

	"lib/validator"
	"univents/internal/openapi"
)

// ValidateMiddleware returns a strict-server middleware that validates
// every request body before the handler runs. The generated request
// object's Body field is validated against its `validate` struct tags;
// operations without a body pass through. Validation errors are returned
// as *fun.AppError and flow through the strict server's
// ResponseErrorHandlerFunc.
func ValidateMiddleware() openapi.StrictMiddlewareFunc {
	return func(f openapi.StrictHandlerFunc, _ string) openapi.StrictHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
			body := bodyOf(request)
			if body != nil {
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
