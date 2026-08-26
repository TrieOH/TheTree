package checkouts_test

import (
	"context"
	"testing"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"

	"univents/internal/repos"
	"univents/internal/services/checkouts"
	"univents/models"

	"sdk/payssage"
)

// pixFn is the happy-path intent fake shared by the gift tests.
func pixFn(_ uuid.UUID, _ payssage.CreateIntentRequest) (*payssage.Intent, error) {
	return pixIntent(), nil
}

// giftRegistration loads the registration materialized for a checkout's
// first ticket line.
func giftRegistration(t *testing.T, r *repos.Repos, res *checkouts.CheckoutResult) *models.Registration {
	t.Helper()
	items, err := r.Purchases.ListItemsByPurchase(context.Background(), res.PurchaseID)
	if err != nil {
		t.Fatalf("load items: %v", err)
	}
	if len(items) == 0 || items[0].RegistrationID == nil {
		t.Fatal("checkout did not materialize a registration")
	}
	reg, err := r.Registrations.GetByID(context.Background(), *items[0].RegistrationID)
	if err != nil {
		t.Fatalf("load registration: %v", err)
	}
	return reg
}

// TestCheckout_GiftEmailOnlyUnknownAccount pins the email-only gift: a
// recipient with no IdentityX account is materialized with a nil
// attendee_user_id (the ticket is claimable after they create an account)
// and the email is normalized (trim + lowercase) at checkout.
func TestCheckout_GiftEmailOnlyUnknownAccount(t *testing.T) {
	r, ops, fs, _ := newOps(t, pixFn)
	fx := seedStore(t, r)
	buyerID := uuid.New()
	fs.actors.seed("buyer@example.com", buyerID)

	res, err := ops.Checkout(context.Background(), fx.editionID, buyerID, pixInput(
		ticketLine(fx.ticketID, &checkouts.Attendee{Email: "Friend@Example.com ", Name: "John Doe"}),
	))
	if err != nil {
		t.Fatalf("Checkout: %v", err)
	}

	reg := giftRegistration(t, r, res)
	if reg.AttendeeUserID != nil {
		t.Fatalf("attendee_user_id = %v, want nil (no account yet)", *reg.AttendeeUserID)
	}
	if reg.AttendeeEmail != "friend@example.com" {
		t.Fatalf("attendee_email = %q, want normalized friend@example.com", reg.AttendeeEmail)
	}
}

// TestCheckout_GiftEmailOnlyExistingAccount pins the tie: an email-only
// attendee whose email already has an account is bound to that actor id at
// checkout (the gift is never left unclaimed for an existing user).
func TestCheckout_GiftEmailOnlyExistingAccount(t *testing.T) {
	r, ops, fs, _ := newOps(t, pixFn)
	fx := seedStore(t, r)
	buyerID := uuid.New()
	recipientID := uuid.New()
	fs.actors.seed("buyer@example.com", buyerID)
	fs.actors.seed("friend@example.com", recipientID)

	res, err := ops.Checkout(context.Background(), fx.editionID, buyerID, pixInput(
		ticketLine(fx.ticketID, &checkouts.Attendee{Email: "friend@example.com", Name: "John Doe"}),
	))
	if err != nil {
		t.Fatalf("Checkout: %v", err)
	}

	reg := giftRegistration(t, r, res)
	if reg.AttendeeUserID == nil || *reg.AttendeeUserID != recipientID {
		t.Fatalf("attendee_user_id = %v, want the resolved account %s", reg.AttendeeUserID, recipientID)
	}
}

// TestCheckout_GiftIdEmailMismatch400 pins the id+email consistency check:
// an attendee whose user_id does not belong to the email is rejected 400
// (the pair cannot both be true).
func TestCheckout_GiftIdEmailMismatch400(t *testing.T) {
	r, ops, fs, _ := newOps(t, pixFn)
	fx := seedStore(t, r)
	buyerID := uuid.New()
	accountID := uuid.New()
	wrongID := uuid.New()
	fs.actors.seed("buyer@example.com", buyerID)
	fs.actors.seed("friend@example.com", accountID)

	_, err := ops.Checkout(context.Background(), fx.editionID, buyerID, pixInput(
		ticketLine(fx.ticketID, &checkouts.Attendee{UserID: &wrongID, Email: "friend@example.com", Name: "John Doe"}),
	))
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("err = %v, want BAD_REQUEST (user_id does not match the email)", err)
	}
}

// TestCheckout_GiftIdWithoutAccountEmail400 pins the id+email consistency
// check's other side: a user_id claimed for an email that has no account at
// all is inconsistent — the email cannot belong to the claimed user.
func TestCheckout_GiftIdWithoutAccountEmail400(t *testing.T) {
	r, ops, fs, _ := newOps(t, pixFn)
	fx := seedStore(t, r)
	buyerID := uuid.New()
	fs.actors.seed("buyer@example.com", buyerID)

	_, err := ops.Checkout(context.Background(), fx.editionID, buyerID, pixInput(
		ticketLine(fx.ticketID, &checkouts.Attendee{UserID: new(uuid.New()), Email: "nobody@example.com", Name: "John Doe"}),
	))
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("err = %v, want BAD_REQUEST (email has no account)", err)
	}
}

