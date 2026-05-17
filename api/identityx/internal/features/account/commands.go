package account

import (
	"IdentityX/contracts"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/feature_deps"
	"IdentityX/internal/shared/ports"
	"IdentityX/internal/shared/security"
	"context"
	"lib/database"
	"lib/errx"
	"time"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type CommandService struct {
	encryptionKey  []byte
	issuer         string
	users          ports.UserRepository
	accounts       ports.AccountRepository
	sessions       ports.SessionRepository
	keys           ports.KeysRepository
	tokenReuseList ports.TokenReuseListRepository
	mailRenderer   ports.EmailRenderer
	mailSender     ports.Mailer
	logger         *zap.Logger
	tracer         trace.Tracer
	tx             database.TxRunner
}

func NewCommandService(deps feature_deps.AccountCommandDeps) *CommandService {
	return errx.MustProvide(&CommandService{
		encryptionKey:  deps.EncryptionKey,
		issuer:         deps.Issuer,
		users:          deps.Users,
		accounts:       deps.Accounts,
		sessions:       deps.Sessions,
		keys:           deps.Keys,
		tokenReuseList: deps.TokenReuseList,
		mailRenderer:   deps.MailRenderer,
		mailSender:     deps.MailSender,
		logger:         deps.Logger,
		tracer:         deps.Tracer,
		tx:             deps.Tx,
	})
}

func (uc *CommandService) Verify(ctx context.Context, token string) (err error) {
	ctx, span := uc.tracer.Start(ctx, "AuthService.Verify")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("verify.success", err == nil))
		}
	}()

	var principal *authz.Principal
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	pair, err := uc.keys.GetActiveSigningKey(ctx, nil)
	if err != nil {
		return err
	}

	var claims *contracts.VerificationClaims
	claims, err = security.VerifyVerificationToken(token, pair)
	if err != nil {
		return err
	}
	if claims.Sub.Subject != principal.UserID {
		return fun.ErrUnauthorized("verification token user mismatch")
	}

	if principal.ProjectID != nil {
		span.SetAttributes(attribute.String("user.project_id", principal.ProjectID.String()))
	}
	span.SetAttributes(attribute.String("user.type", string(principal.UserType)))

	wasAlreadyVerified, err := uc.accounts.Verify(ctx, claims.Sub.Subject)
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.Bool("user.was_already_verified", wasAlreadyVerified))

	return nil
}

func (uc *CommandService) ResendVerificationEmail(ctx context.Context) (err error) {
	ctx, span := uc.tracer.Start(ctx, "AccountService.ResendVerificationEmail")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("resend_verification.success", err == nil))
		}
	}()

	var principal *authz.Principal
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var u *contracts.User
	var pu *contracts.User
	if principal.ProjectID != nil {
		pu, err = uc.users.GetUserByID(ctx, principal.UserID)
		if err != nil {
			return err
		}
		if pu.IsVerified == true {
			return fun.ErrBadRequest("user already verified")
		}
	} else {
		u, err = uc.users.GetUserByID(ctx, principal.UserID)
		if err != nil {
			return err
		}
		if u.IsVerified == true {
			return fun.ErrBadRequest("user already verified")
		}
	}

	if pu != nil {
		if pu.IsVerified {
			return fun.ErrBadRequest("user already verified")
		}
	} else {
		if u.IsVerified {
			return fun.ErrBadRequest("user already verified")
		}
	}

	var SigningKid string
	SigningKid, err = uc.keys.GetActiveSigningKID(ctx, nil)
	if err != nil {
		return err
	}

	var verificationPayload []byte
	verificationPayload, err = security.NewVerificationToken(contracts.NewVerificationTokenInput{
		KID:       SigningKid,
		Subject:   principal.UserID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}, uc.issuer)
	if err != nil {
		return err
	}

	pair, err := uc.keys.GetActiveSigningKey(ctx, nil)
	if err != nil {
		return err
	}

	var verificationSig []byte
	verificationSig, err = security.SignKey(verificationPayload, pair, uc.encryptionKey)
	if err != nil {
		return err
	}

	verificationTokenStr := security.AssembleJWT(verificationPayload, verificationSig)

	var verificationEmail ports.Email
	if pu != nil {
		verificationEmail, err = uc.mailRenderer.Verification(ctx, ports.VerificationEmailData{
			UserID: pu.ID,
			Email:  pu.Email,
			Token:  verificationTokenStr,
			Locale: "en",
		})
	} else {
		verificationEmail, err = uc.mailRenderer.Verification(ctx, ports.VerificationEmailData{
			UserID: u.ID,
			Email:  u.Email,
			Token:  verificationTokenStr,
			Locale: "en",
		})
	}

	if err != nil {
		return err
	}

	err = uc.mailSender.Send(ctx, verificationEmail)
	if err != nil {
		return err
	}

	return nil
}

