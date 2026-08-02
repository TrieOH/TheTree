package app

import (
	"lib/httpserver"
	"univents/internal/openapi"
)

// adapt bridges the harness's strict middleware shape to the generated one.
// The generated handler func type is distinct from the harness's, but their
// underlying types are identical, so the leaf conversion is a plain cast.
func adapt(mw httpserver.StrictMiddlewareFunc) openapi.StrictMiddlewareFunc {
	return func(f openapi.StrictHandlerFunc, operationID string) openapi.StrictHandlerFunc {
		return openapi.StrictHandlerFunc(mw(httpserver.StrictHandlerFunc(f), operationID))
	}
}
