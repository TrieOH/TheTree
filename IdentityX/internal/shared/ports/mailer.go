package ports

import "context"

type Mailer interface {
	Send(ctx context.Context, email Email) error
}

type Email struct {
	To       string
	Subject  string
	TextBody string
	HTMLBody string
}
