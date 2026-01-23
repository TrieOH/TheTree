package outbounds

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

type ErrServiceUnavailable struct {
	ServiceName string
}

func (e ErrServiceUnavailable) Error() string {
	return "service unavailable: " + e.ServiceName
}
