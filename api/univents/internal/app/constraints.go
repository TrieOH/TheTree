package app

import "lib/database"

func SetupConstraintMessages() {
	database.SetConstraintErrorRegistry(database.ConstraintRegistry{
		// events
		"chk_event_status_valid":             "Event status must be one of: draft, active, discontinued.",
		"chk_event_payments_config_complete": "Both Payssage seller ID and wallet ID must be set together, or neither.",

		// event_members
		"chk_event_members_role_valid": "Member role must be one of: owner, admin, staff.",

		// editions
		"chk_editions_dates_valid":               "Edition end date must be after the start date.",
		"chk_editions_registration_before_start": "Registration opening date must be before or equal to the edition start date.",
		"excl_editions_no_overlap":               "This edition's dates overlap with another edition of the same event.",

		// registrations
		"chk_registrations_status_valid": "Registration status must be one of: pending, confirmed, cancelled, expired.",

		// products
		"uniq_products_edition_vendor_code":         "A product with this vendor code already exists in this edition.",
		"uniq_product_variants_edition_vendor_code": "A product variant with this vendor code already exists in this edition.",

		// product_purchases
		"chk_product_purchases_status_valid": "Product purchase status must be one of: pending, confirmed, cancelled, expired.",

		// programs
		"chk_programs_kind_valid": "Program kind must be one of: activity, checkpoint.",

		// program_occurrences
		"chk_program_occurrences_dates_valid": "Program occurrence end time must be after the start time.",

		// program_participations
		"chk_program_participations_status_valid": "Participation status must be one of: registered, attended, no_show, cancelled.",

		// signature_requests
		"chk_signature_requests_status_valid":     "Signature request status must be one of: pending, completed, expired, cancelled.",
		"uniq_signature_requests_idempotency_key": "This idempotency key has already been used.",

		// signatures
		"chk_signatures_status_valid":    "Signature status must be valid.",
		"chk_signatures_ready_has_image": "A ready signature must have an image URL.",

		// certification_templates
		"chk_certification_templates_kind_valid": "Certification template kind must be one of: edition_attendance, program_attendance.",
	})
}
