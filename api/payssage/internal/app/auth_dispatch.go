package app

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"payssage/internal/openapi"
)

// authChains maps every operationId to the middleware chain it must run
// through, mirroring the parity-test matrix. Everything is JWT-protected
// except the provider webhook receive and the OAuth callback (browser
// redirect target), which are public.
func authChains(mw middlewares) map[string][]func(http.Handler) http.Handler {
	chains := map[string][]func(http.Handler) http.Handler{
		"providerCallback": {},
		"receiveWebhook":   {},

		"listOrganizations":            {mw.jwtAuth},
		"createOrganization":           {mw.jwtAuth},
		"listOrganizationMembers":      {mw.jwtAuth},
		"addOrganizationMember":        {mw.jwtAuth},
		"removeOrganizationMember":     {mw.jwtAuth},
		"getOrganizationMemberByID":    {mw.jwtAuth},
		"getOrganizationMemberByEmail": {mw.jwtAuth},

		"createWallet":            {mw.jwtAuth},
		"listWallets":             {mw.jwtAuth},
		"getWallet":               {mw.jwtAuth},
		"setWalletFee":            {mw.jwtAuth},
		"setWalletSandbox":        {mw.jwtAuth},
		"listOrganizationWallets": {mw.jwtAuth},
		"bindCollector":           {mw.jwtAuth},
		"unbindCollector":         {mw.jwtAuth},

		"listCollectors":             {mw.jwtAuth},
		"getCollector":               {mw.jwtAuth},
		"listOrganizationCollectors": {mw.jwtAuth},

		"listWalletSellers": {mw.jwtAuth},

		"listIntentsByProfile":    {mw.jwtAuth},
		"getIntent":               {mw.jwtAuth},
		"cancelIntent":            {mw.jwtAuth},
		"listWalletIntents":       {mw.jwtAuth},
		"listOrganizationIntents": {mw.jwtAuth},
		"checkout":                {mw.jwtAuth},
		"hardCreateIntent":        {mw.jwtAuth},

		"connectProvider": {mw.jwtAuth},
		"revokeProvider":  {mw.jwtAuth},

		"createWebhookEndpoint": {mw.jwtAuth},
		"listWebhookEndpoints":  {mw.jwtAuth},
		"getWebhookEndpoint":    {mw.jwtAuth},
		"deleteWebhookEndpoint": {mw.jwtAuth},

		"listWebhookEvents": {mw.jwtAuth},
		"getWebhookEvent":   {mw.jwtAuth},

		"listWebhookDeliveries": {mw.jwtAuth},
		"getWebhookDelivery":    {mw.jwtAuth},
	}
	return chains
}

// authDispatch is the strict-server middleware that resolves the auth
// chain for each operation (by operationId) and runs it around the
// handler. Public operations (empty chain) pass through untouched.
func authDispatch(mw middlewares) openapi.StrictMiddlewareFunc {
	chains := authChains(mw)
	return func(f openapi.StrictHandlerFunc, operationID string) openapi.StrictHandlerFunc {
		if operationID != "" {
			operationID = strings.ToLower(operationID[:1]) + operationID[1:]
		}
		chain := chains[operationID]
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
