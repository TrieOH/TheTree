// Package errmap is Univents' problem switch (ADR-0019). It sits apart
// from internal/app, which wires the features, so service tests can check
// what status an error answers with without importing the features they
// belong to.
package errmap

import (
	"errors"

	"univents/errs"
	"univents/internal/problems"

	fp "github.com/MintzyG/problematic"
)

// Problem answers what Univents' problems can't name for themselves
// (ADR-0019): the typed errors, which carry the problem's members, the
// named errors answered with a shared problem, the not-found errors and the
// check constraints. A problem that names its sentinel (x-problem's error)
// is answered by problems.Mapper itself, which calls Problem first and
// tries the shared cases last; Problem is nil for an error it doesn't
// know. Named errors don't wrap named errors, so case order never breaks a
// tie. The detail is the named error's message, never the chain's.
func Problem(err error) *fp.Problem {
	var (
		rejected       *errs.PaymentRejectedError
		purchaseStatus *errs.PurchaseStatusError
		accessLevel    *errs.AccessLevelError
	)
	switch {
	// checkout
	case errors.As(err, &rejected):
		return problems.PaymentRejected(rejected.Error(), rejected.Message)

	// refunds
	case errors.As(err, &purchaseStatus):
		return problems.PurchaseStatus(purchaseStatus.Error(), purchaseStatus.Status)

	// programs
	case errors.As(err, &accessLevel):
		return problems.AccessLevel(accessLevel.Error(), accessLevel.Kind)

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
	errs.ErrEditionNotFound,
	errs.ErrBadgeTemplateNotFound,
	errs.ErrBadgeEmissionNotFound,
	errs.ErrProgramNotFound,
	errs.ErrPurchaseNotFound,
	errs.ErrTicketTypeNotFound,
	errs.ErrRegistrationNotFound,
	errs.ErrCertificationNotFound,
	errs.ErrEventNotFound,
	errs.ErrProductNotFound,
	errs.ErrSignatureNotFound,
	errs.ErrWSTokenNotFound,
}

// invalid are the check constraints, answered as the shared validation
// problem.
var invalid = []error{
	errs.ErrInvalidEventStatus,
	errs.ErrPublicKeyNeedsSeller,
	errs.ErrInvalidEventMemberRole,
	errs.ErrEditionDates,
	errs.ErrRegistrationAfterStart,
	errs.ErrInvalidRegistrationStatus,
	errs.ErrInvalidProductPurchaseStatus,
	errs.ErrInvalidProgramKind,
	errs.ErrOccurrenceDates,
	errs.ErrInvalidParticipationStatus,
	errs.ErrInvalidSignatureRequestStatus,
	errs.ErrInvalidSignatureStatus,
	errs.ErrReadySignatureNeedsImage,
	errs.ErrInvalidCertificationKind,
	errs.ErrInvalidBadgeEmissionOrigin,
	errs.ErrInvalidBadgeEmissionStatus,
	errs.ErrInvalidBadgeTemplateScope,
	errs.ErrInvalidPurchaseStatus,
	errs.ErrInvalidPurchaseItemType,
	errs.ErrPurchaseItemQuantity,
}
