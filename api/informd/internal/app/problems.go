package app

import (
	"errors"

	"Informd/errs"
	"Informd/internal/problems"

	fp "github.com/MintzyG/problematic"
)

// Problem answers what Informd's problems can't name for themselves
// (ADR-0019): the typed errors, which carry the problem's members, the
// named errors answered with a shared problem, the not-found errors and the
// check constraints. A problem that names its sentinel (x-problem's error)
// is answered by problems.Mapper itself, which calls Problem first and
// tries the shared cases last; Problem is nil for an error it doesn't
// know. Named errors don't wrap named errors, so case order never breaks a
// tie. The detail is the named error's message, never the chain's.
func Problem(err error) *fp.Problem {
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
