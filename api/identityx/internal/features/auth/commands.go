package auth

import (
	"IdentityX/contracts"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/feature_deps"
	"IdentityX/internal/shared/ports"
	"IdentityX/internal/shared/security"
	"context"
	"errors"
	"lib/database"
	"lib/errx"
	"strings"
	"time"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type CommandService struct {
	encryptionKey []byte
	issuer        string
	users         ports.UserRepository
	sessions      ports.SessionRepository
	projects      ports.ProjectRepository
	keys          ports.KeysRepository
	mailRenderer  ports.EmailRenderer
	mailSender    ports.Mailer
	logger        *zap.Logger
	tracer        trace.Tracer
	tx            database.TxRunner
}

func NewCommandService(deps feature_deps.AuthCommandDeps) *CommandService {
	return errx.MustProvide(&CommandService{
		encryptionKey: deps.EncryptionKey,
		issuer:        deps.Issuer,
		users:         deps.Users,
		sessions:      deps.Sessions,
		projects:      deps.Projects,
		keys:          deps.Keys,
		mailRenderer:  deps.Renderer,
		mailSender:    deps.Mailer,
		logger:        deps.Logger,
		tracer:        deps.Tracer,
		tx:            deps.Tx,
	})
}

// Register handles the business logic for creating a new user.
// It validates the input, hashes the password, and then attempts to create the user in the database.
// It returns an error if the email is already in use or if there is a problem with the database.
func (uc *CommandService) Register(ctx context.Context, in contracts.RegisterInput) error {
	var err error
	ctx, span := uc.tracer.Start(ctx, "AuthService.Register")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("register.success", err == nil))
		}
	}()

	var verificationEmail ports.Email
	err = uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		verificationEmail, err = uc.registerInternal(ctx, in)
		return err
	})
	if err != nil {
		return err
	}

	return uc.mailSender.Send(ctx, verificationEmail)
}

