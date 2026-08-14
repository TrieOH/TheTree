package mercado_pago_provider

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"unicode/utf8"

	"payssage/internal/config"
	"payssage/internal/providers/mercado_pago"
	"payssage/models"
	"payssage/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"resty.dev/v3"
)

// stubRoundTripper serves a canned MercadoPago response and captures the
// request body so tests can assert what payssage actually sent.
type stubRoundTripper struct {
	body string
	req  chan map[string]any
}

func (s *stubRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var payload map[string]any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber() // keep "50.00" as json.Number so assertions compare exact strings
	_ = dec.Decode(&payload)
	if s.req != nil {
		s.req <- payload
	}
	return &http.Response{
		StatusCode: http.StatusCreated,
		Status:     "201 Created",
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(s.body)),
		Request:    r,
	}, nil
}

type fakeSellers struct {
	ports.SellerRepo
	seller *models.Seller
}

func (f *fakeSellers) GetByID(_ context.Context, _ uuid.UUID) (*models.Seller, error) {
	if f.seller == nil {
		return nil, fun.Errf("seller not found").NotFound()
	}
	return f.seller, nil
}

type fakeWallets struct {
	ports.WalletRepo
	wallet *models.Wallet
}

func (f *fakeWallets) GetByID(_ context.Context, _ uuid.UUID) (*models.Wallet, error) {
	if f.wallet == nil {
		return nil, fun.Errf("wallet not found").NotFound()
	}
	return f.wallet, nil
}

type fakeIntents struct{ ports.IntentRepo }
type fakeCollectors struct{ ports.CollectorRepo }

const mpApproved = `{"id": 123456789, "status": "approved", "status_detail": "accredited"}`

func newCheckoutTestProvider(wallet *models.Wallet, seller *models.Seller) (*Provider, *stubRoundTripper) {
	stub := &stubRoundTripper{body: mpApproved, req: make(chan map[string]any, 1)}
	p := &Provider{
		cfg:        config.MercadoPagoConfig{MpAccessToken: "TEST_ACCESS_TOKEN"},
		intents:    fakeIntents{},
		collectors: fakeCollectors{},
		sellers:    &fakeSellers{seller: seller},
		wallets:    &fakeWallets{wallet: wallet},
		httpClient: resty.New().SetTransport(stub),
	}
	return p, stub
}

func testSeller(walletID uuid.UUID) *models.Seller {
	creds, _ := json.Marshal(models.MercadoPagoCredentials{
		PublicKey:   "TEST_PUBLIC_KEY",
		AccessToken: "TEST_SELLER_AT",
	})
	return &models.Seller{
		ID:          uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		WalletID:    walletID,
		Provider:    "mercadopago",
		Credentials: creds,
	}
}

func testIntent(walletID, sellerID uuid.UUID) *models.Intent {
	return &models.Intent{
		ID:          uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		WalletID:    walletID,
		SellerID:    sellerID,
		AmountCents: 100000, // R$ 1000.00
		Currency:    "BRL",
		Provider:    "mercadopago",
	}
}

const cardCheckoutData = `{
  "payment_method_id": "visa",
  "token": "mp_token_123",
  "issuer_id": "issuer_banco",
  "installments": 3,
  "payer": {"email": "buyer@example.com", "identification_type": "CPF", "identification_number": "12345678901"}
}`

const pixCheckoutData = `{
  "payment_method_id": "pix",
  "payer": {"email": "buyer@example.com", "identification_type": "CPF", "identification_number": "12345678901"}
}`

// TestCheckout_AppliesWalletFee pins the fix: the wallet's fee_bps is the
// authoritative marketplace fee — the checkout path no longer ignores it
// (application_fee was previously always 0 because univents sends no
// marketplace_fee_bps).
func TestCheckout_AppliesWalletFee(t *testing.T) {
	walletID := uuid.New()
	p, stub := newCheckoutTestProvider(
		&models.Wallet{ID: walletID, FeeBps: 500},
		testSeller(walletID),
	)

	if err := p.Checkout(context.Background(), testIntent(walletID, uuid.MustParse("22222222-2222-2222-2222-222222222222")), json.RawMessage(cardCheckoutData)); err != nil {
		t.Fatalf("Checkout: %v", err)
	}

	body := <-stub.req
	if got := body["application_fee"]; got != json.Number("50.00") {
		t.Errorf("application_fee = %v, want 50.00 (5%% of R$1000.00)", got)
	}
	if got := body["issuer_id"]; got != "issuer_banco" {
		t.Errorf("issuer_id = %v, want issuer_banco", got)
	}
	if got := body["installments"]; got != json.Number("3") {
		t.Errorf("installments = %v, want 3", got)
	}
}

