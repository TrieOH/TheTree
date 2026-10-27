package errs

import "errors"

// Request shapes the spec cannot express.
var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrInvalidDesignData  = errors.New("invalid design_data")
	ErrInvalidGalleryURLs = errors.New("invalid gallery_urls")
)

// Events, editions and templates.
var (
	ErrEventNotDraft                   = errors.New("cannot publish non draft event")
	ErrEventNotActive                  = errors.New("cannot discontinue non active event")
	ErrEditionAlreadyPublished         = errors.New("edition is already published")
	ErrInvalidBadgeOrigin              = errors.New("origin must be null or \"staff\"")
	ErrStaffTemplateWithTicketType     = errors.New("a staff badge template cannot target a ticket type")
	ErrSellerNotInPlatformWallet       = errors.New("seller does not belong to the platform payssage wallet")
	ErrInvalidPayssageWebhookSignature = errors.New("invalid webhook signature")
)
