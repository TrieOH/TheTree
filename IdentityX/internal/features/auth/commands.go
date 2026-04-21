package auth

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/security"
	"IdentityX/internal/platform/telemetry"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"IdentityX/internal/shared/validation"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	fun "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type CommandService struct {
	users          ports.UserRepository
	sessions       ports.SessionRepository
	projects       ports.ProjectRepository
	keys           ports.KeysRepository
	tokenReuseList ports.TokenReuseListRepository
	redis          ports.RedisCacheService
	mailRenderer   ports.EmailRenderer
	mailSender     ports.Mailer
	logger         *zap.Logger
	tracer         trace.Tracer
	tx             database.TxRunner
}

func NewCommandService(
	Users ports.UserRepository,
	Sessions ports.SessionRepository,
	Projects ports.ProjectRepository,
	Keys ports.KeysRepository,
	TokenReuseList ports.TokenReuseListRepository,
	Redis ports.RedisCacheService,
	renderer ports.EmailRenderer,
	mailer ports.Mailer,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		users:          Users,
		sessions:       Sessions,
		projects:       Projects,
		keys:           Keys,
		tokenReuseList: TokenReuseList,
		redis:          Redis,
		mailRenderer:   renderer,
		mailSender:     mailer,
		logger:         logger,
		tracer:         tracer,
		tx:             tx,
	}
}

type RegisterInput struct {
	Email     string
	Password  string
	ProjectID *uuid.UUID // nil = client
}

type UserTokensOutput struct {
	AccessTokenString  string
	RefreshTokenString string

	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time

	Domain string
}

