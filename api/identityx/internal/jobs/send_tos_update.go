package jobs

import (
	"context"

	"IdentityX/internal/emails"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/email"
	"lib/telemetry"

	"github.com/riverqueue/river"
	"go.uber.org/zap"
)

// SendTosUpdateWorker fans a ToS-change notification out to every human
// actor of the project. The job is one per update; per-recipient failures
// are logged and skipped so one bad address cannot block the others, and
// the job completes best-effort (notifications are not transactional
// promises — a full outage is visible in the logs and the update can be
// re-announced by bumping a correction version).
type SendTosUpdateWorker struct {
	river.WorkerDefaults[emails.SendTosUpdateArgs]

	emailClient *email.Client
	actors      ports.ActorRepo
	templates   ports.EmailTemplateRepo
}

func NewSendTosUpdateWorker(emailClient *email.Client, actors ports.ActorRepo, templates ports.EmailTemplateRepo) *SendTosUpdateWorker {
	return &SendTosUpdateWorker{
		emailClient: emailClient,
		actors:      actors,
		templates:   templates,
	}
}

func (w *SendTosUpdateWorker) Work(ctx context.Context, job *river.Job[emails.SendTosUpdateArgs]) error {
	args := job.Args

	actors, err := w.actors.List(ctx, args.ProjectID)
	if err != nil {
		return err
	}

	tpl, err := emails.ResolveTemplate(ctx, w.templates, &args.ProjectID, models.TosEmailTemplateKind)
	if err != nil {
		return err
	}

	for _, actor := range actors {
		// Service and machine actors have no inbox and no consent to
		// notify; only humans with an address on file get the notice.
		if actor.Type != models.HumanActorType || actor.Email == nil {
			continue
		}

		subject, body, err := emails.Render(tpl, emails.Data{
			ProjectName:   args.ProjectName,
			ProjectDomain: emails.DomainHost(args.BaseDomain),
			Email:         *actor.Email,
			TosVersion:    args.TosVersion,
			EffectiveDate: args.EffectiveDate,
			TosContent:    args.TosContent,
		})
		if err != nil {
			telemetry.Log().Error("failed to render tos update email",
				zap.String("actor_id", actor.ID.String()),
				zap.Error(err),
			)
			continue
		}

		err = w.emailClient.Send(email.Message{
			To:      []string{*actor.Email},
			Subject: subject,
			Body:    body,
			HTML:    true,
		})
		if err != nil {
			telemetry.Log().Error("failed to send tos update email",
				zap.String("actor_id", actor.ID.String()),
				zap.String("to", *actor.Email),
				zap.Error(err),
			)
		}
	}
	return nil
}

var _ river.Worker[emails.SendTosUpdateArgs] = (*SendTosUpdateWorker)(nil)
