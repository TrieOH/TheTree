package app

import (
	"errors"

	"Informd/errs"
	"Informd/internal/problems"

	fp "github.com/MintzyG/problematic"
)

// Problem answers Informd's named errors with their declared problems
// (ADR-0019), and nil for an error it doesn't know; problems.Mapper then
// tries the shared cases. Named errors don't wrap named errors, so case
// order never breaks a tie. The detail is the named error's message, never
// the chain's.
func Problem(err error) *fp.Problem {
	switch {
	// form lifecycle
	case errors.Is(err, errs.ErrOpenNotDraft):
		return problems.OpenNotDraft(errs.ErrOpenNotDraft.Error())
	case errors.Is(err, errs.ErrCloseNotOpen):
		return problems.CloseNotOpen(errs.ErrCloseNotOpen.Error())
	case errors.Is(err, errs.ErrRedraftNotOpen):
		return problems.RedraftNotOpen(errs.ErrRedraftNotOpen.Error())
	case errors.Is(err, errs.ErrRedraftWithResponses):
		return problems.RedraftWithResponses(errs.ErrRedraftWithResponses.Error())
	case errors.Is(err, errs.ErrArchiveNotClosed):
		return problems.ArchiveNotClosed(errs.ErrArchiveNotClosed.Error())
	case errors.Is(err, errs.ErrFormNotAnswerable):
		return problems.FormNotAnswerable(errs.ErrFormNotAnswerable.Error())
	case errors.Is(err, errs.ErrFormNotOpenForResponses):
		return problems.FormNotOpenForResponses(errs.ErrFormNotOpenForResponses.Error())

	// form contents
	case errors.Is(err, errs.ErrStepsAcrossForms):
		return problems.StepsAcrossForms(errs.ErrStepsAcrossForms.Error())
	case errors.Is(err, errs.ErrFieldsAcrossSteps):
		return problems.FieldsAcrossSteps(errs.ErrFieldsAcrossSteps.Error())
	case errors.Is(err, errs.ErrInvalidSelectOptions):
		return problems.InvalidSelectOptions(errs.ErrInvalidSelectOptions.Error())

	// membership
	case errors.Is(err, errs.ErrAddSelfToForm):
		return problems.AddSelfToForm(errs.ErrAddSelfToForm.Error())
	case errors.Is(err, errs.ErrFormOwnerAlreadyMember):
		return problems.FormOwnerAlreadyMember(errs.ErrFormOwnerAlreadyMember.Error())
	case errors.Is(err, errs.ErrNamespaceOwnerAlreadyInForm):
		return problems.NamespaceOwnerAlreadyInForm(errs.ErrNamespaceOwnerAlreadyInForm.Error())
	case errors.Is(err, errs.ErrNamespaceMemberAlreadyInForm):
		return problems.NamespaceMemberAlreadyInForm(errs.ErrNamespaceMemberAlreadyInForm.Error())
	case errors.Is(err, errs.ErrAlreadyFormMember):
		return problems.AlreadyFormMember(errs.ErrAlreadyFormMember.Error())
	case errors.Is(err, errs.ErrRemoveSelfFromForm):
		return problems.RemoveSelfFromForm(errs.ErrRemoveSelfFromForm.Error())
	case errors.Is(err, errs.ErrRemoveFormOwner):
		return problems.RemoveFormOwner(errs.ErrRemoveFormOwner.Error())
	case errors.Is(err, errs.ErrRemoveNamespaceOwnerFromForm):
		return problems.RemoveNamespaceOwnerFromForm(errs.ErrRemoveNamespaceOwnerFromForm.Error())
	case errors.Is(err, errs.ErrRemoveNamespaceMemberFromForm):
		return problems.RemoveNamespaceMemberFromForm(errs.ErrRemoveNamespaceMemberFromForm.Error())
	case errors.Is(err, errs.ErrNotFormMember):
		return problems.NotFormMember(errs.ErrNotFormMember.Error())
	case errors.Is(err, errs.ErrAddSelfToNamespace):
		return problems.AddSelfToNamespace(errs.ErrAddSelfToNamespace.Error())
	case errors.Is(err, errs.ErrAddOwnerToNamespace):
		return problems.AddOwnerToNamespace(errs.ErrAddOwnerToNamespace.Error())
	case errors.Is(err, errs.ErrAlreadyNamespaceMember):
		return problems.AlreadyNamespaceMember(errs.ErrAlreadyNamespaceMember.Error())
	case errors.Is(err, errs.ErrNotNamespaceMember):
		return problems.NotNamespaceMember(errs.ErrNotNamespaceMember.Error())

	// constraints
	case errors.Is(err, errs.ErrFormNameTaken):
		return problems.FormNameTaken(errs.ErrFormNameTaken.Error())
	case errors.Is(err, errs.ErrAPIKeyNameTaken):
		return problems.APIKeyNameTaken(errs.ErrAPIKeyNameTaken.Error())
	case errors.Is(err, errs.ErrFieldKeyTaken):
		return problems.FieldKeyTaken(errs.ErrFieldKeyTaken.Error())
	case errors.Is(err, errs.ErrNamespaceNameTaken):
		return problems.NamespaceNameTaken(errs.ErrNamespaceNameTaken.Error())
	case errors.Is(err, errs.ErrFieldAlreadyAnswered):
		return problems.FieldAlreadyAnswered(errs.ErrFieldAlreadyAnswered.Error())
	case errors.Is(err, errs.ErrResponderEmailTaken):
		return problems.ResponderEmailTaken(errs.ErrResponderEmailTaken.Error())
	case errors.Is(err, errs.ErrInviteTokenTaken):
		return problems.InviteTokenTaken(errs.ErrInviteTokenTaken.Error())
	}
	for _, target := range notFound {
		if errors.Is(err, target) {
			return problems.NotFound(target.Error())
		}
	}
	for _, target := range invalid {
		if errors.Is(err, target) {
			return problems.Validation(target.Error())
		}
	}
	return nil
}

// notFound are the named missing resources, answered as the shared
// not-found problem.
var notFound = []error{
	errs.ErrFormNotFound,
	errs.ErrNamespaceNotFound,
	errs.ErrStepNotFound,
	errs.ErrFieldNotFound,
	errs.ErrAnswerNotFound,
	errs.ErrResponseNotFound,
	errs.ErrResponderNotFound,
}

// invalid are the check constraints, answered as the shared validation
// problem.
var invalid = []error{
	errs.ErrInvalidFormStatus,
	errs.ErrInvalidFormStatusState,
	errs.ErrInvalidFieldType,
	errs.ErrInvalidFieldKeyFormat,
	errs.ErrInvalidSelectBehavior,
	errs.ErrEmptySelectOptions,
	errs.ErrInvalidSelectType,
	errs.ErrInvalidNamespaceMemberRole,
	errs.ErrInviteBothNamed,
}
