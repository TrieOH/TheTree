// Package email is the kit module for outgoing email over SMTP.
package email

import (
	"context"

	"lib/backend"
	libconfig "lib/config"
	libemail "lib/email"
)

// Module fills Deps.Email with an SMTP sender. Nothing connects at boot:
// the sender dials per message.
func Module(cfg libconfig.SMTP) backend.Module {
	return backend.Module{
		Name: "email",
		Open: func(_ context.Context, d *backend.Deps) (backend.Close, error) {
			d.Email = libemail.NewClient(libconfig.EmailConfig(cfg))
			return backend.NoClose, nil
		},
	}
}
