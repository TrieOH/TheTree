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

	// registrations
	ErrInvalidRegistrationStatus = errors.New("registration status must be one of: pending, confirmed, cancelled, expired")

	// products, product_purchases
	ErrInvalidProductPurchaseStatus = errors.New("product purchase status must be one of: pending, confirmed, cancelled, expired")

	// programs, occurrences, participations
	ErrInvalidProgramKind         = errors.New("program kind must be one of: activity, checkpoint")
	ErrOccurrenceDates            = errors.New("program occurrence end time must be after the start time")
	ErrInvalidParticipationStatus = errors.New("participation status must be one of: registered, attended, no_show, cancelled")

	// signature_requests, signatures
	ErrInvalidSignatureRequestStatus = errors.New("signature request status must be one of: pending, completed, expired, cancelled")
	ErrInvalidSignatureStatus        = errors.New("signature status must be valid")
	ErrReadySignatureNeedsImage      = errors.New("a ready signature must have an image URL")

	// certification_templates
	ErrInvalidCertificationKind = errors.New("certification template kind must be one of: edition_attendance, program_attendance")

	// badge_emissions, badge_templates
	ErrInvalidBadgeEmissionOrigin = errors.New("badge emission origin must be one of: participant, staff")
	ErrInvalidBadgeEmissionStatus = errors.New("badge emission status must be one of: active, revoked")
	ErrInvalidBadgeTemplateScope  = errors.New("badge template scope must be the edition default (null) or staff (staff, without a ticket type)")

	// purchases, purchase_items
	ErrInvalidPurchaseStatus   = errors.New("purchase status must be one of: pending, approved, expired, cancelled, failed, rejected, refunded")
	ErrInvalidPurchaseItemType = errors.New("purchase item type must be one of: ticket, product, program_occurrence")
	ErrPurchaseItemQuantity    = errors.New("purchase item quantity must be greater than zero")
)
