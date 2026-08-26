package jobs

import (
	"context"
	"fmt"
	"lib/crypto"
	"lib/email"
	"lib/telemetry"
	"os"
	"univents/assets"
	"univents/models"
	"univents/ports"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

type GrantCertsForEditionArgs struct {
	EditionID uuid.UUID `json:"edition_id"`
}

func (GrantCertsForEditionArgs) Kind() string { return "cert.grant_edition" }

type GrantCertsForOccurrenceArgs struct {
	EditionID uuid.UUID `json:"edition_id"`
	ProgramID uuid.UUID `json:"program_id"`
}

func (GrantCertsForOccurrenceArgs) Kind() string { return "cert.grant_occurrence" }

type grantCertsDeps struct {
	certs    ports.CertificationRepo
	editions ports.EditionRepo
	events   ports.EventRepo
	email    *email.Client
}

type GrantCertsWorker struct {
	river.WorkerDefaults[GrantCertsForEditionArgs]
	grantCertsDeps
}

func NewGrantCertsWorker(
	certs ports.CertificationRepo,
	editions ports.EditionRepo,
	events ports.EventRepo,
	email *email.Client,
) *GrantCertsWorker {
	return &GrantCertsWorker{grantCertsDeps: grantCertsDeps{
		certs: certs, editions: editions, events: events, email: email,
	}}
}

func (w *GrantCertsWorker) Work(ctx context.Context, job *river.Job[GrantCertsForEditionArgs]) error {
	errs := w.grantEditionAttendanceCerts(ctx, job.Args.EditionID)
	for _, e := range errs {
		telemetry.Log().Warn("cert emission had errors",
			zap.String("user_id", e.UserID.String()),
			zap.String("error", e.Error),
		)
	}
	if len(errs) > 0 {
		return fmt.Errorf("cert emission for edition %s had %d errors", job.Args.EditionID, len(errs))
	}
	return nil
}

type GrantCertsForOccurrenceWorker struct {
	river.WorkerDefaults[GrantCertsForOccurrenceArgs]
	grantCertsDeps
}

func NewGrantCertsForOccurrenceWorker(
	certs ports.CertificationRepo,
	editions ports.EditionRepo,
	events ports.EventRepo,
	email *email.Client,
) *GrantCertsForOccurrenceWorker {
	return &GrantCertsForOccurrenceWorker{grantCertsDeps: grantCertsDeps{
		certs: certs, editions: editions, events: events, email: email,
	}}
}

func (w *GrantCertsForOccurrenceWorker) Work(ctx context.Context, job *river.Job[GrantCertsForOccurrenceArgs]) error {
	errs := w.grantProgramCerts(ctx, job.Args.EditionID, job.Args.ProgramID)
	for _, e := range errs {
		telemetry.Log().Warn("cert emission had errors",
			zap.String("user_id", e.UserID.String()),
			zap.String("error", e.Error),
		)
	}
	if len(errs) > 0 {
		return fmt.Errorf("cert emission for program %s had %d errors", job.Args.ProgramID, len(errs))
	}
	return nil
}

type grantCertsError struct {
	UserID uuid.UUID `json:"user_id"`
	Error  string    `json:"error"`
}

func (w *grantCertsDeps) grantEditionAttendanceCerts(ctx context.Context, editionID uuid.UUID) []grantCertsError {
	template := w.findTemplate(ctx, editionID, nil)

	attendees, err := w.certs.ListDistinctRegistrationsByEdition(ctx, editionID)
	if err != nil {
		return []grantCertsError{{UserID: uuid.Nil, Error: fmt.Sprintf("failed to list attendees: %v", err)}}
	}

	var templateID *uuid.UUID
	if template != nil {
		templateID = &template.ID
	}

	eligible := make([]models.CertEligibleAttendee, 0, len(attendees))
	for _, a := range attendees {
		if a.UserID == uuid.Nil {
			// Accountless gifted ticket (email-only recipient): no profile
			// to own the cert yet — skipped until the claim flow ties an
			// actor id and the grant job re-runs.
			continue
		}
		hasCert, err := w.certs.HasCertForRegistration(ctx, a.RegistrationID, templateID)
		if err != nil {
			w.recordEmissionError(ctx, editionID, a.UserID, templateID, nil, fmt.Sprintf("failed to check existing cert: %v", err))
			continue
		}
		if hasCert {
			continue
		}
		eligible = append(eligible, a)
	}

	// TODO: in the future, require check-in for edition attendance certs (not just confirmed registration)
	return w.emitCerts(ctx, editionID, template, eligible, nil)
}

func (w *grantCertsDeps) grantProgramCerts(ctx context.Context, editionID, programID uuid.UUID) []grantCertsError {
	template := w.findTemplate(ctx, editionID, &programID)

	participants, err := w.certs.ListDistinctParticipantsByProgram(ctx, programID)
	if err != nil {
		return []grantCertsError{{UserID: uuid.Nil, Error: fmt.Sprintf("failed to list participants: %v", err)}}
	}

	eligible := make([]models.CertEligibleAttendee, 0, len(participants))
	for _, p := range participants {
		if p.UserID == uuid.Nil {
			// Accountless gifted ticket (email-only recipient): no profile
			// to own the cert yet — skipped until the claim flow ties an
			// actor id and the grant job re-runs.
			continue
		}
		hasCert, err := w.certs.HasCertForProgram(ctx, p.UserID, programID)
		if err != nil {
			w.recordEmissionError(ctx, editionID, p.UserID, nil, &programID, fmt.Sprintf("failed to check existing cert: %v", err))
			continue
		}
		if hasCert {
			continue
		}
		eligible = append(eligible, p)
	}

	return w.emitCerts(ctx, editionID, template, eligible, &programID)
}

func (w *grantCertsDeps) findTemplate(ctx context.Context, editionID uuid.UUID, programID *uuid.UUID) *models.CertificationTemplate {
	templates, err := w.certs.ListTemplatesForEmission(ctx, editionID)
	if err != nil || len(templates) == 0 {
		return nil
	}
	kind := models.CertificationTemplateKindEditionAttendance
	if programID != nil {
		kind = models.CertificationTemplateKindProgramAttendance
		for _, t := range templates {
			if t.Kind != kind {
				continue
			}
			programs, err := w.certs.ListCertTemplateLinks(ctx, t.ID)
			if err != nil {
				continue
			}
			for _, tp := range programs {
				if tp.ProgramID == *programID {
					return &t
				}
			}
		}
		return nil
	}
	for _, t := range templates {
		if t.Kind == kind {
			return &t
		}
	}
	return nil
}

func (w *grantCertsDeps) emitCerts(ctx context.Context, editionID uuid.UUID, template *models.CertificationTemplate, attendees []models.CertEligibleAttendee, programID *uuid.UUID) []grantCertsError {
	var errors []grantCertsError
	var templateID *uuid.UUID
	templateName := "Certificate"
	if template != nil {
		templateID = &template.ID
		templateName = template.Name
	}

	for _, a := range attendees {
		cert, err := w.certs.Certify(ctx, models.CertifyInput{
			EditionID:        editionID,
			TemplateID:       templateID,
			RegistrationID:   a.RegistrationID,
			UserID:           a.UserID,
			ProgramID:        programID,
			VerificationHash: crypto.HashHMACSHA256(uuid.New().String()),
		})
		if err != nil {
			errors = append(errors, grantCertsError{UserID: a.UserID, Error: fmt.Sprintf("failed to emit cert: %v", err)})
			w.recordEmissionError(ctx, editionID, a.UserID, templateID, programID, fmt.Sprintf("failed to emit cert: %v", err))
			continue
		}

		//nolint:gosec // context.WithoutCancel is safe — it detaches cancellation while preserving values
		go w.sendCertEmail(context.WithoutCancel(ctx), cert, templateName, a.AttendeeEmail, a.AttendeeName)
	}

	return errors
}

func (w *grantCertsDeps) sendCertEmail(ctx context.Context, cert *models.Certification, certName, toEmail, toName string) {
	edition, err := w.editions.GetByID(ctx, cert.EditionID)
	if err != nil {
		telemetry.Log().Error("failed to get edition for cert email", zap.Error(err))
		return
	}

	event, err := w.events.GetByID(ctx, edition.EventID)
	if err != nil {
		telemetry.Log().Error("failed to get event for cert email", zap.Error(err))
		return
	}

	appURL := os.Getenv("APP_URL")
	certLink := fmt.Sprintf("%s/certifications/%s", appURL, cert.ID)
	verifyLink := fmt.Sprintf("%s/verify/%s", appURL, cert.VerificationHash)

	body, err := assets.RenderCertGrantedEmail(assets.CertGrantedEmailData{
		AttendeeName: toName,
		EventName:    event.FullName,
		EditionName:  edition.Name,
		CertName:     certName,
		CertLink:     certLink,
		VerifyLink:   verifyLink,
	})
	if err != nil {
		telemetry.Log().Error("failed to render cert email", zap.Error(err))
		return
	}

	err = w.email.Send(email.Message{
		To:      []string{toEmail},
		Subject: fmt.Sprintf("Your certificate for %s — %s", event.FullName, edition.Name),
		Body:    body,
		HTML:    true,
	})
	if err != nil {
		telemetry.Log().Error("failed to send cert email", zap.Error(err))
		return
	}

	_ = w.certs.MarkEmailSent(context.Background(), cert.ID)
}

func (w *grantCertsDeps) recordEmissionError(ctx context.Context, editionID, userID uuid.UUID, templateID *uuid.UUID, programID *uuid.UUID, errMsg string) {
	_ = w.certs.RecordEmissionError(ctx, &models.CertEmissionError{
		EditionID:    editionID,
		UserID:       userID,
		TemplateID:   templateID,
		ProgramID:    programID,
		ErrorMessage: errMsg,
	})
}