// TestCheckout_GiftTwoEmailOnlySameAddress400 pins the cart dedup for
// accountless recipients: two lines for the same email are one person.
func TestCheckout_GiftTwoEmailOnlySameAddress400(t *testing.T) {
	r, ops, fs, _ := newOps(t, pixFn)
	fx := seedStore(t, r)
	buyerID := uuid.New()
	fs.actors.seed("buyer@example.com", buyerID)

	_, err := ops.Checkout(context.Background(), fx.editionID, buyerID, pixInput(
		ticketLine(fx.ticketID, &checkouts.Attendee{Email: "friend@example.com", Name: "John Doe"}),
		ticketLine(fx.ticketID, &checkouts.Attendee{Email: "friend@example.com", Name: "John Doe"}),
	))
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("err = %v, want BAD_REQUEST (one ticket per person)", err)
	}
}

// TestCheckout_GiftResolvedAndExplicitSamePerson400 pins the post-resolution
// dedup: an email-only line that resolves to an account collides with an
// explicit id line for the same person.
func TestCheckout_GiftResolvedAndExplicitSamePerson400(t *testing.T) {
	r, ops, fs, _ := newOps(t, pixFn)
	fx := seedStore(t, r)
	buyerID := uuid.New()
	recipientID := uuid.New()
	fs.actors.seed("buyer@example.com", buyerID)
	fs.actors.seed("friend@example.com", recipientID)

	_, err := ops.Checkout(context.Background(), fx.editionID, buyerID, pixInput(
		ticketLine(fx.ticketID, &checkouts.Attendee{UserID: &recipientID, Email: "friend@example.com", Name: "John Doe"}),
		ticketLine(fx.ticketID, &checkouts.Attendee{Email: "friend@example.com", Name: "John Doe"}),
	))
	if err == nil || !fun.Is(err, fun.CodeBadRequest) {
		t.Fatalf("err = %v, want BAD_REQUEST (resolved to the same person)", err)
	}
}

// TestCheckout_GiftEmailOnlyHolderBlocks409 pins the one-ticket-per-person
// rule for accountless recipients: a second email-only gift to an address
// that already holds an active registration is 409, whether the first gift
// is pending or confirmed.
func TestCheckout_GiftEmailOnlyHolderBlocks409(t *testing.T) {
	r, ops, fs, _ := newOps(t, pixFn)
	fx := seedStore(t, r)
	buyerID := uuid.New()
	fs.actors.seed("buyer@example.com", buyerID)

	// First gift: email-only reservation (paid path stays pending).
	res, err := ops.Checkout(context.Background(), fx.editionID, buyerID, pixInput(
		ticketLine(fx.ticketID, &checkouts.Attendee{Email: "friend@example.com", Name: "John Doe"}),
	))
	if err != nil {
		t.Fatalf("first gift: %v", err)
	}
	_ = res

	// A second purchaser gifting to the same address → 409 (the email-only
	// slot is held), not a duplicate registration.
	otherBuyerID := uuid.New()
	fs.actors.seed("buyer2@example.com", otherBuyerID)
	_, err = ops.Checkout(context.Background(), fx.editionID, otherBuyerID, pixInput(
		ticketLine(fx.ticketID, &checkouts.Attendee{Email: "friend@example.com", Name: "John Doe"}),
	))
	if err == nil || !fun.Is(err, fun.CodeConflict) {
		t.Fatalf("second gift: err = %v, want CONFLICT (one ticket per person)", err)
	}
}

// TestCheckout_GiftEmailOnlyFreeOrder pins the free-order path for
// email-only recipients: total 0 confirms immediately, the registration is
// confirmed with a nil attendee_user_id, and no badge is emitted yet (no
// account to attach it to — the claim flow re-emits).
func TestCheckout_GiftEmailOnlyFreeOrder(t *testing.T) {
	r, ops, fs, _ := newOps(t, pixFn)
	fx := seedStore(t, r)
	buyerID := uuid.New()
	fs.actors.seed("buyer@example.com", buyerID)

	freeTicket, err := r.TicketTypes.Create(context.Background(), &models.TicketType{
		EditionID:   fx.editionID,
		Name:        "Free",
		AccessLevel: 0,
		PriceCents:  0,
		MaxQuantity: new(int(1)),
	})
	if err != nil {
		t.Fatalf("seed free ticket: %v", err)
	}

	res, err := ops.Checkout(context.Background(), fx.editionID, buyerID, pixInput(
		ticketLine(freeTicket.ID, &checkouts.Attendee{Email: "friend@example.com", Name: "John Doe"}),
	))
	if err != nil {
		t.Fatalf("Checkout: %v", err)
	}
	if res.Status != models.PurchaseStatusApproved {
		t.Fatalf("status = %s, want approved (free order)", res.Status)
	}

	reg := giftRegistration(t, r, res)
	if reg.Status != models.RegistrationStatusConfirmed || reg.AttendeeUserID != nil {
		t.Fatalf("registration = %+v, want confirmed with no account", reg)
	}

	// The gifted-ticket email is enqueued atomically with the confirmation
	// (the recipient must be told to create an account and claim).
	giftJobs := fs.river.giftEmailArgs()
	if len(giftJobs) != 1 || giftJobs[0].RegistrationID != reg.ID {
		t.Fatalf("gifts.send_email jobs = %+v, want [{registration: %s}]", giftJobs, reg.ID)
	}
	// Badge emission for accountless recipients is deferred (no profile to
	// attach it to) — covered by the badges-package unit test.
}
