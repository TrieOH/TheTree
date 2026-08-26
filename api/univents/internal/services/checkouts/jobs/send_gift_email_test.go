package jobs_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"lib/database"
	"lib/email"
	"lib/testdb"

	"univents/internal/repos"
	"univents/internal/services/checkouts/jobs"
	"univents/internal/sqlc"
	"univents/models"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// recordingSender fakes the SMTP seam and captures every send.
type recordingSender struct {
	mu     sync.Mutex
	sent   []email.Message
	failAt int // fail the Nth send (1-based); 0 = never fail
}

func (s *recordingSender) Send(msg email.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.failAt > 0 && len(s.sent)+1 == s.failAt {
		return errors.New("smtp down")
	}
	s.sent = append(s.sent, msg)
	return nil
}

func (s *recordingSender) messages() []email.Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]email.Message{}, s.sent...)
}

// seedGiftedReg seeds an event → edition → ticket type → confirmed
// registration for the given attendee (nil AttendeeUserID = accountless
// gift). Returns the registration and its event slug.
func seedGiftedReg(t *testing.T, r *repos.Repos, attendeeUserID *uuid.UUID) (models.Registration, string) {
	t.Helper()
	ctx := context.Background()

	event, err := r.Events.Create(ctx, &models.Event{
		OwnerID:  uuid.New(),
		FullName: "Gift Email Test Event",
		Slug:     "gift-email-test-" + uuid.NewString()[:8],
		Status:   models.EventStatusActive,
	})
	if err != nil {
		t.Fatalf("seed event: %v", err)
	}
	edition, err := r.Editions.Create(ctx, &models.Edition{
		EventID:   event.ID,
		Name:      "Gift Email Test Edition",
		Slug:      "gift-email-test-ed-" + uuid.NewString()[:8],
		StartsAt:  time.Now().Add(-time.Hour),
		EndsAt:    time.Now().Add(24 * time.Hour),
		CreatedBy: event.OwnerID,
	})
	if err != nil {
		t.Fatalf("seed edition: %v", err)
	}
	ticket, err := r.TicketTypes.Create(ctx, &models.TicketType{
		EditionID:   edition.ID,
		Name:        "Standard",
		AccessLevel: 0,
		PriceCents:  1000,
		MaxQuantity: new(int(10)),
	})
	if err != nil {
		t.Fatalf("seed ticket: %v", err)
	}
	reg, err := r.Registrations.Create(ctx, &models.Registration{
		EditionID:      edition.ID,
		TicketTypeID:   ticket.ID,
		PurchaserID:    uuid.New(),
		AttendeeUserID: attendeeUserID,
		AttendeeEmail:  "friend@example.com",
		AttendeeName:   "John Doe",
		Status:         models.RegistrationStatusConfirmed,
	})
	if err != nil {
		t.Fatalf("seed registration: %v", err)
	}
	return *reg, event.Slug
}

func newGiftWorker(t *testing.T, r *repos.Repos, sender *recordingSender) *jobs.SendGiftEmailWorker {
	t.Helper()
	return jobs.NewSendGiftEmailWorker(r.Registrations, r.Editions, r.Events, r.TicketTypes, sender)
}

