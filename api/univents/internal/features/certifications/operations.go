// Package certifications is the certifications feature: its repo, handlers
// and operations. Build wires it.
package certifications

import (
	"context"
	"lib/email"
	libjobs "lib/jobs"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/internal/features/certifications/jobs"
	"univents/models"
	"univents/ports"
	"uuid"
)

// Jobs is the cert-grant scheduling surface: a plain enqueue outside any tx
// (the emit endpoint is a fire-and-forget enqueue of the grant job; the job
// itself re-checks existing certs, so re-emission is safe). Satisfied by
// the lib/jobs carriers.
type Jobs interface {
	Enqueue(ctx context.Context, job libjobs.Job) error
}

type Deps struct {
	Events   ports.EventRepo
	Editions ports.EditionRepo
	Certs    ports.CertificationRepo
	Programs ports.ProgramRepo
	Email    *email.Client
	Jobs     Jobs
	Authz    *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) CreateTemplate(ctx context.Context, input models.CreateCertificationTemplateInput) (*models.CertificationTemplate, error) {
	edition, err := o.Editions.GetByID(ctx, input.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Certs.CreateTemplate(ctx, input)
}

func (o *Operations) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	template, err := o.Certs.GetTemplateByID(ctx, id)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Certs.DeleteTemplate(ctx, id)
}

// EmitCertsForProgram enqueues the cert grant job for a program's
// participants. The actor must be an owner or admin of the event. The job
// (cert.grant_occurrence) is idempotent — it only emits certificates for
// participants without one for the program — so re-enqueuing is safe and
// acts as a per-program retry for previously failed emissions.
func (o *Operations) EmitCertsForProgram(ctx context.Context, programID uuid.UUID) error {
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

func (o *Operations) GetCertByHash(ctx context.Context, hash string) (*models.Certification, error) {
	return o.Certs.GetByHash(ctx, hash)
}

func (o *Operations) GetCertByID(ctx context.Context, id uuid.UUID) (*models.Certification, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	cert, err := o.Certs.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if cert.UserID == ident.Sub.ID {
		return cert, nil
	}

	edition, err := o.Editions.GetByID(ctx, cert.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleStaff)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func (o *Operations) GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.CertificationTemplate, error) {
	return o.Certs.GetTemplateByID(ctx, id)
}

func (o *Operations) InvalidateCert(ctx context.Context, id uuid.UUID, reason *string) error {
	cert, err := o.Certs.GetByID(ctx, id)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, cert.EditionID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Certs.Invalidate(ctx, id, reason)
}

func (o *Operations) LinkCertTemplate(ctx context.Context, templateID, programID uuid.UUID) error {
	template, err := o.Certs.GetTemplateByID(ctx, templateID)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Certs.LinkCertTemplate(ctx, templateID, programID)
}

func (o *Operations) ListCertTemplateLinks(ctx context.Context, templateID uuid.UUID) ([]models.CertificationTemplateProgram, error) {
	return o.Certs.ListCertTemplateLinks(ctx, templateID)
}

func (o *Operations) ListCertsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Certification, error) {
	return o.Certs.ListByEdition(ctx, editionID)
}

func (o *Operations) ListCertsByUser(ctx context.Context, userID uuid.UUID) ([]models.Certification, error) {
	return o.Certs.ListByUser(ctx, userID)
}

func (o *Operations) ListEmissionErrors(ctx context.Context, editionID uuid.UUID) ([]models.CertEmissionError, error) {
	return o.Certs.ListEmissionErrorsByEdition(ctx, editionID)
}

func (o *Operations) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.CertificationTemplate, error) {
	return o.Certs.ListTemplates(ctx, editionID)
}

func (o *Operations) UnlinkCertTemplate(ctx context.Context, templateID, programID uuid.UUID) error {
	template, err := o.Certs.GetTemplateByID(ctx, templateID)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Certs.UnlinkCertTemplate(ctx, templateID, programID)
}

func (o *Operations) UpdateTemplate(ctx context.Context, input models.UpdateCertificationTemplateInput) (*models.CertificationTemplate, error) {
	template, err := o.Certs.GetTemplateByID(ctx, input.TemplateID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Certs.UpdateTemplate(ctx, input)
}
