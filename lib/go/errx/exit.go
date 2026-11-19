package errx

import (
	"log/slog"
	"os"
)

// Exit logs err and ends the process when err is set. It belongs at the
// process edge (cmd/main, config loading); boot code returns errors (ADR-0014).
func Exit(err error, msg string) {
	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stderr, nil)).Error(msg, "err", err)
		os.Exit(1)
	}
}