// TestSendGiftEmail_SendsToAccountlessRecipient pins the gifted-ticket
// email: a confirmed email-only registration gets the claim instructions at
// the recipient's address, naming the event/edition and linking the event
// page (whose my-ticket read performs the lazy claim).
func TestSendGiftEmail_SendsToAccountlessRecipient(t *testing.T) {
	pool := testdb.Postgres(t, "../../../../db/migrations")
	q := sqlc.New(pool)
	database.SetDefaultRunner(database.NewPGXTxRunner(pool))
	r := repos.New(q)

	reg, slug := seedGiftedReg(t, r, nil)
	sender := &recordingSender{}
	worker := newGiftWorker(t, r, sender)

	err := worker.Work(context.Background(), &river.Job[jobs.SendGiftEmailArgs]{
		Args: jobs.SendGiftEmailArgs{RegistrationID: reg.ID},
	})
	if err != nil {
		t.Fatalf("Work: %v", err)
	}

	msgs := sender.messages()
	if len(msgs) != 1 {
		t.Fatalf("sends = %d, want 1", len(msgs))
	}
	msg := msgs[0]
	if len(msg.To) != 1 || msg.To[0] != "friend@example.com" {
		t.Fatalf("to = %v, want [friend@example.com]", msg.To)
	}
	if !strings.Contains(msg.Subject, "Gift Email Test Event") {
		t.Fatalf("subject = %q, want the event name", msg.Subject)
	}
	for _, want := range []string{"John Doe", "Standard", "Gift Email Test Edition", "/events/" + slug} {
		if !strings.Contains(msg.Body, want) {
			t.Fatalf("body missing %q", want)
		}
	}
	if !msg.HTML {
		t.Fatal("message must be HTML")
	}
}

// TestSendGiftEmail_SkipsClaimedRegistration pins the idempotent skip: once
// the recipient claims (attendee_user_id set), the job is a no-op — the
// badge email covers the confirmation.
func TestSendGiftEmail_SkipsClaimedRegistration(t *testing.T) {
	pool := testdb.Postgres(t, "../../../../db/migrations")
	q := sqlc.New(pool)
	database.SetDefaultRunner(database.NewPGXTxRunner(pool))
	r := repos.New(q)

	claimedID := uuid.New()
	reg, _ := seedGiftedReg(t, r, &claimedID)
	sender := &recordingSender{}
	worker := newGiftWorker(t, r, sender)

	err := worker.Work(context.Background(), &river.Job[jobs.SendGiftEmailArgs]{
		Args: jobs.SendGiftEmailArgs{RegistrationID: reg.ID},
	})
	if err != nil {
		t.Fatalf("Work: %v", err)
	}
	if n := len(sender.messages()); n != 0 {
		t.Fatalf("sends = %d, want 0 (recipient already has an account)", n)
	}
}

// TestSendGiftEmail_SkipsUnconfirmedRegistration pins the guard: only
// confirmed registrations announce the gift (an abandoned reservation must
// not email the recipient).
func TestSendGiftEmail_SkipsUnconfirmedRegistration(t *testing.T) {
	pool := testdb.Postgres(t, "../../../../db/migrations")
	q := sqlc.New(pool)
	database.SetDefaultRunner(database.NewPGXTxRunner(pool))
	r := repos.New(q)

	reg, _ := seedGiftedReg(t, r, nil)
	_, err := r.Registrations.UpdateStatus(context.Background(), reg.ID, models.RegistrationStatusPending, nil)
	if err != nil {
		t.Fatalf("flip to pending: %v", err)
	}
	sender := &recordingSender{}
	worker := newGiftWorker(t, r, sender)

	err = worker.Work(context.Background(), &river.Job[jobs.SendGiftEmailArgs]{
		Args: jobs.SendGiftEmailArgs{RegistrationID: reg.ID},
	})
	if err != nil {
		t.Fatalf("Work: %v", err)
	}
	if n := len(sender.messages()); n != 0 {
		t.Fatalf("sends = %d, want 0 (not confirmed)", n)
	}
}

// TestSendGiftEmail_RetriesOnSendFailure pins the retry contract: a send
// failure is returned so River retries with backoff.
func TestSendGiftEmail_RetriesOnSendFailure(t *testing.T) {
	pool := testdb.Postgres(t, "../../../../db/migrations")
	q := sqlc.New(pool)
	database.SetDefaultRunner(database.NewPGXTxRunner(pool))
	r := repos.New(q)

	reg, _ := seedGiftedReg(t, r, nil)
	sender := &recordingSender{failAt: 1}
	worker := newGiftWorker(t, r, sender)

	err := worker.Work(context.Background(), &river.Job[jobs.SendGiftEmailArgs]{
		Args: jobs.SendGiftEmailArgs{RegistrationID: reg.ID},
	})
	if err == nil {
		t.Fatal("Work: err = nil, want the send error (River retries)")
	}
}

