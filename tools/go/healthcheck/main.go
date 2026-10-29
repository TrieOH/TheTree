// Command healthcheck is a static, dependency-free HTTP health probe for the
// distroless runtime images (which have no shell, no wget, no curl). It GETs
// the URL given as the first argument (default http://127.0.0.1:8080/health)
// and exits 0 on any 2xx response, 1 otherwise. It is the command behind the
// HEALTHCHECK instruction in api/<svc>/Dockerfile.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	os.Exit(run(os.Args))
}

// run performs the probe and returns the process exit code; it is separate
// from main so deferred cleanups run before os.Exit.
func run(args []string) int {
	url := "http://127.0.0.1:8080/health"
	if len(args) > 1 {
		url = args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil) //nolint:gosec // CLI probe invoked by the Dockerfile HEALTHCHECK; the URL is operator-provided, not untrusted input
	if err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck:", err)
		return 1
	}

	resp, err := http.DefaultClient.Do(req) //nolint:gosec // ditto — tainted-URL analysis only; see comment above
	if err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck:", err)
		return 1
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "healthcheck: %s returned %s\n", url, resp.Status)
		return 1
	}
	return 0
}