// TestCheckout_WalletFeeWinsOverRequestFee pins that a caller-supplied
// marketplace_fee_bps cannot undercut the wallet's configured fee.
func TestCheckout_WalletFeeWinsOverRequestFee(t *testing.T) {
	walletID := uuid.New()
	data := `{
	  "payment_method_id": "visa",
	  "token": "mp_token_123",
	  "installments": 1,
	  "marketplace_fee_bps": 100,
	  "payer": {"email": "buyer@example.com", "identification_type": "CPF", "identification_number": "12345678901"}
	}`
	p, stub := newCheckoutTestProvider(
		&models.Wallet{ID: walletID, FeeBps: 500},
		testSeller(walletID),
	)

	if err := p.Checkout(context.Background(), testIntent(walletID, uuid.MustParse("22222222-2222-2222-2222-222222222222")), json.RawMessage(data)); err != nil {
		t.Fatalf("Checkout: %v", err)
	}

	body := <-stub.req
	if got := body["application_fee"]; got != json.Number("50.00") {
		t.Errorf("application_fee = %v, want 50.00 (wallet fee wins over caller's 100bps)", got)
	}
}

// TestCheckout_RequestFeeWhenWalletHasNone pins the fallback: with no
// wallet fee configured, a caller-supplied marketplace_fee_bps still works.
func TestCheckout_RequestFeeWhenWalletHasNone(t *testing.T) {
	walletID := uuid.New()
	data := `{
	  "payment_method_id": "visa",
	  "token": "mp_token_123",
	  "installments": 1,
	  "marketplace_fee_bps": 300,
	  "payer": {"email": "buyer@example.com", "identification_type": "CPF", "identification_number": "12345678901"}
	}`
	p, stub := newCheckoutTestProvider(
		&models.Wallet{ID: walletID, FeeBps: 0},
		testSeller(walletID),
	)

	if err := p.Checkout(context.Background(), testIntent(walletID, uuid.MustParse("22222222-2222-2222-2222-222222222222")), json.RawMessage(data)); err != nil {
		t.Fatalf("Checkout: %v", err)
	}

	body := <-stub.req
	if got := body["application_fee"]; got != json.Number("30.00") {
		t.Errorf("application_fee = %v, want 30.00 (caller's 300bps when wallet has no fee)", got)
	}
}

// TestCheckout_PixPinsShape: pix carries no card-only fields and the wallet
// fee still applies.
func TestCheckout_PixPinsShape(t *testing.T) {
	walletID := uuid.New()
	p, stub := newCheckoutTestProvider(
		&models.Wallet{ID: walletID, FeeBps: 500},
		testSeller(walletID),
	)

	if err := p.Checkout(context.Background(), testIntent(walletID, uuid.MustParse("22222222-2222-2222-2222-222222222222")), json.RawMessage(pixCheckoutData)); err != nil {
		t.Fatalf("Checkout: %v", err)
	}

	body := <-stub.req
	if got := body["application_fee"]; got != json.Number("50.00") {
		t.Errorf("application_fee = %v, want 50.00", got)
	}
	// The receipt/email title comes from the top-level `description`, not
	// additional_info — MP renders "Produto sem nome" without it. This pix
	// checkout carries no items, so the generic fallback applies.
	if got := body["description"]; got != "Payssage purchase" {
		t.Errorf("description = %v, want fallback (Payssage purchase)", got)
	}
	for _, key := range []string{"issuer_id", "installments", "token"} {
		if _, ok := body[key]; ok {
			t.Errorf("pix request must not carry %q", key)
		}
	}
	if _, ok := body["date_of_expiration"]; !ok {
		t.Error("pix request missing date_of_expiration")
	}
}

