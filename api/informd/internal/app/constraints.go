package app

import (
	"lib/database"
)

func SetupConstraintMessages() {
	database.SetConstraintErrorRegistry(database.ConstraintRegistry{
		// forms
		"chk_forms_valid_status":       "status must be one of: draft, open, closed, archived",
		"chk_forms_valid_status_state": "opened_at, closed_at or archived_at must be set when status is open, closed or archived",
		"uniq_form_name_per_namespace": "a form with this name already exists in this namespace",
		"uniq_name_per_user":           "an API key with this name already exists",

		// fields
		"chk_fields_type":       "field type must be one of: string, email, int, float, bool, date, time, datetime, select, file, phone, url",
		"chk_fields_key_format": "field key must start with a letter or underscore and contain only lowercase letters, digits and underscores",
		"uniq_key_per_step":     "a field with this key already exists in this step",

		// select
		"chk_select_behaviour":  "select fields must have a behaviour of checkbox or radio",
		"chk_select_options":    "select fields must have a non-empty options array",
		"chk_select_value_type": "select_type must be one of: email, int, float, date, time, datetime, phone, url",

		// namespaces
		"uniq_namespace_name_per_user":    "a namespace with this name already exists",
		"chk_valid_namespace_member_role": "invalid namespace member role",

		// response
		"uniq_answer_per_field_per_response": "this field already has an answer in this response",

		// responders
		"uniq_responder_email_on_system": "a responder with this email already exists on this system",
		"chk_invite_not_both_named":      "an invite cannot have both responder_id and email set",
		"form_invites_token_key":         "an invite with this token already exists",
	})
}
