package outbounds

type VerificationEmailData struct {
	UserID string
	Email  string
	Token  string
	Locale string
}

type PasswordResetEmailData struct {
	UserID string
	Email  string
	Token  string
	Locale string
}

type EmailRenderer interface {
	Verification(data VerificationEmailData) (Email, error)
	PasswordReset(data PasswordResetEmailData) (Email, error)
}
