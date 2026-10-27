// Package email is the kit module for outgoing email over SMTP.
package email

import (
	"context"

	"lib/backend"
	libemail "lib/email"
)

// Module fills Deps.Email with an SMTP sender. Nothing connects at boot:
// the sender dials per message.
func Module(cfg libemail.Config) backend.Module {
	return backend.Module{
		Name: "email",
		Open: func(_ context.Context, d *backend.Deps) (backend.Close, error) {
			d.Email = libemail.NewClient(cfg)
			return backend.NoClose, nil
		},
	}
}