func (uc *CommandService) registerInternal(ctx context.Context, in contracts.RegisterInput) (ports.Email, error) {
	var err error
	var spanAttrs []attribute.KeyValue
	if in.ProjectID != nil {
		spanAttrs = append(spanAttrs, attribute.String("project.id", in.ProjectID.String()))
	}

	ctx, span := uc.tracer.Start(ctx, "AuthService.registerInternal", trace.WithAttributes(spanAttrs...))
	defer span.End()

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	if len(in.Password) > 72 {
		return ports.Email{}, fun.ErrBadRequest("password can't be longer than 72 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return ports.Email{}, fun.Errf("invalid password: %s", err).BadRequest()
	}

	userType := contracts.UserTypeClient
	if in.ProjectID != nil {
		userType = contracts.UserTypeProject
	}

	u, err := uc.users.Register(ctx, in.Email, string(hashedPassword), in.ProjectID, userType)
	if err != nil {
		return ports.Email{}, err
	}

	span.SetAttributes(
		attribute.String("user.id", u.ID.String()),
		attribute.Int64("user.created_at", u.CreatedAt.Unix()),
		attribute.String("user.type", string(u.UserType)),
	)

	kid, err := uc.keys.GetActiveSigningKID(ctx, nil)
	if err != nil {
		return ports.Email{}, err
	}

	verificationPayload, err := security.NewVerificationToken(contracts.NewVerificationTokenInput{
		KID:       kid,
		Subject:   u.ID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}, uc.issuer)
	if err != nil {
		return ports.Email{}, err
	}

	var pair *contracts.Pair
	pair, err = uc.keys.GetActiveSigningKey(ctx, in.ProjectID)
	if err != nil {
		return ports.Email{}, err
	}

	verificationSig, err := security.SignKey(verificationPayload, pair, uc.encryptionKey)
	if err != nil {
		return ports.Email{}, err
	}

	verificationEmail, err := uc.mailRenderer.Verification(ctx, ports.VerificationEmailData{
		UserID: u.ID,
		Email:  u.Email,
		Token:  security.AssembleJWT(verificationPayload, verificationSig),
		Locale: "en",
	})
	if err != nil {
		return ports.Email{}, err
	}

	return verificationEmail, nil
}

// Login handles the business logic for logging in a user.
// It finds the user by email, compares the password, and if successful,
// creates a new session and returns a new set of access and refresh tokens.
func (uc *CommandService) Login(ctx context.Context, in contracts.LoginInput) (tokens *contracts.UserTokensOutput, err error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	var spanAttrs []attribute.KeyValue
	if in.ProjectID != nil {
		spanAttrs = append(spanAttrs, attribute.String("project.id", in.ProjectID.String()))
	}

	ctx, span := uc.tracer.Start(ctx, "AuthService.Login", trace.WithAttributes(spanAttrs...))
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("login.success", err == nil))
		}
	}()

	u, err := uc.users.GetUserByEmail(ctx, in.Email, in.ProjectID)
	if fun.Is(err, fun.CodeNotFound) {
		return nil, fun.ErrUnauthorized("invalid email or password")
	}
	if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("user.id", u.ID.String()),
		attribute.String("user.type", string(u.UserType)),
		attribute.Int64("user.created_at_unix", u.CreatedAt.Unix()),
	)

	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return nil, fun.ErrUnauthorized("invalid email or password")
	}

	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	sess, err := uc.sessions.Create(ctx, contracts.Session{
		UserID:    u.ID,
		ProjectID: u.ProjectID,
		IssuedAt:  time.Now(),
		UserAgent: in.Agent,
		UserIP:    in.IP,
		ExpiresAt: refreshExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	if err = uc.users.UpdateLastLogin(ctx, u.ID); err != nil {
		return nil, err
	}

	accessJTI, err := uuid.NewV7()
	if err != nil {
		return nil, fun.ErrInternal("error generating UUIDv7 at auth/login")
	}

	kid, err := uc.keys.GetActiveSigningKID(ctx, u.ProjectID)
	if err != nil {
		return nil, err
	}

	pair, err := uc.keys.GetActiveSigningKey(ctx, u.ProjectID)
	if err != nil {
		return nil, err
	}

	accessExpiresAt := time.Now().Add(15 * time.Minute)

	accessPayload, err := security.NewAccessToken(contracts.NewAccessTokenInput{
		KID:       kid,
		User:      *u,
		IP:        in.IP,
		Agent:     in.Agent,
		AccessJTI: accessJTI.String(),
		SessionID: sess.SessionID,
		FamilyID:  sess.FamilyID,
		ExpiresAt: accessExpiresAt,
	}, uc.issuer)
	if err != nil {
		return nil, err
	}

	accessSig, err := security.SignKey(accessPayload, pair, uc.encryptionKey)
	if err != nil {
		return nil, err
	}

	refreshPayload, err := security.NewRefreshToken(contracts.NewRefreshTokenInput{
		KID:        kid,
		AccessJTI:  accessJTI,
		RefreshJTI: sess.TokenID,
		ExpiresAt:  refreshExpiresAt,
		FamilyID:   sess.FamilyID,
	}, uc.issuer)
	if err != nil {
		return nil, err
	}

	refreshSig, err := security.SignKey(refreshPayload, pair, uc.encryptionKey)
	if err != nil {
		return nil, err
	}

	domain := "https://dev.trieauth.trieoh.com"
	if u.ProjectID != nil {
		project, err := uc.projects.GetByIDInternal(ctx, *u.ProjectID)
		if err != nil {
			return nil, err
		}
		domain = project.Domain
	}

	return &contracts.UserTokensOutput{
		AccessTokenString:  security.AssembleJWT(accessPayload, accessSig),
		RefreshTokenString: security.AssembleJWT(refreshPayload, refreshSig),
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
		Domain:             domain,
	}, nil
}

// Logout handles the business logic for logging out a user.
// It retrieves the principal from the context, deletes the session, and revokes the refresh token.
func (uc *CommandService) Logout(ctx context.Context) error {
	ctx, span := uc.tracer.Start(ctx, "AuthService.Logout")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("logout.success", err == nil))
		}
	}()

	var principal *authz.Principal
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	if principal.Method == authz.AuthMethodApiKey {
		return errors.New("can't logout an api key, please revoke it instead")
	}

	var userType contracts.UserType
	if principal.ProjectID == nil {
		userType = contracts.UserTypeClient
	} else {
		userType = contracts.UserTypeProject
	}

	var sess *contracts.Session
	sess, err = uc.sessions.MarkRevokedByID(ctx, principal.UserID, *principal.SessionID, userType)
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.String("session.id", sess.SessionID.String()))

	return nil
}

// Refresh handles the business logic for refreshing a user's tokens.
// It parses the refresh token, checks if it's revoked, and if not,
// determines whether to refresh the tokens for a client or a project user.
func (uc *CommandService) Refresh(ctx context.Context, in contracts.RefreshInput) (*contracts.UserTokensOutput, error) {
	txOptions := database.TxOptions{
		Isolation: pgx.ReadCommitted,
		ReadOnly:  pgx.ReadWrite,
	}

	var out *contracts.UserTokensOutput
	err := uc.tx.WithinTxWithOptions(ctx, txOptions, func(ctx context.Context) error {
		var err error
		out, err = uc.refreshInternal(ctx, in)
		return err
	})

	return out, err
}

