package app

import (
	"net/http"

	"Informd/errs"
	"lib/httpserver"
)

// Errors is Informd's one table from named errors to HTTP (ADR-0016). The
// most specific error in a chain wins, and an error missing here is a
// generic 500. The shared lib errors (missing rows, constraint classes,
// batches, insufficient role) are mapped by lib/httpserver after these.
func Errors() httpserver.ErrorMap {
	const (
		badRequest = http.StatusBadRequest
		forbidden  = http.StatusForbidden
		notFound   = http.StatusNotFound
		conflict   = http.StatusConflict
	)
	is := httpserver.Is
	invalid := func(target error) httpserver.ErrorMapping { return is(target, badRequest).Validation() }

	return httpserver.ErrorMap{
		// form lifecycle
		is(errs.ErrOpenNotDraft, badRequest),
		is(errs.ErrCloseNotOpen, badRequest),
		is(errs.ErrRedraftNotOpen, badRequest),
		is(errs.ErrRedraftWithResponses, badRequest),
		is(errs.ErrArchiveNotClosed, badRequest),
		is(errs.ErrFormNotAnswerable, forbidden),
		is(errs.ErrFormNotOpenForResponses, badRequest),

		// form contents
		is(errs.ErrStepsAcrossForms, badRequest),
		is(errs.ErrFieldsAcrossSteps, badRequest),
		is(errs.ErrInvalidSelectOptions, badRequest),

		// membership
		is(errs.ErrAddSelfToForm, badRequest),
		is(errs.ErrFormOwnerAlreadyMember, badRequest),
		is(errs.ErrNamespaceOwnerAlreadyInForm, badRequest),
		is(errs.ErrNamespaceMemberAlreadyInForm, badRequest),
		is(errs.ErrAlreadyFormMember, badRequest),
		is(errs.ErrRemoveSelfFromForm, badRequest),
		is(errs.ErrRemoveFormOwner, badRequest),
		is(errs.ErrRemoveNamespaceOwnerFromForm, badRequest),
		is(errs.ErrRemoveNamespaceMemberFromForm, badRequest),
		is(errs.ErrNotFormMember, badRequest),
		is(errs.ErrAddSelfToNamespace, badRequest),
		is(errs.ErrAddOwnerToNamespace, badRequest),
		is(errs.ErrAlreadyNamespaceMember, badRequest),
		is(errs.ErrNotNamespaceMember, badRequest),

		// not found
		is(errs.ErrFormNotFound, notFound),
		is(errs.ErrNamespaceNotFound, notFound),
		is(errs.ErrStepNotFound, notFound),
		is(errs.ErrFieldNotFound, notFound),
		is(errs.ErrAnswerNotFound, notFound),
		is(errs.ErrResponseNotFound, notFound),
		is(errs.ErrResponderNotFound, notFound),

		// constraints: unique → 409, check → 400 validation
		invalid(errs.ErrInvalidFormStatus),
		invalid(errs.ErrInvalidFormStatusState),
		is(errs.ErrFormNameTaken, conflict),
		is(errs.ErrAPIKeyNameTaken, conflict),
		invalid(errs.ErrInvalidFieldType),
		invalid(errs.ErrInvalidFieldKeyFormat),
		is(errs.ErrFieldKeyTaken, conflict),
		invalid(errs.ErrInvalidSelectBehavior),
		invalid(errs.ErrEmptySelectOptions),
		invalid(errs.ErrInvalidSelectType),
		is(errs.ErrNamespaceNameTaken, conflict),
		invalid(errs.ErrInvalidNamespaceMemberRole),
		is(errs.ErrFieldAlreadyAnswered, conflict),
		is(errs.ErrResponderEmailTaken, conflict),
		invalid(errs.ErrInviteBothNamed),
		is(errs.ErrInviteTokenTaken, conflict),
	}
}