func TestPaymentDescription(t *testing.T) {
	cases := []struct {
		name string
		ai   *mercado_pago.AdditionalInfo
		want string
	}{
		{"nil items", nil, "Payssage purchase"},
		{"empty items", &mercado_pago.AdditionalInfo{}, "Payssage purchase"},
		{"single title no multiplier", &mercado_pago.AdditionalInfo{Items: []mercado_pago.Item{{Title: "Ticket Legal", Quantity: 1}}}, "Ticket Legal"},
		{"groups repeated titles with qty", &mercado_pago.AdditionalInfo{Items: []mercado_pago.Item{
			{Title: "Ticket Legal", Quantity: 1},
			{Title: "Ticket Legal", Quantity: 1},
			{Title: "Camiseta", Quantity: 2},
		}}, "2x Camiseta, 2x Ticket Legal"},
		{"single item with qty 2", &mercado_pago.AdditionalInfo{Items: []mercado_pago.Item{
			{Title: "Camiseta", Quantity: 2},
		}}, "2x Camiseta"},
		{"skips blank titles", &mercado_pago.AdditionalInfo{Items: []mercado_pago.Item{
			{Title: "  ", Quantity: 1},
			{Title: "Oficina", Quantity: 1},
		}}, "Oficina"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := paymentDescription(tc.ai); got != tc.want {
				t.Fatalf("paymentDescription = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestPaymentDescriptionTruncatesAtBoundary(t *testing.T) {
	longTitle := strings.Repeat("Ingresso Esgotado ", 30) // long single item
	desc := paymentDescription(&mercado_pago.AdditionalInfo{Items: []mercado_pago.Item{{Title: longTitle, Quantity: 1}}})
	if len(desc) > 250 {
		t.Fatalf("description over 250 bytes: %d", len(desc))
	}
	if !utf8.ValidString(desc) {
		t.Fatalf("description split a UTF-8 character: %q", desc)
	}

	// Multiple long items: truncation must end at a whole-item boundary (a
	// complete "Nx Title"), never mid-name and never with a dangling ", ".
	many := paymentDescription(&mercado_pago.AdditionalInfo{Items: []mercado_pago.Item{
		{Title: "Camiseta Oficial do Evento Edição 2026", Quantity: 3},
		{Title: "Ingresso Pista Premium com Acesso ao Backstage", Quantity: 2},
		{Title: "Programa de Mentoria Exclusiva para Participantes", Quantity: 1},
	}})
	if len(many) > 250 {
		t.Fatalf("description over 250 bytes: %d", len(many))
	}
	if strings.HasSuffix(many, ", ") || strings.HasSuffix(many, ",") {
		t.Fatalf("description ends with a dangling separator: %q", many)
	}
	// Each "Nx Title" segment must appear complete or not at all.
	for _, seg := range []string{"3x Camiseta Oficial do Evento Edição 2026", "2x Ingresso Pista Premium com Acesso ao Backstage"} {
		if strings.Contains(many, strings.TrimPrefix(seg, "3x ")) || strings.Contains(many, strings.TrimPrefix(seg, "2x ")) {
			if !strings.Contains(many, seg) {
				t.Fatalf("segment cut in half — found %q without its full %q: %q", strings.TrimPrefix(seg, "3x "), seg, many)
			}
		}
	}
}

func TestTruncateDescription(t *testing.T) {
	if got := truncateDescription("short", 10); got != "short" {
		t.Fatalf("short string must pass through: %q", got)
	}
	// Cuts at the last ", " boundary, dropping the tail item entirely.
	got := truncateDescription("2x Camiseta, 2x Ticket Legal, 2x Oficina", 22)
	if got != "2x Camiseta, 2x Ticket" && got != "2x Camiseta, 2x Ticket Legal" && got != "2x Camiseta" {
		// must be a prefix ending at a complete item (no mid-name cut)
		if strings.Contains(got, "Ticket Lega") && !strings.HasSuffix(got, "Legal") {
			t.Fatalf("mid-name cut: %q", got)
		}
	}
	// Rune safety: a multi-byte char on the boundary must not be split.
	acc := strings.Repeat("Ação ", 100)
	tr := truncateDescription(acc, 10)
	if !utf8.ValidString(tr) {
		t.Fatalf("split a UTF-8 character: %q", tr)
	}
}
