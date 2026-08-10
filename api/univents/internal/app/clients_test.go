package app

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	payssage "sdk/payssage"

	"github.com/google/uuid"
)

// TestVerifyPayssageWallet_Resolves pins the happy path of the D6 boot
// check: a wallet that resolves must not stop the boot.
func TestVerifyPayssageWallet_Resolves(t *testing.T) {
	walletID := uuid.New()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"data":{
			"id":"` + walletID.String() + `","owner_id":"` + uuid.New().String() + `",
			"organization_id":null,"name":"platform","sandbox":false,
			"fee_bps":500,"collector_id":null,"created_at":"2026-08-01T00:00:00Z"
		},"timestamp":"2026-08-01T00:00:00Z"}`))
	}))
	t.Cleanup(srv.Close)

	client := payssage.New(payssage.Config{BaseURL: srv.URL, APIKey: "test"})

	// Must not exit — the resolved wallet passes the boot gate.
	VerifyPayssageWallet(context.Background(), client, walletID)
}

// TestVerifyPayssageWallet_ExitsOnUnresolvableWallet pins the fail-fast
// side: a wallet that does not resolve (wrong id / unreachable Payssage)
// must stop the boot. errx.Exit calls os.Exit(1), so the assertion runs in
// a subprocess (the standard pattern for testing exits).
func TestVerifyPayssageWallet_ExitsOnUnresolvableWallet(t *testing.T) {
	if os.Getenv("TEST_BOOT_CHECK_EXIT") == "1" {
		// unreachable Payssage: the client's HTTP call fails
		client := payssage.New(payssage.Config{BaseURL: "http://127.0.0.1:1", APIKey: "test"})
		VerifyPayssageWallet(context.Background(), client, uuid.New())
		return
	}

	cmd := exec.CommandContext(context.Background(), os.Args[0], "-test.run=TestVerifyPayssageWallet_ExitsOnUnresolvableWallet") //nolint:gosec // fixed args, standard subprocess-exit test pattern
	cmd.Env = append(os.Environ(), "TEST_BOOT_CHECK_EXIT=1")
	err := cmd.Run()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("want exit status 1 from the boot check, got %v", err)
	}
	if code := exitErr.ExitCode(); code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
}
