package errx

import (
	"log/slog"
	"os"
)

func Must(err error, msg string) {
	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stderr, nil)).Error(msg, "err", err)
		os.Exit(1)
	}
}
