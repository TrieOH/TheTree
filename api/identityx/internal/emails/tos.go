package emails

import (
	"context"
	"fmt"
	"strings"
	"time"

	"IdentityX/models"

	"github.com/MintzyG/fun"
)

// TosNotifier enqueues the ToS-change notification job: one River job per
// update, fanned out to the project's users by the worker. It is the ToS
// module's only contact with the queue; the token-minting Sender stays
// untouched because ToS notifications carry no action token.
type TosNotifier struct {
	enqueuer Enqueuer
}

func NewTosNotifier(enqueuer Enqueuer) *TosNotifier {
	return &TosNotifier{enqueuer: enqueuer}
}

// NotifyUpdate enqueues the notification for a created or updated ToS.
// project must carry a domain — a tenant without one is misconfigured and
// the link-less email would still mislead, so fail loudly instead.
func (n *TosNotifier) NotifyUpdate(ctx context.Context, project *models.Project, tos *models.TermsOfService) error {
	if project.Domain == nil {
		return fun.ErrInternal(fmt.Sprintf(
			"cannot enqueue ToS notification for project %s: no domain configured",
			project.ID,
		))
	}
	domain := *project.Domain
	if !strings.Contains(domain, "://") {
		domain = "https://" + domain
	}

	_, err := n.enqueuer.Insert(ctx, SendTosUpdateArgs{
		ProjectID:     project.ID,
		ProjectName:   project.Name,
		BaseDomain:    domain,
		TosVersion:    tos.Version,
		TosContent:    tos.Content,
		EffectiveDate: tos.EffectiveAt.Format(time.DateOnly),
	}, nil)
	return err
}