// TestSendGiftEmail_DeliversViaSMTP proves the email actually arrives: the
// worker runs against a real SMTP server (the same mailpit image the dev
// compose stack uses, in testcontainers) and the message is verified
// through mailpit's HTTP API. Skipped when no Docker daemon is reachable,
// like the rest of the integration suite.
func TestSendGiftEmail_DeliversViaSMTP(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)

	ctx := context.Background()
	mailpit, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "axllent/mailpit:v1.30.7",
			ExposedPorts: []string{"1025/tcp", "8025/tcp"},
			WaitingFor:   wait.ForHTTP("/api/v1/messages").WithPort("8025/tcp").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("start mailpit: %v", err)
	}
	t.Cleanup(func() { _ = mailpit.Terminate(ctx) })

	host, err := mailpit.Host(ctx)
	if err != nil {
		t.Fatalf("mailpit host: %v", err)
	}
	smtpPort, err := mailpit.MappedPort(ctx, "1025/tcp")
	if err != nil {
		t.Fatalf("mailpit smtp port: %v", err)
	}
	apiPort, err := mailpit.MappedPort(ctx, "8025/tcp")
	if err != nil {
		t.Fatalf("mailpit api port: %v", err)
	}

	pool := testdb.Postgres(t, "../../../../db/migrations")
	q := sqlc.New(pool)
	database.SetDefaultRunner(database.NewPGXTxRunner(pool))
	r := repos.New(q)

	reg, _ := seedGiftedReg(t, r, nil)

	// The real SMTP client, pointed at the mailpit container (same code
	// path the deployed worker uses: lib/email → SMTP).
	smtpPortNum, err := strconv.Atoi(smtpPort.Port())
	if err != nil {
		t.Fatalf("smtp port %q: %v", smtpPort.Port(), err)
	}
	sender := email.NewClient(email.Config{Host: host, Port: smtpPortNum, From: "univents@test", TLS: false})
	worker := jobs.NewSendGiftEmailWorker(r.Registrations, r.Editions, r.Events, r.TicketTypes, sender)

	err = worker.Work(ctx, &river.Job[jobs.SendGiftEmailArgs]{
		Args: jobs.SendGiftEmailArgs{RegistrationID: reg.ID},
	})
	if err != nil {
		t.Fatalf("Work: %v", err)
	}

	// Ask mailpit what it received — this is the delivery proof.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"http://"+net.JoinHostPort(host, apiPort.Port())+"/api/v1/messages", nil)
	if err != nil {
		t.Fatalf("mailpit request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("mailpit api: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read mailpit response: %v", err)
	}
	var inbox struct {
		Total    int `json:"total"`
		Messages []struct {
			Subject string `json:"Subject"`
			To      []struct {
				Address string `json:"Address"`
			} `json:"To"`
		} `json:"messages"`
	}
	err = json.Unmarshal(body, &inbox)
	if err != nil {
		t.Fatalf("decode mailpit response %q: %v", body, err)
	}

	if inbox.Total != 1 || len(inbox.Messages) != 1 {
		t.Fatalf("mailpit messages = %d (%d listed), want exactly 1 delivered", inbox.Total, len(inbox.Messages))
	}
	msg := inbox.Messages[0]
	if len(msg.To) != 1 || msg.To[0].Address != "friend@example.com" {
		t.Fatalf("delivered to = %+v, want [friend@example.com]", msg.To)
	}
	if !strings.Contains(msg.Subject, "Gift Email Test Event") {
		t.Fatalf("delivered subject = %q, want the event name", msg.Subject)
	}
}
