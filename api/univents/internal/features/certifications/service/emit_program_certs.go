package service

import (
	"context"
	"lib/telemetry"
	"univents/internal/features/certifications/jobs"
	"univents/models"

	"github.com/google/uuid"
)

// EmitCertsForProgram enqueues the cert grant job for a program's
// participants. The actor must be an owner or admin of the event. The job
// (cert.grant_occurrence) is idempotent — it only emits certificates for
// participants without one for the program — so re-enqueuing is safe and
// acts as a per-program retry for previously failed emissions.
func (o *Operations) EmitCertsForProgram(ctx context.Context, programID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "CertificationsService.EmitCertsForProgram")
	defer span.End()

	program, err := o.Programs.GetByID(ctx, programID)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, program.EditionID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Jobs.Enqueue(ctx, jobs.GrantCertsForOccurrenceArgs{
		EditionID: program.EditionID,
		ProgramID: program.ID,
	})
}
