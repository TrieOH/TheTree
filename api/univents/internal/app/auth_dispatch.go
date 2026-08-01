package app

import (
	"context"
	"net/http"
	"slices"

	"univents/internal/openapi"
)

// authChains maps every operationId to the middleware chain it must run
// through, mirroring the parity-test matrix. Public read routes plus the
// signature fulfill/deny/revoke token-authenticated routes have no chain;
// everything else requires a JWT.
func authChains(mw middlewares) map[string][]func(http.Handler) http.Handler {
	chains := map[string][]func(http.Handler) http.Handler{
		"addEventMember":                  {mw.jwt},
		"cancelSignatureRequest":          {mw.jwt},
		"createBadgeTemplate":             {mw.jwt},
		"createCertificationTemplate":     {mw.jwt},
		"createEdition":                   {mw.jwt},
		"createEvent":                     {mw.jwt},
		"createInitialProduct":            {mw.jwt},
		"createProductVariant":            {mw.jwt},
		"createProgram":                   {mw.jwt},
		"createProgramOccurrence":         {mw.jwt},
		"createSignature":                 {mw.jwt},
		"createSignatureRequest":          {mw.jwt},
		"createTicketType":                {mw.jwt},
		"deleteBadgeTemplate":             {mw.jwt},
		"deleteCertificationTemplate":     {mw.jwt},
		"deleteOccurrence":                {mw.jwt},
		"deleteProduct":                   {mw.jwt},
		"deleteProductVariant":            {mw.jwt},
		"deleteProgram":                   {mw.jwt},
		"deleteSignature":                 {mw.jwt},
		"denySignatureRequest":            {},
		"discontinueEvent":                {mw.jwt},
		"fulfillSignatureRequest":         {},
		"getActiveEdition":                {},
		"getBadgeTemplate":                {mw.jwt},
		"getCertification":                {mw.jwt},
		"getCertificationTemplate":        {},
		"getEditionBySlug":                {},
		"getEventBySlug":                  {},
		"getHealth":                       {},
		"getOccurrence":                   {},
		"getOpenAPISpec":                  {},
		"getProduct":                      {},
		"getProductByVendorCode":          {},
		"getProgram":                      {},
		"getSignature":                    {},
		"getSignatureRequest":             {},
		"getTicketType":                   {},
		"getVariantByVendorCode":          {},
		"invalidateCertification":         {mw.jwt},
		"linkCertificationTemplate":       {mw.jwt},
		"listBadgeTemplates":              {mw.jwt},
		"listCertificationEmissionErrors": {mw.jwt},
		"listCertificationTemplateLinks":  {},
		"listCertificationTemplates":      {},
		"listDraftEditions":               {mw.jwt},
		"listEditionCertifications":       {mw.jwt},
		"listEditionOccurrences":          {},
		"listEditionProducts":             {},
		"listEditionPrograms":             {},
		"listEditionSignatureRequests":    {},
		"listEditionSignatures":           {},
		"listEventMembers":                {mw.jwt},
		"listJoinedEvents":                {mw.jwt},
		"listMyCertifications":            {mw.jwt},
		"listOwnedEvents":                 {mw.jwt},
		"listPastEditions":                {},
		"listProductVariants":             {},
		"listProgramOccurrences":          {},
		"listPublicEditions":              {},
		"listPublicEvents":                {},
		"listTicketTypes":                 {},
		"listUpcomingEditions":            {},
		"patchEdition":                    {mw.jwt},
		"patchEvent":                      {mw.jwt},
		"patchOccurrence":                 {mw.jwt},
		"patchProduct":                    {mw.jwt},
		"patchProductVariant":             {mw.jwt},
		"patchProgram":                    {mw.jwt},
		"patchTicketType":                 {mw.jwt},
		"publishEdition":                  {mw.jwt},
		"publishEvent":                    {mw.jwt},
		"removeEventMember":               {mw.jwt},
		"revokeSignature":                 {},
		"unlinkCertificationTemplate":     {mw.jwt},
		"updateCertificationTemplate":     {mw.jwt},
		"verifyCertification":             {},
	}
	return chains
}

// authDispatch is the strict-server middleware that resolves the auth
// chain for each operation (by operationId) and runs it around the
// handler. Public operations (empty chain) pass through untouched.
func authDispatch(mw middlewares) openapi.StrictMiddlewareFunc {
	chains := authChains(mw)
	return func(f openapi.StrictHandlerFunc, operationID string) openapi.StrictHandlerFunc {
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
				resp, ferr = f(ctx, w, r, request)
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
