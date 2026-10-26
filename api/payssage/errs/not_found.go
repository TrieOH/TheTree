package errs

import "errors"

// Each resource's not-found error. Repos return them as
// lib/database.NotFoundError, so a caller can match the resource or any
// missing row.
var (
	ErrOAuthStateNotFound      = errors.New("state not found")
	ErrWebhookEventNotFound    = errors.New("webhook event not found")
	ErrWalletNotFound          = errors.New("wallet not found")
	ErrWebhookEndpointNotFound = errors.New("webhook endpoint not found")
	ErrIntentNotFound          = errors.New("intent not found")
	ErrOrganizationNotFound    = errors.New("organization not found")
	ErrCollectorNotFound       = errors.New("collector not found")
	ErrWebhookDeliveryNotFound = errors.New("webhook delivery not found")
	ErrSellerNotFound          = errors.New("seller not found")
)
