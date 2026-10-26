package app

import (
	"lib/database"
	"univents/errs"
)

// Constraints maps every named constraint in the schema to the named error
// repos surface on violation (ADR-0016); the postgres module takes it at
// boot, and boot fails when the schema has a constraint missing here.
func Constraints() database.ConstraintErrors {
	return database.ConstraintErrors{
		// events
		"chk_event_status_valid":                        errs.ErrInvalidEventStatus,
		"chk_event_payments_public_key_requires_seller": errs.ErrPublicKeyNeedsSeller,

		// event_members
		"chk_event_members_role_valid": errs.ErrInvalidEventMemberRole,

		// editions
		"chk_editions_dates_valid":               errs.ErrEditionDates,
		"chk_editions_registration_before_start": errs.ErrRegistrationAfterStart,
		"excl_editions_no_overlap":               errs.ErrEditionOverlap,

		// registrations
		"chk_registrations_status_valid":                 errs.ErrInvalidRegistrationStatus,
		"uniq_registrations_active_per_edition_attendee": errs.ErrTicketHeld,
		"uniq_registrations_active_email_per_edition":    errs.ErrTicketHeld,

		// products
		"uniq_products_edition_vendor_code":         errs.ErrProductCodeTaken,
		"uniq_product_variants_edition_vendor_code": errs.ErrVariantCodeTaken,

		// product_purchases
		"chk_product_purchases_status_valid": errs.ErrInvalidProductPurchaseStatus,

		// programs
		"chk_programs_kind_valid": errs.ErrInvalidProgramKind,

		// program_occurrences
		"chk_program_occurrences_dates_valid": errs.ErrOccurrenceDates,

		// program_participations
		"chk_program_participations_status_valid":                    errs.ErrInvalidParticipationStatus,
		"uniq_program_participations_active_per_occurrence_attendee": errs.ErrOccurrenceRegistrationTaken,

		// signature_requests
		"chk_signature_requests_status_valid":     errs.ErrInvalidSignatureRequestStatus,
		"uniq_signature_requests_idempotency_key": errs.ErrIdempotencyKeyUsed,

		// signatures
		"chk_signatures_status_valid":    errs.ErrInvalidSignatureStatus,
		"chk_signatures_ready_has_image": errs.ErrReadySignatureNeedsImage,

		// certification_templates
		"chk_certification_templates_kind_valid": errs.ErrInvalidCertificationKind,

		// badge_emissions
		"chk_badge_emissions_origin_valid":         errs.ErrInvalidBadgeEmissionOrigin,
		"chk_badge_emissions_status_valid":         errs.ErrInvalidBadgeEmissionStatus,
		"uniq_badge_emissions_edition_user_origin": errs.ErrBadgeEmitted,

		// badge_templates
		"chk_badge_template_scope_valid":   errs.ErrInvalidBadgeTemplateScope,
		"idx_one_badge_template_per_scope": errs.ErrBadgeTemplateScopeTaken,

		// purchases (store)
		"chk_purchases_status_valid":                   errs.ErrInvalidPurchaseStatus,
		"uniq_purchases_pending_per_purchaser_edition": errs.ErrPendingPurchaseExists,

		// purchase_items
		"chk_purchase_items_type_valid":        errs.ErrInvalidPurchaseItemType,
		"chk_purchase_items_quantity_positive": errs.ErrPurchaseItemQuantity,
		"uniq_purchase_items_ticket_unit":      errs.ErrTicketUnitInPurchase,
		"uniq_purchase_items_product_line":     errs.ErrProductLineInPurchase,
		"uniq_purchase_items_program_line":     errs.ErrProgramLineInPurchase,

		// ws_tokens
		"uniq_ws_tokens_token_hash": errs.ErrWSTokenIssued,
	}
}
