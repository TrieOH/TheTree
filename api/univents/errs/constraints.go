package errs

import "errors"

// One error per named constraint in the schema. internal/app's constraint
// registry maps each constraint to its error, and lib/database pairs it
// with the violation's class. Boot fails when the schema has a constraint
// the registry does not list.
var (
	// events, event_members
	ErrInvalidEventStatus     = errors.New("event status must be one of: draft, active, discontinued")
	ErrPublicKeyNeedsSeller   = errors.New("an event's Payssage public key requires a connected seller")
	ErrInvalidEventMemberRole = errors.New("member role must be one of: owner, admin, staff")

	// editions
	ErrEditionDates           = errors.New("edition end date must be after the start date")
	ErrRegistrationAfterStart = errors.New("registration opening date must be before or equal to the edition start date")
	ErrEditionOverlap         = errors.New("this edition's dates overlap with another edition of the same event")

	// registrations
	ErrInvalidRegistrationStatus = errors.New("registration status must be one of: pending, confirmed, cancelled, expired")
	ErrTicketHeld                = errors.New("this person already has a ticket for this edition")

	// products, product_purchases
	ErrProductCodeTaken             = errors.New("a product with this vendor code already exists in this edition")
	ErrVariantCodeTaken             = errors.New("a product variant with this vendor code already exists in this edition")
	ErrInvalidProductPurchaseStatus = errors.New("product purchase status must be one of: pending, confirmed, cancelled, expired")

	// programs, occurrences, participations
	ErrInvalidProgramKind          = errors.New("program kind must be one of: activity, checkpoint")
	ErrOccurrenceDates             = errors.New("program occurrence end time must be after the start time")
	ErrInvalidParticipationStatus  = errors.New("participation status must be one of: registered, attended, no_show, cancelled")
	ErrOccurrenceRegistrationTaken = errors.New("this person is already registered for this activity occurrence")

	// signature_requests, signatures
	ErrInvalidSignatureRequestStatus = errors.New("signature request status must be one of: pending, completed, expired, cancelled")
	ErrIdempotencyKeyUsed            = errors.New("this idempotency key has already been used")
	ErrInvalidSignatureStatus        = errors.New("signature status must be valid")
	ErrReadySignatureNeedsImage      = errors.New("a ready signature must have an image URL")

	// certification_templates
	ErrInvalidCertificationKind = errors.New("certification template kind must be one of: edition_attendance, program_attendance")

	// badge_emissions, badge_templates
	ErrInvalidBadgeEmissionOrigin = errors.New("badge emission origin must be one of: participant, staff")
	ErrInvalidBadgeEmissionStatus = errors.New("badge emission status must be one of: active, revoked")
	ErrBadgeEmitted               = errors.New("this person already has a badge for this edition and origin")
	ErrInvalidBadgeTemplateScope  = errors.New("badge template scope must be the edition default (null) or staff (staff, without a ticket type)")
	ErrBadgeTemplateScopeTaken    = errors.New("this edition already has a template for this scope (default, ticket type, or staff)")

	// purchases, purchase_items
	ErrInvalidPurchaseStatus   = errors.New("purchase status must be one of: pending, approved, expired, cancelled, failed, rejected, refunded")
	ErrPendingPurchaseExists   = errors.New("this person already has a pending purchase for this edition")
	ErrInvalidPurchaseItemType = errors.New("purchase item type must be one of: ticket, product, program_occurrence")
	ErrPurchaseItemQuantity    = errors.New("purchase item quantity must be greater than zero")
	ErrTicketUnitInPurchase    = errors.New("this ticket is already in this purchase for this person")
	ErrProductLineInPurchase   = errors.New("this product is already in this purchase")
	ErrProgramLineInPurchase   = errors.New("this program occurrence is already in this purchase")

	// ws_tokens
	ErrWSTokenIssued = errors.New("this webhook token has already been issued")
)