func (uc *CommandService) refreshInternal(ctx context.Context, in contracts.RefreshInput) (*contracts.UserTokensOutput, error) {
	ctx, span := uc.tracer.Start(ctx, "AuthService.Refresh")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("refresh.success", err == nil))
		}
	}()

	refreshToken := &contracts.RefreshClaims{}
	_, err = security.ParseJWTUnverified[*contracts.RefreshClaims](in.RefreshCookie, refreshToken)
	if err != nil {
		return nil, err
	}

	var oldJTI uuid.UUID
	oldJTI, err = uuid.Parse(refreshToken.ID)
	if err != nil {
		return nil, err
	}
	if oldJTI == uuid.Nil {
		return nil, fun.ErrBadRequest("invalid refresh token ID")
	}

	var uid uuid.UUID
	uid, err = uuid.NewV7()
	if err != nil {
		return nil, fun.ErrInternal("error generating UUIDv7 at auth/refreshInternal")
	}

	var newRefreshJTI = uid
	var refreshExp = time.Now().Add(7 * 24 * time.Hour)

	span.SetAttributes(attribute.String("old_token.id", oldJTI.String()))
	span.SetAttributes(attribute.String("new_token.id", newRefreshJTI.String()))

	var sess *contracts.Session
	sess, err = uc.sessions.GetByFamilyID(ctx, refreshToken.Sub.FamilyID)
	if err != nil {
		return nil, fun.ErrUnauthorized("session not found or revoked")
	}

	now := time.Now()
	if sess.ExpiresAt.Before(now) || sess.RevokedAt != nil {
		// FIXME Record suspicious behaviour on audit when it is implemented
		return nil, fun.ErrUnauthorized("session not found or revoked")
	}

	// should revoke the session because of replay attacks
	// FIXME Add suspicious behaviour to audit when it is implemented
	if sess.TokenID != oldJTI {
		_ = uc.sessions.MarkRevokedByFamilyID(ctx, sess.FamilyID)
		return nil, fun.ErrUnauthorized("refresh token reuse not allowed")
	}

	sess, err = uc.sessions.RotateToken(ctx, refreshToken.Sub.FamilyID, newRefreshJTI, oldJTI, refreshExp)
	if fun.Is(err, fun.CodeNotFound) {
		// sql.ErrNoRows → raced / reused / revoked
		return nil, fun.ErrUnauthorized("session not found or revoked")
	} else if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("session.id", sess.SessionID.String()),
		attribute.String("session.token_id", sess.TokenID.String()),
		attribute.String("session.user_type", string(sess.UserType)),
	)

	u, err := uc.users.GetUserByID(ctx, sess.UserID)
	if err != nil {
		return nil, err
	}

	newAccessJTI, err := uuid.NewV7()
	if err != nil {
		return nil, fun.ErrInternal("error generating UUIDv7 at auth/refreshInternal")
	}

	var (
		kid    string
		pair   *contracts.Pair
		domain string
	)

	kid, err = uc.keys.GetActiveSigningKID(ctx, sess.ProjectID)
	if err != nil {
		return nil, err
	}

	pair, err = uc.keys.GetActiveSigningKey(ctx, sess.ProjectID)
	if err != nil {
		return nil, err
	}

	domain = "https://dev.trieauth.trieoh.com"
	if sess.ProjectID != nil {
		project, err := uc.projects.GetByIDInternal(ctx, *sess.ProjectID)
		if err != nil {
			return nil, err
		}
		domain = project.Domain
	}

	// sign := func(payload []byte) ([]byte, error) — abstracts the GoAuth vs Project signer
	sign := func(payload []byte) ([]byte, error) {
		return security.SignKey(payload, pair, uc.encryptionKey)
	}

	accessExpiresAt := time.Now().Add(15 * time.Minute)

	var accessPayload []byte
	accessPayload, err = security.NewAccessToken(contracts.NewAccessTokenInput{
		KID: kid, User: *u, IP: in.IP, Agent: in.Agent,
		AccessJTI: newAccessJTI.String(), SessionID: sess.SessionID,
		FamilyID: sess.FamilyID, ExpiresAt: accessExpiresAt,
	}, uc.issuer)
	if err != nil {
		return nil, err
	}

	accessSig, err := sign(accessPayload)
	if err != nil {
		return nil, err
	}

	refreshPayload, err := security.NewRefreshToken(contracts.NewRefreshTokenInput{
		KID: kid, AccessJTI: newAccessJTI, RefreshJTI: newRefreshJTI,
		ExpiresAt: refreshExp, FamilyID: sess.FamilyID,
	}, uc.issuer)
	if err != nil {
		return nil, err
	}

	refreshSig, err := sign(refreshPayload)
	if err != nil {
		return nil, err
	}

	return &contracts.UserTokensOutput{
		AccessTokenString:  security.AssembleJWT(accessPayload, accessSig),
		RefreshTokenString: security.AssembleJWT(refreshPayload, refreshSig),
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExp,
		Domain:             domain,
	}, nil
}
