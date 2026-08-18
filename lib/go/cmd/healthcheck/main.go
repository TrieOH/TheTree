// Command healthcheck is a static, dependency-free HTTP health probe for the
// distroless runtime images (which have no shell, no wget, no curl). It GETs
// the URL given as the first argument (default http://127.0.0.1:8080/health)
// and exits 0 on any 2xx response, 1 otherwise. It is the command behind the
// HEALTHCHECK instruction in api/<svc>/Dockerfile.
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	url := "http://127.0.0.1:8080/health"
	if len(os.Args) > 1 {
		url = os.Args[1]
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "healthcheck:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "healthcheck: %s returned %s\n", url, resp.Status)
		os.Exit(1)
	}
}
