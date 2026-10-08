package main

import (
	"fmt"
	"os"
	"time"

	"IdentityX/internal/app"
	"IdentityX/internal/jobs"
	"lib/errx"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "schedules" {
		errx.Exit(printSchedules(), "identityx schedules")
		return
	}
	errx.Exit(app.Run(), "identityx exited with error")
}

// printSchedules lists the periodic jobs as "<kind> <seconds>" lines, for
// the managed-mode deploy to turn into scheduler entries. It reads only
// ROTATE_KEYS_JOB_DURATION, so it runs without the rest of the config.
func printSchedules() error {
	rotate, err := time.ParseDuration(os.Getenv("ROTATE_KEYS_JOB_DURATION"))
	if err != nil {
		return fmt.Errorf("ROTATE_KEYS_JOB_DURATION: %w", err)
	}
	for _, p := range jobs.Periodic(rotate) {
		fmt.Printf("%s %d\n", p.Job.Kind(), int(p.Every.Seconds()))
	}
	return nil
}
