// Package jobs holds the one declaration of Payssage's background work:
// which job kinds exist and how often the periodic ones run. The handlers
// live with their feature (services/webhooks/jobs); the carrier is chosen
// by the app wiring, and nothing here knows it.
package jobs

import (
	"time"

	libjobs "lib/jobs"
	webhooks "payssage/internal/services/webhooks/jobs"
	"payssage/ports"

	"resty.dev/v3"
)

// Deps are what the handlers need.
type Deps struct {
	Deliveries ports.WebhookDeliveryRepo
	Events     ports.WebhookEventRepo
	Endpoints  ports.WebhookEndpointRepo
	HTTPClient *resty.Client
	// Enqueuer is what the redrive job enqueues through.
	Enqueuer libjobs.Enqueuer
}

// Registry registers every Payssage job kind.
func Registry(d Deps) *libjobs.Registry {
	reg := libjobs.NewRegistry()
	libjobs.Register(reg, webhooks.NewDeliverWebhookWorker(d.Deliveries, d.Events, d.Endpoints, d.HTTPClient).Work)
	libjobs.Register(reg, webhooks.NewFailWebhookDeliveryWorker(d.Deliveries).Work)
	libjobs.Register(reg, webhooks.NewRedriveWebhookDeliveriesWorker(d.Deliveries, d.Enqueuer).Work)
	return reg
}

// RedriveInterval is how often unsent deliveries are re-enqueued. The sweep
// only catches enqueues lost between commit and enqueue, a rare failure, so
// it runs often enough to bound that delay without keeping the Cloudflare
// container awake.
const RedriveInterval = 15 * time.Minute

// Periodic declares the recurring jobs. The in-process runner runs them on
// tickers; on Cloudflare the Worker's cron triggers fire the same kinds, so
// keep its crons in step with these cadences.
func Periodic() []libjobs.Periodic {
	return []libjobs.Periodic{
		{Job: webhooks.RedriveWebhookDeliveriesArgs{}, Every: RedriveInterval, RunOnStart: true},
	}
}