// Register handles the business logic for creating a new user.
// It validates the input, hashes the password, and then attempts to create the user in the database.
// It returns an error if the email is already in use or if there is a problem with the database.
func (uc *CommandService) Register(ctx context.Context, in RegisterInput) error {
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

func (uc *CommandService) registerInternal(ctx context.Context, in RegisterInput) (ports.Email, error) {
	var err error
	var spanAttrs []attribute.KeyValue
	if in.ProjectID != nil {
		spanAttrs = append(spanAttrs, attribute.String("project.id", in.ProjectID.String()))
	}

	ctx, span := uc.tracer.Start(ctx, "AuthService.registerInternal", trace.WithAttributes(spanAttrs...))
	defer span.End()

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	if len(in.Password) > 72 {
		return ports.Email{}, fun.NewError("password can't be longer than 72 characters").BadRequest()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return ports.Email{}, fun.NewError("invalid password").WithErr(err).BadRequest()
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
	})
	if err != nil {
		return ports.Email{}, err
	}

	var pair *contracts.Pair
	pair, err = uc.keys.GetActiveSigningKey(ctx, in.ProjectID)
	if err != nil {
		return ports.Email{}, err
	}

	verificationSig, err := security.SignKey(verificationPayload, pair)
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

type LoginInput struct {
	Email     string
	Password  string
	IP        string
	Agent     string
	ProjectID *uuid.UUID // nil = client
}

// Login handles the business logic for logging in a user.
// It finds the user by email, compares the password, and if successful,
// creates a new session and returns a new set of access and refresh security.
func (uc *CommandService) Login(ctx context.Context, in LoginInput) (tokens *UserTokensOutput, err error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	spanAttrs := []attribute.KeyValue{}
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
	if fail.Is(err, errx.SQLNotFound) {
		return nil, fail.New(errx.AuthInvalidCredentials).RecordCtx(ctx)
	} else if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("user.id", u.ID.String()),
		attribute.String("user.type", string(u.UserType)),
		attribute.Int64("user.created_at_unix", u.CreatedAt.Unix()),
	)

	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return nil, fail.New(errx.AuthInvalidCredentials).Trace(err.Error()).RecordCtx(ctx)
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
		return nil, fail.New(errx.SYSUUIDV7GenerationError).With(err).WithArgs("auth/login").RecordCtx(ctx)
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
	})
	if err != nil {
		return nil, err
	}

	accessSig, err := security.SignKey(accessPayload, pair)
	if err != nil {
		return nil, err
	}

	refreshPayload, err := security.NewRefreshToken(contracts.NewRefreshTokenInput{
		KID:        kid,
		AccessJTI:  accessJTI,
		RefreshJTI: sess.TokenID,
		ExpiresAt:  refreshExpiresAt,
		FamilyID:   sess.FamilyID,
	})
	if err != nil {
		return nil, err
	}

	refreshSig, err := security.SignKey(refreshPayload, pair)
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

	return &UserTokensOutput{
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

type RefreshInput struct {
	RefreshCookie *http.Cookie
	Agent         string
	IP            string
}

// Refresh handles the business logic for refreshing a user's security.
// It parses the refresh token, checks if it's revoked, and if not,
// determines whether to refresh the security for a client or a project user.
func (uc *CommandService) Refresh(ctx context.Context, in RefreshInput) (*UserTokensOutput, error) {
	txOptions := database.TxOptions{
		Isolation: pgx.ReadCommitted,
		ReadOnly:  pgx.ReadWrite,
	}

	var out *UserTokensOutput
	err := uc.tx.WithinTxWithOptions(ctx, txOptions, func(ctx context.Context) error {
		var err error
		out, err = uc.refreshInternal(ctx, in)
		return err
	})

	return out, err
}

func (uc *CommandService) refreshInternal(ctx context.Context, in RefreshInput) (*UserTokensOutput, error) {
	ctx, span := uc.tracer.Start(ctx, "AuthService.Refresh")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("refresh.success", err == nil))
		}
	}()

	refreshToken := &contracts.RefreshClaims{}
	_, err = security.ParseJWTUnverified[*contracts.RefreshClaims](in.RefreshCookie.Value, refreshToken)
	if err != nil {
		return nil, err
	}

	var oldJTI uuid.UUID
	oldJTI, err = validation.RequireRefreshJTI(&refreshToken.ID)
	if err != nil {
		return nil, err
	}

	var uid uuid.UUID
	uid, err = uuid.NewV7()
	if err != nil {
		return nil, fail.New(errx.SYSUUIDV7GenerationError).With(err).WithArgs("auth/refreshInternal").RecordCtx(ctx)
	}

	var newRefreshJTI = uid
	var refreshExp = time.Now().Add(7 * 24 * time.Hour)

	span.SetAttributes(attribute.String("old_token.id", oldJTI.String()))
	span.SetAttributes(attribute.String("new_token.id", newRefreshJTI.String()))

	var sess *contracts.Session
	sess, err = uc.sessions.GetByFamilyID(ctx, refreshToken.Sub.FamilyID)
	if err != nil {
		return nil, fail.New(errx.SessionNotFound).RecordCtx(ctx)
	}

	now := time.Now()
	if sess.ExpiresAt.Before(now) || sess.RevokedAt != nil {
		// FIXME Record suspicious behaviour on audit when it is implemented
		return nil, fail.New(errx.SessionNotFound).RecordCtx(ctx)
	}

	// should revoke the session because of replay attacks
	// FIXME Add suspicious behaviour to audit when it is implemented
	if sess.TokenID != oldJTI {
		_ = uc.sessions.MarkRevokedByFamilyID(ctx, sess.FamilyID)
		return nil, fail.New(errx.TokenReuseIdentified).WithArgs("refresh").RecordCtx(ctx)
	}

	sess, err = uc.sessions.RotateToken(ctx, refreshToken.Sub.FamilyID, newRefreshJTI, oldJTI, refreshExp)
	if fail.Is(err, errx.SQLNotFound) {
		// sql.ErrNoRows → raced / reused / revoked
		return nil, fail.New(errx.SessionNotFound).RecordCtx(ctx)
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
		return nil, fail.New(errx.SYSUUIDV7GenerationError).With(err).WithArgs("auth/finishRefresh").RecordCtx(ctx)
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
		return security.SignKey(payload, pair)
	}

	accessExpiresAt := time.Now().Add(15 * time.Minute)

	var accessPayload []byte
	accessPayload, err = security.NewAccessToken(contracts.NewAccessTokenInput{
		KID: kid, User: *u, IP: in.IP, Agent: in.Agent,
		AccessJTI: newAccessJTI.String(), SessionID: sess.SessionID,
		FamilyID: sess.FamilyID, ExpiresAt: accessExpiresAt,
	})
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
	})
	if err != nil {
		return nil, err
	}

	refreshSig, err := sign(refreshPayload)
	if err != nil {
		return nil, err
	}

	return &UserTokensOutput{
		AccessTokenString:  security.AssembleJWT(accessPayload, accessSig),
		RefreshTokenString: security.AssembleJWT(refreshPayload, refreshSig),
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExp,
		Domain:             domain,
	}, nil
}

func (uc *CommandService) GetJWKS(ctx context.Context) (map[string]any, error) {
	ekeys, err := uc.keys.ListPublicKeys(ctx, nil)
	if err != nil {
		telemetry.Log().Error("Failed listing GoAuth public security", zap.Error(err))
		return nil, fail.New(errx.SYSJWKSRetrievalFailed).With(err).RecordCtx(ctx)
	}

	jwkKeys := make([]any, 0, len(ekeys))
	var jwk map[string]any
	for _, k := range ekeys {
		jwk, err = contracts.PublicKeyToJWK(k)
		if err != nil {
			return nil, err
		}
		jwkKeys = append(jwkKeys, jwk)
	}

	return map[string]any{
		"security": jwkKeys,
	}, nil
}
