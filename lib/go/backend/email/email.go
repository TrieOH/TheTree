// Package email is the kit module for outgoing email over SMTP.
package email

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lib/backend"
	libconfig "lib/config"
	libemail "lib/email"
)

// Module fills Deps.Email with an SMTP sender. Nothing connects at boot:
// the sender dials per message. With SMTP_MONTHLY_LIMIT set, every send
// first reserves its recipients in the service's email_monthly_sends table,
// so it needs the postgres module before it.
func Module(cfg libconfig.SMTP) backend.Module {
	return backend.Module{
		Name: "email",
		Open: func(_ context.Context, d *backend.Deps) (backend.Close, error) {
			client := libemail.NewClient(libconfig.EmailConfig(cfg))
			if cfg.MonthlyLimit > 0 {
				if d.DB == nil {
					return nil, errors.New("email: SMTP_MONTHLY_LIMIT needs the postgres module first")
				}
				client.WithQuota(monthlyQuota{db: d.DB, limit: cfg.MonthlyLimit})
			}
			d.Email = client
			return backend.NoClose, nil
		},
	}
}

// monthlyQuota counts recipients per calendar month in one row. The update
// only applies while the total stays within the limit, so concurrent sends
// can't overshoot it.
type monthlyQuota struct {
	db    *pgxpool.Pool
	limit int
}

const reserveSQL = `
INSERT INTO email_monthly_sends AS s (month, sent)
VALUES (date_trunc('month', now())::date, $1)
ON CONFLICT (month) DO UPDATE SET sent = s.sent + $1
WHERE s.sent + $1 <= $2
RETURNING sent`

func (q monthlyQuota) Reserve(ctx context.Context, n int) error {
	if n > q.limit {
		return libemail.ErrMonthlyLimit
	}
	var sent int
	err := q.db.QueryRow(ctx, reserveSQL, n, q.limit).Scan(&sent)
	if errors.Is(err, pgx.ErrNoRows) {
		return libemail.ErrMonthlyLimit
	}
	if err != nil {
		return fmt.Errorf("email: reserve monthly quota: %w", err)
	}
	return nil
}