func (uc *CommandService) ForgotPassword(ctx context.Context, in contracts.ForgotPasswordInput) (err error) {
	ctx, span := uc.tracer.Start(ctx, "AccountService.ForgotPassword")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("forgot_password.success", err == nil))
	}()

	var u *contracts.User
	u, err = uc.users.GetUserByEmail(ctx, in.Email, in.ProjectID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return nil // silent success (no enumeration)
		}
		return err
	}

	var SigningKid string
	SigningKid, err = uc.keys.GetActiveSigningKID(ctx, nil)
	if err != nil {
		return err
	}

	var resetPayload []byte
	resetPayload, err = security.NewResetPasswordToken(contracts.NewResetPasswordInput{
		KID:       SigningKid,
		Subject:   u.ID,
		ProjectID: in.ProjectID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}, uc.issuer)
	if err != nil {
		return err
	}

	pair, err := uc.keys.GetActiveSigningKey(ctx, nil)
	if err != nil {
		return err
	}

	var resetSig []byte
	resetSig, err = security.SignKey(resetPayload, pair, uc.encryptionKey)
	if err != nil {
		return err
	}

	resetPassTokenStr := security.AssembleJWT(resetPayload, resetSig)

	var e ports.Email
	e, err = uc.mailRenderer.PasswordReset(ctx, ports.PasswordResetEmailData{
		UserID: u.ID.String(),
		Email:  u.Email,
		Token:  resetPassTokenStr,
		Locale: "en",
	})

	err = uc.mailSender.Send(ctx, e)
	if err != nil {
		return err
	}
	return nil
}

func (uc *CommandService) ResetPassword(ctx context.Context, in contracts.ResetPasswordInput) (err error) {
	ctx, span := uc.tracer.Start(ctx, "AccountService.ResetPassword")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("reset_password.success", err == nil))
	}()

	err = uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		return uc.resetPasswordInternal(ctx, in)
	})

	return err
}

func (uc *CommandService) resetPasswordInternal(ctx context.Context, in contracts.ResetPasswordInput) (err error) {
	ctx, span := uc.tracer.Start(ctx, "AccountService.resetPasswordInternal")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("reset_password.success", err == nil))
	}()

	pair, err := uc.keys.GetActiveSigningKey(ctx, nil)
	if err != nil {
		return err
	}

	var claims *contracts.ResetPasswordClaims
	claims, err = security.VerifyResetPasswordToken(in.Token, pair)
	if err != nil {
		return err
	}

	var jti uuid.UUID
	jti, err = uuid.Parse(claims.ID)
	if err != nil {
		return fun.Errf("invalid jti in claims: %s", err).BadRequest()
	}

	var exists bool
	exists, err = uc.tokenReuseList.Exists(ctx, jti, claims.Sub.Subject)
	if err != nil {
		return err
	}
	if exists {
		// FIXME when the audit is implemented add this to the audit
		return fun.ErrUnauthorized("token already used")
	}
	if len(in.NewPassword) > 72 {
		return fun.ErrBadRequest("password length exceeds 72 bytes")
	}

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fun.Errf("invalid password: %s", err).BadRequest()
	}

	if claims.Sub.ProjectID == nil {
		err = uc.accounts.ResetPassword(ctx, claims.Sub.Subject, hashedPassword)
		if err != nil {
			return err
		}
		_, err = uc.sessions.MarkRevokedByFilter(ctx, contracts.Filter{
			UserID:   claims.Sub.Subject,
			UserType: contracts.UserTypeClient,
		})
		if err != nil {
			return err
		}
	} else {
		err = uc.accounts.ResetPassword(ctx, claims.Sub.Subject, hashedPassword)
		if err != nil {
			return err
		}
		_, err = uc.sessions.MarkRevokedByFilter(ctx, contracts.Filter{
			UserID:   claims.Sub.Subject,
			UserType: contracts.UserTypeProject,
		})
		if err != nil {
			return err
		}
	}

	err = uc.tokenReuseList.Append(ctx, jti, claims.Sub.Subject, claims.ExpiresAt.Time)
	if err != nil {
		return err
	}

	return nil
}
