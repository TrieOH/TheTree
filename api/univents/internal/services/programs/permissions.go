package programs

import (
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"

	"univents/models"
)

// checkEligible resolves the caller's confirmed registration in the
// occurrence's edition and enforces the activity's access rules (D1):
//
//   - a confirmed (paid or free) ticket is required — pending/unpaid
//     reservations cannot hold activity spots;
//   - staff_only programs additionally require an event staff role
//     (owner/admin/staff);
//   - min_access_level gates on the ticket's access level (null = open).
//
// Returns the registration and program, or a typed fun error (403).
func (o *Operations) checkEligible(ctx context.Context, occ *models.ProgramOccurrence, userID uuid.UUID) (*models.Registration, *models.Program, error) {
	reg, err := o.registrations.GetActiveByEditionAndAttendee(ctx, occ.EditionID, userID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return nil, nil, fun.ErrForbidden("a confirmed ticket is required to register for activities")
		}
		return nil, nil, err
	}
	if reg.Status != models.RegistrationStatusConfirmed {
		return nil, nil, fun.ErrForbidden("a confirmed ticket is required to register for activities")
	}

	program, err := o.programs.GetByID(ctx, occ.ProgramID)
	if err != nil {
		return nil, nil, err
	}

	err = o.checkProgramAccess(ctx, occ, program, reg, userID)
	if err != nil {
		return nil, nil, err
	}
	return reg, program, nil
}

// checkProgramAccess enforces the program's access rules for an attendee
// (shared by register and checkpoint check-in):
//
//   - staff_only programs additionally require an event staff role
//     (owner/admin/staff);
//   - min_access_level gates on the ticket's access level (null = open).
func (o *Operations) checkProgramAccess(ctx context.Context, occ *models.ProgramOccurrence, program *models.Program, reg *models.Registration, userID uuid.UUID) error {
	if program.StaffOnly {
		edition, err := o.editions.GetByID(ctx, occ.EditionID)
		if err != nil {
			return err
		}
		err = o.authz.CheckEvent(ctx, userID, edition.EventID, models.EventMemberRoleStaff)
		if err != nil {
			return err
		}
	}

	if program.MinAccessLevel != nil {
		ticket, err := o.ticketTypes.GetByID(ctx, reg.TicketTypeID)
		if err != nil {
			return err
		}
		if ticket.AccessLevel < *program.MinAccessLevel {
			return fun.ErrForbidden("your ticket level does not grant access to this " + string(program.Kind))
		}
	}

	return nil
}
