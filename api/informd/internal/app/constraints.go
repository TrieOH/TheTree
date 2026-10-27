package app

import (
	"Informd/errs"
	"lib/database"
)

// Constraints maps Informd's constraint and unique-index names to the named
// error repos surface on violation (ADR-0016). They travel into
// database.SetupDB as data — boot fails loudly when a migration adds a
// constraint missing here.
func Constraints() database.ConstraintErrors {
	return database.ConstraintErrors{
		// forms
		"chk_forms_valid_status":       errs.ErrInvalidFormStatus,
		"chk_forms_valid_status_state": errs.ErrInvalidFormStatusState,
		"uniq_form_name_per_namespace": errs.ErrFormNameTaken,
		"uniq_name_per_user":           errs.ErrAPIKeyNameTaken,

		// fields
		"chk_fields_type":       errs.ErrInvalidFieldType,
		"chk_fields_key_format": errs.ErrInvalidFieldKeyFormat,
		"uniq_key_per_step":     errs.ErrFieldKeyTaken,

		// select
		"chk_select_behaviour":  errs.ErrInvalidSelectBehavior,
		"chk_select_options":    errs.ErrEmptySelectOptions,
		"chk_select_value_type": errs.ErrInvalidSelectType,

		// namespaces
		"uniq_namespace_name_per_user":    errs.ErrNamespaceNameTaken,
		"chk_valid_namespace_member_role": errs.ErrInvalidNamespaceMemberRole,

		// response
		"uniq_answer_per_field_per_response": errs.ErrFieldAlreadyAnswered,

		// responders
		"uniq_responder_email_on_system": errs.ErrResponderEmailTaken,
		"chk_invite_not_both_named":      errs.ErrInviteBothNamed,
		"form_invites_token_key":         errs.ErrInviteTokenTaken,
	}
}
