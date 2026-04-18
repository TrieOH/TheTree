package users

import (
	"IdentityX/internal/features/keys"
	"IdentityX/internal/features/tokens"
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/telemetry"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"IdentityX/internal/shared/validation"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

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
	projectUsers   ports.ProjectUserRepository
	keys           ports.KeysRepository
	tokenReuseList ports.TokenReuseListRepository
	redis          ports.RedisCacheService
	keysS          keys.CommandService
	tokens         tokens.CommandService
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
	ProjectUsers ports.ProjectUserRepository,
	Keys ports.KeysRepository,
	TokenReuseList ports.TokenReuseListRepository,
	Redis ports.RedisCacheService,
	keysS keys.CommandService,
	tokenBundle tokens.CommandService,
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
		projectUsers:   ProjectUsers,
		keys:           Keys,
		tokenReuseList: TokenReuseList,
		redis:          Redis,
		keysS:          keysS,
		tokens:         tokenBundle,
		mailRenderer:   renderer,
		mailSender:     mailer,
		logger:         logger,
		tracer:         tracer,
		tx:             tx,
	}
}

type RegisterUserInput struct {
	Email    string
	Password string
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
func (uc *CommandService) Register(ctx context.Context, in RegisterUserInput) error {
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

	err = uc.mailSender.Send(ctx, verificationEmail)
	if err != nil {
		return err
	}

	return nil
}

func (uc *CommandService) registerInternal(ctx context.Context, in RegisterUserInput) (ports.Email, error) {
	var err error
	ctx, span := uc.tracer.Start(ctx, "AuthService.registerInternal")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("register.success", err == nil))
		}
	}()

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	if len(in.Password) > 72 {
		return ports.Email{}, fail.New(errx.AuthInvalidPassword).RecordCtx(ctx)
	}

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return ports.Email{}, fail.New(errx.RequestInvalidPassword).With(err).RecordCtx(ctx)
	}

	var u *contracts.User
	u, err = uc.users.Register(ctx, in.Email, string(hashedPassword))
	if err != nil {
		return ports.Email{}, err
	}

	span.SetAttributes(
		attribute.String("user.id", u.ID.String()),
		attribute.Int64("user.created_at", u.CreatedAt.Unix()),
		attribute.String("user.type", u.UserType),
	)

	var identity *contracts.Identity
	identity, err = uc.sessions.CreateIdentity(ctx, contracts.ClientIdentity, u.ID)
	if err != nil {
		return ports.Email{}, err
	}

	span.SetAttributes(
		attribute.String("user.identity.id", identity.ID.String()),
		attribute.String("user.identity.type", string(identity.IdentityType)),
	)

	var SigningKid string
	SigningKid, err = uc.keys.GetActiveGoAuthSigningKID(ctx)
	if err != nil {
		return ports.Email{}, err
	}

	var verificationPayload []byte
	verificationPayload, err = uc.tokens.NewVerificationToken(contracts.NewVerificationTokenInput{
		KID:       SigningKid,
		Subject:   u.ID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return ports.Email{}, err
	}

	var verificationSig []byte
	verificationSig, err = uc.keysS.SignGoAuth(ctx, verificationPayload)
	if err != nil {
		return ports.Email{}, err
	}

	verificationTokenStr := uc.tokens.AssembleJWT(verificationPayload, verificationSig)

	var verificationEmail ports.Email
	verificationEmail, err = uc.mailRenderer.Verification(ctx, ports.VerificationEmailData{
		UserID: u.ID,
		Email:  u.Email,
		Token:  verificationTokenStr,
		Locale: "en",
	})
	if err != nil {
		return ports.Email{}, err
	}

	return verificationEmail, nil
}

type LoginUserInput struct {
	Email    string
	Password string

	Agent string
	IP    string
}

// Login handles the business logic for logging in a user.
// It finds the user by email, compares the password, and if successful,
// creates a new session and returns a new set of access and refresh tokens.
func (uc *CommandService) Login(ctx context.Context, in LoginUserInput) (tokens *UserTokensOutput, err error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	ctx, span := uc.tracer.Start(ctx, "AuthService.Login")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("login.success", err == nil))
		}
	}()

	var u *contracts.User
	u, err = uc.users.GetUserByEmail(ctx, in.Email)
	if fail.Is(err, errx.SQLNotFound) {
		return nil, fail.New(errx.AuthInvalidCredentials).RecordCtx(ctx)
	} else if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("user.id", u.ID.String()),
		attribute.String("user.type", u.UserType),
		attribute.Int64("user.created_at_unix", u.CreatedAt.Unix()),
	)

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password))
	if err != nil {
		return nil, fail.New(errx.AuthInvalidCredentials).Trace(err.Error()).RecordCtx(ctx)
	}

	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	var identity *contracts.Identity
	identity, err = uc.sessions.GetIdentityByEntityIDAndType(ctx, u.ID, contracts.ClientIdentity)
	if err != nil {
		return nil, err
	}

	var sess *contracts.Session
	sess, err = uc.sessions.Create(ctx, contracts.Session{
		IdentityID: identity.ID,
		IssuedAt:   time.Now(),
		UserAgent:  in.Agent,
		UserIP:     in.IP,
		ExpiresAt:  refreshExpiresAt,
	})

	if err != nil {
		return nil, err
	}

	var accessJTI uuid.UUID
	accessJTI, err = uuid.NewV7()
	if err != nil {
		return nil, fail.New(errx.SYSUUIDV7GenerationError).With(err).WithArgs("auth/login").RecordCtx(ctx)
	}

	var SigningKid string
	SigningKid, err = uc.keys.GetActiveGoAuthSigningKID(ctx)
	if err != nil {
		return nil, err
	}

	accessExpiresAt := time.Now().Add(15 * time.Minute)
	var accessPayload []byte
	accessPayload, err = uc.tokens.NewAccessToken(contracts.NewAccessTokenInput{
		KID:       SigningKid,
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

	var accessSig []byte
	accessSig, err = uc.keysS.SignGoAuth(ctx, accessPayload)
	if err != nil {
		return nil, err
	}

	accessTokenStr := uc.tokens.AssembleJWT(accessPayload, accessSig)

	var refreshPayload []byte
	refreshPayload, err = uc.tokens.NewRefreshToken(contracts.NewRefreshTokenInput{
		KID:        SigningKid,
		AccessJTI:  accessJTI,
		RefreshJTI: sess.TokenID,
		ExpiresAt:  refreshExpiresAt,
		FamilyID:   sess.FamilyID,
	})
	if err != nil {
		return nil, err
	}

	var refreshSig []byte
	refreshSig, err = uc.keysS.SignGoAuth(ctx, refreshPayload)
	if err != nil {
		return nil, err
	}

	refreshTokenStr := uc.tokens.AssembleJWT(refreshPayload, refreshSig)

	return &UserTokensOutput{
		AccessTokenString:  accessTokenStr,
		RefreshTokenString: refreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
		Domain:             "https://dev.trieauth.trieoh.com",
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

	var identityType contracts.IdentityType
	if principal.ProjectID == nil {
		identityType = contracts.ClientIdentity
	} else {
		identityType = contracts.ProjectIdentity
	}

	var sess *contracts.Session
	sess, err = uc.sessions.MarkRevokedByID(ctx, principal.UserID, *principal.SessionID, identityType)
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

// Refresh handles the business logic for refreshing a user's tokens.
// It parses the refresh token, checks if it's revoked, and if not,
// determines whether to refresh the tokens for a client or a project user.
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

	var refreshToken *contracts.RefreshClaims
	refreshToken, err = uc.tokens.VerifyRefreshToken(ctx, in.RefreshCookie.Value)
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
		attribute.String("session.user_type", sess.UserType),
	)

	if sess.ProjectID == nil {
		return uc.finishClientRefresh(ctx, sess, in, newRefreshJTI, refreshExp)
	}

	return uc.finishProjectUserRefresh(ctx, sess, in, newRefreshJTI, refreshExp)
}

func (uc *CommandService) finishClientRefresh(
	ctx context.Context,
	sess *contracts.Session,
	in RefreshInput,
	refreshJTI uuid.UUID,
	refreshExpiresAt time.Time,
) (tokens *UserTokensOutput, err error) {
	ctx, span := uc.tracer.Start(ctx, "AuthService.finishClientRefresh")
	defer span.End()

	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("finishClientRefresh.success", err == nil))
		}
	}()

	var identity *contracts.Identity
	identity, err = uc.sessions.GetIdentityByIDAndType(ctx, sess.IdentityID, contracts.ClientIdentity)
	if err != nil {
		return nil, err
	}

	var u *contracts.User
	u, err = uc.users.GetUserByID(ctx, identity.EntityID)
	if err != nil {
		return nil, err
	}

	var newAccessJTI uuid.UUID
	newAccessJTI, err = uuid.NewV7()
	if err != nil {
		return nil, fail.New(errx.SYSUUIDV7GenerationError).With(err).WithArgs("auth/finishClientRefresh").RecordCtx(ctx)
	}

	SigningKid, err := uc.keys.GetActiveGoAuthSigningKID(ctx)
	if err != nil {
		return nil, err
	}

	accessExpiresAt := time.Now().Add(15 * time.Minute)
	var accessPayload []byte
	accessPayload, err = uc.tokens.NewAccessToken(contracts.NewAccessTokenInput{
		KID:       SigningKid,
		User:      *u,
		IP:        in.IP,
		Agent:     in.Agent,
		AccessJTI: newAccessJTI.String(),
		SessionID: sess.SessionID,
		FamilyID:  sess.FamilyID,
		ExpiresAt: accessExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	var accessSig []byte
	accessSig, err = uc.keysS.SignGoAuth(ctx, accessPayload)
	if err != nil {
		return nil, err
	}

	accessTokenStr := uc.tokens.AssembleJWT(
		accessPayload,
		accessSig,
	)

	var refreshPayload []byte
	refreshPayload, err = uc.tokens.NewRefreshToken(contracts.NewRefreshTokenInput{
		KID:        SigningKid,
		AccessJTI:  newAccessJTI,
		RefreshJTI: refreshJTI,
		ExpiresAt:  refreshExpiresAt,
		FamilyID:   sess.FamilyID,
	})
	if err != nil {
		return nil, err
	}

	var refreshSig []byte
	refreshSig, err = uc.keysS.SignGoAuth(ctx, refreshPayload)
	if err != nil {
		return nil, err
	}

	refreshTokenStr := uc.tokens.AssembleJWT(
		refreshPayload,
		refreshSig,
	)

	return &UserTokensOutput{
		AccessTokenString:  accessTokenStr,
		RefreshTokenString: refreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
		Domain:             "https://dev.trieauth.trieoh.com",
	}, nil
}

func (uc *CommandService) finishProjectUserRefresh(
	ctx context.Context,
	sess *contracts.Session,
	in RefreshInput,
	refreshJTI uuid.UUID,
	refreshExpiresAt time.Time,
) (*UserTokensOutput, error) {
	ctx, span := uc.tracer.Start(ctx, "AuthService.finishProjectUserRefresh")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("finishProjectUserRefresh.success", err == nil))
		}
	}()

	var identity *contracts.Identity
	identity, err = uc.sessions.GetIdentityByIDAndType(ctx, sess.IdentityID, contracts.ProjectIdentity)
	if err != nil {
		return nil, err
	}

	projectUser, err := uc.projectUsers.GetByIDInternal(ctx, identity.EntityID, *sess.ProjectID)
	if err != nil {
		return nil, err
	}

	var newAccessJTI uuid.UUID
	newAccessJTI, err = uuid.NewV7()
	if err != nil {
		return nil, fail.New(errx.SYSUUIDV7GenerationError).With(err).WithArgs("auth/finishProjectUserRefresh").RecordCtx(ctx)
	}

	SigningKid, err := uc.keys.GetActiveProjectSigningKID(ctx, *sess.ProjectID)
	if err != nil {
		return nil, err
	}

	var accessPayload []byte
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	accessPayload, err = uc.tokens.NewProjectAccessToken(contracts.NewProjectAccessTokenInput{
		KID:       SigningKid,
		User:      *projectUser,
		IP:        in.IP,
		Agent:     in.Agent,
		AccessJTI: newAccessJTI.String(),
		SessionID: sess.SessionID,
		FamilyID:  sess.FamilyID,
		ExpiresAt: accessExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	var accessSig []byte
	accessSig, err = uc.keysS.SignProject(ctx, *sess.ProjectID, accessPayload)
	if err != nil {
		return nil, err
	}

	accessTokenStr := uc.tokens.AssembleJWT(accessPayload, accessSig)

	var refreshPayload []byte
	refreshPayload, err = uc.tokens.NewRefreshToken(contracts.NewRefreshTokenInput{
		KID:        SigningKid,
		AccessJTI:  newAccessJTI,
		RefreshJTI: refreshJTI,
		ExpiresAt:  refreshExpiresAt,
		FamilyID:   sess.FamilyID,
	})
	if err != nil {
		return nil, err
	}

	var refreshSig []byte
	refreshSig, err = uc.keysS.SignProject(ctx, *sess.ProjectID, refreshPayload)
	if err != nil {
		return nil, err
	}

	refreshTokenStr := uc.tokens.AssembleJWT(refreshPayload, refreshSig)

	var project *contracts.Project
	project, err = uc.projects.GetByIDInternal(ctx, *sess.ProjectID)
	if err != nil {
		return nil, err
	}

	return &UserTokensOutput{
		AccessTokenString:  accessTokenStr,
		RefreshTokenString: refreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
		Domain:             project.Domain,
	}, nil
}

type ProjectRegisterInput struct {
	Email        string
	Password     string
	CustomFields *json.RawMessage
	ProjectID    uuid.UUID
	SchemaType   string
	FlowID       string
}

// RegisterProjectUser handles the business logic for creating a new project user.
// It validates the input, hashes the password, and then attempts to create the user in the database.
func (uc *CommandService) RegisterProjectUser(ctx context.Context, in ProjectRegisterInput) error {
	var err error
	ctx, span := uc.tracer.Start(ctx, "AuthService.RegisterProjectUser")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("register.success", err == nil))
		}
	}()

	var verificationEmail ports.Email
	err = uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		verificationEmail, err = uc.registerProjectUserInternal(ctx, in)
		return err
	})
	if err != nil {
		return err
	}

	err = uc.mailSender.Send(ctx, verificationEmail)
	if err != nil {
		return err
	}

	return nil
}

func (uc *CommandService) registerProjectUserInternal(ctx context.Context, in ProjectRegisterInput) (ports.Email, error) {
	ctx, span := uc.tracer.Start(ctx, "AuthService.registerProjectUserInternal",
		trace.WithAttributes(attribute.String("project.id", in.ProjectID.String())),
	)
	defer span.End()

	if in.FlowID == "" {
		return ports.Email{}, fail.New(errx.SchemaEmptyFlowID).RecordCtx(ctx)
	}

	if in.SchemaType == "" {
		return ports.Email{}, fail.New(errx.SchemaEmptySchemaType).RecordCtx(ctx)
	}

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.FlowID = strings.TrimSpace(strings.ToLower(in.FlowID))
	in.SchemaType = strings.TrimSpace(strings.ToLower(in.SchemaType))

	if len(in.Password) > 72 {
		return ports.Email{}, fail.New(errx.AuthInvalidPassword).RecordCtx(ctx)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return ports.Email{}, fail.New(errx.RequestInvalidPassword).With(err).RecordCtx(ctx)
	}

	var usr *contracts.ProjectUser
	usr, err = uc.projectUsers.Register(ctx, contracts.ProjectUser{
		ProjectID:    in.ProjectID,
		Email:        in.Email,
		PasswordHash: string(hashedPassword),
		Metadata:     in.CustomFields,
	})
	if err != nil {
		return ports.Email{}, err
	}

	span.SetAttributes(
		attribute.String("user.id", usr.ID.String()),
		attribute.Int64("user.created_at", usr.CreatedAt.Unix()),
		attribute.String("user.type", usr.UserType),
	)

	var identity *contracts.Identity
	identity, err = uc.sessions.CreateIdentity(ctx, contracts.ProjectIdentity, usr.ID)
	if err != nil {
		return ports.Email{}, err
	}

	span.SetAttributes(
		attribute.String("user.identity.id", identity.ID.String()),
		attribute.String("user.identity.type", string(identity.IdentityType)),
	)

	SigningKid, err := uc.keys.GetActiveGoAuthSigningKID(ctx)
	if err != nil {
		return ports.Email{}, err
	}

	var verificationPayload []byte
	verificationPayload, err = uc.tokens.NewVerificationToken(contracts.NewVerificationTokenInput{
		KID:       SigningKid,
		Subject:   usr.ID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return ports.Email{}, err
	}

	var verificationSig []byte
	verificationSig, err = uc.keysS.SignGoAuth(ctx, verificationPayload)
	if err != nil {
		return ports.Email{}, err
	}

	verificationTokenStr := uc.tokens.AssembleJWT(
		verificationPayload,
		verificationSig,
	)

	var verificationEmail ports.Email
	verificationEmail, err = uc.mailRenderer.Verification(ctx, ports.VerificationEmailData{
		UserID: usr.ID,
		Email:  usr.Email,
		Token:  verificationTokenStr,
		Locale: "en",
	})
	if err != nil {
		return ports.Email{}, err
	}

	return verificationEmail, nil
}

type ProjectLoginInput struct {
	Email     string
	Password  string
	ProjectID uuid.UUID
	IP        string
	Agent     string
}

// LoginProjectUser handles the business logic for logging in a project user.
// It finds the user by email, compares the password, and if successful,
// creates a new session and returns a new set of access and refresh tokens.
func (uc *CommandService) LoginProjectUser(
	ctx context.Context,
	in ProjectLoginInput,
) (
	tokens *UserTokensOutput,
	err error,
) {
	ctx, span := uc.tracer.Start(ctx, "AuthService.LoginProjectUser",
		trace.WithAttributes(attribute.String("project.id", in.ProjectID.String())),
	)
	defer span.End()

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	var usr *contracts.ProjectUser
	usr, err = uc.projectUsers.GetByEmailInternal(ctx, in.ProjectID, in.Email)
	if fail.Is(err, errx.SQLNotFound) {
		return nil, fail.New(errx.AuthInvalidCredentials).RecordCtx(ctx)
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(in.Password))
	if err != nil {
		return nil, fail.New(errx.AuthInvalidCredentials).Trace(err.Error()).RecordCtx(ctx)
	}

	var identity *contracts.Identity
	identity, err = uc.sessions.GetIdentityByEntityIDAndType(ctx, usr.ID, contracts.ProjectIdentity)
	if err != nil {
		return nil, err
	}

	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)
	var sess *contracts.Session
	sess, err = uc.sessions.Create(ctx, contracts.Session{
		IssuedAt:   time.Now(),
		UserAgent:  in.Agent,
		UserIP:     in.IP,
		ExpiresAt:  refreshExpiresAt,
		IdentityID: identity.ID,
		ProjectID:  &in.ProjectID,
	})
	if err != nil {
		return nil, err
	}

	if err = uc.projectUsers.UpdateLastLogin(ctx, identity.EntityID); err != nil {
		return nil, err
	}

	var accessJTI uuid.UUID
	accessJTI, err = uuid.NewV7()
	if err != nil {
		return nil, fail.New(errx.SYSUUIDV7GenerationError).With(err).WithArgs("auth/LoginProjectUser").RecordCtx(ctx)
	}

	var SigningKid string
	SigningKid, err = uc.keys.GetActiveProjectSigningKID(ctx, in.ProjectID)
	if err != nil {
		return nil, err
	}

	accessExpiresAt := time.Now().Add(15 * time.Minute)
	var accessPayload []byte
	accessPayload, err = uc.tokens.NewProjectAccessToken(contracts.NewProjectAccessTokenInput{
		KID:       SigningKid,
		User:      *usr,
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

	var accessSig []byte
	accessSig, err = uc.keysS.SignProject(ctx, in.ProjectID, accessPayload)
	if err != nil {
		return nil, err
	}

	accessTokenStr := uc.tokens.AssembleJWT(accessPayload, accessSig)

	var refreshPayload []byte
	refreshPayload, err = uc.tokens.NewRefreshToken(contracts.NewRefreshTokenInput{
		KID:        SigningKid,
		AccessJTI:  accessJTI,
		RefreshJTI: sess.TokenID,
		ExpiresAt:  refreshExpiresAt,
		FamilyID:   sess.FamilyID,
	})
	if err != nil {
		return nil, err
	}

	var refreshSig []byte
	refreshSig, err = uc.keysS.SignProject(ctx, in.ProjectID, refreshPayload)
	if err != nil {
		return nil, err
	}

	refreshTokenStr := uc.tokens.AssembleJWT(refreshPayload, refreshSig)

	var project *contracts.Project
	project, err = uc.projects.GetByIDInternal(ctx, in.ProjectID)
	if err != nil {
		return nil, err
	}

	return &UserTokensOutput{
		AccessTokenString:  accessTokenStr,
		RefreshTokenString: refreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
		Domain:             project.Domain,
	}, nil
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

	var claims *contracts.VerificationClaims
	claims, err = uc.tokens.VerifyVerificationToken(ctx, token)
	if err != nil {
		return err
	}

	if claims.Sub.Subject != principal.UserID {
		return fail.New(errx.TokenUserMismatch).WithArgs("verification").RecordCtx(ctx)
	}

	var wasAlreadyVerified bool
	if principal.ProjectID == nil {
		span.SetAttributes(attribute.String("user.type", "client"))
		wasAlreadyVerified, err = uc.users.Verify(ctx, claims.Sub.Subject)
		if err != nil {
			return err
		}
	} else {
		span.SetAttributes(attribute.String("user.type", "project"))
		span.SetAttributes(attribute.String("user.project_id", principal.ProjectID.String()))
		wasAlreadyVerified, err = uc.projectUsers.Verify(ctx, claims.Sub.Subject)
		if err != nil {
			return err
		}
	}

	span.SetAttributes(attribute.Bool("user.was_already_verified", wasAlreadyVerified))

	return nil
}

func (uc *CommandService) ResendVerificationEmail(ctx context.Context) (err error) {
	ctx, span := uc.tracer.Start(ctx, "AuthService.ResendVerificationEmail")
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
	var pu *contracts.ProjectUser
	if principal.ProjectID != nil {
		pu, err = uc.projectUsers.GetByIDInternal(ctx, principal.UserID, *principal.ProjectID)
		if err != nil {
			return err
		}
		if pu.IsVerified == true {
			return fail.New(errx.AuthAlreadyVerified).RecordCtx(ctx)
		}
	} else {
		u, err = uc.users.GetUserByID(ctx, principal.UserID)
		if err != nil {
			return err
		}
		if u.IsVerified == true {
			return fail.New(errx.AuthAlreadyVerified).RecordCtx(ctx)
		}
	}

	if pu != nil {
		if pu.IsVerified {
			return fail.New(errx.AuthAlreadyVerified).RecordCtx(ctx)
		}
	} else {
		if u.IsVerified {
			return fail.New(errx.AuthAlreadyVerified).RecordCtx(ctx)
		}
	}

	var SigningKid string
	SigningKid, err = uc.keys.GetActiveGoAuthSigningKID(ctx)
	if err != nil {
		return err
	}

	var verificationPayload []byte
	verificationPayload, err = uc.tokens.NewVerificationToken(contracts.NewVerificationTokenInput{
		KID:       SigningKid,
		Subject:   principal.UserID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return err
	}

	var verificationSig []byte
	verificationSig, err = uc.keysS.SignGoAuth(ctx, verificationPayload)
	if err != nil {
		return err
	}

	verificationTokenStr := uc.tokens.AssembleJWT(verificationPayload, verificationSig)

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

type ForgotPasswordInput struct {
	Email     string
	ProjectID *uuid.UUID
}

func (uc *CommandService) ForgotPassword(ctx context.Context, in ForgotPasswordInput) (err error) {
	ctx, span := uc.tracer.Start(ctx, "AuthService.ForgotPassword")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("forgot_password.success", err == nil))
	}()

	var u *contracts.User
	var pu *contracts.ProjectUser

	u, err = uc.users.GetUserByEmail(ctx, in.Email)
	if err != nil && !errx.IsNotFound(err) {
		return err
	}

	if err != nil && errx.IsNotFound(err) {
		// Global user not found
		if in.ProjectID == nil {
			return nil // silent success
		}

		pu, err = uc.projectUsers.GetByEmailInternal(ctx, *in.ProjectID, in.Email)
		if err != nil {
			if errx.IsNotFound(err) {
				return nil // silent success (no enumeration)
			}
			return err // real failure
		}

	}

	var SigningKid string
	SigningKid, err = uc.keys.GetActiveGoAuthSigningKID(ctx)
	if err != nil {
		return err
	}

	var subjectID uuid.UUID
	var subjectEmail string

	if pu != nil {
		subjectID = pu.ID
		subjectEmail = pu.Email
	} else {
		subjectID = u.ID
		subjectEmail = u.Email
	}

	var resetPayload []byte
	resetPayload, err = uc.tokens.NewResetPasswordToken(contracts.NewResetPasswordInput{
		KID:       SigningKid,
		Subject:   subjectID,
		ProjectID: in.ProjectID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return err
	}

	var resetSig []byte
	resetSig, err = uc.keysS.SignGoAuth(ctx, resetPayload)
	if err != nil {
		return err
	}

	resetPassTokenStr := uc.tokens.AssembleJWT(resetPayload, resetSig)

	var e ports.Email
	e, err = uc.mailRenderer.PasswordReset(ctx, ports.PasswordResetEmailData{
		UserID: subjectID.String(),
		Email:  subjectEmail,
		Token:  resetPassTokenStr,
		Locale: "en",
	})

	err = uc.mailSender.Send(ctx, e)
	if err != nil {
		return err
	}
	return nil
}

type ResetPasswordInput struct {
	NewPassword string
	Token       string
}

func (uc *CommandService) ResetPassword(ctx context.Context, in ResetPasswordInput) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ResetPassword")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("reset_password.success", err == nil))
	}()

	err = uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		return uc.resetPasswordInternal(ctx, in)
	})

	return err
}

func (uc *CommandService) resetPasswordInternal(ctx context.Context, in ResetPasswordInput) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ResetPasswordInternal")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("reset_password.success", err == nil))
	}()

	var claims *contracts.ResetPasswordClaims
	claims, err = uc.tokens.VerifyResetPasswordToken(ctx, in.Token)
	if err != nil {
		return err
	}

	var jti uuid.UUID
	jti, err = uuid.Parse(claims.ID)
	if err != nil {
		return fail.New(errx.RequestParseUUIDError).RecordCtx(ctx)
	}

	var exists bool
	exists, err = uc.tokenReuseList.Exists(ctx, jti, claims.Sub.Subject)
	if err != nil {
		return err
	}
	if exists {
		// FIXME when the audit is implemented add this to the audit
		return fail.New(errx.AuthTokenAlreadyUsed).RecordCtx(ctx)
	}

	if len(in.NewPassword) > 72 {
		return fail.New(errx.AuthInvalidPassword).RecordCtx(ctx)
	}

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fail.New(errx.RequestInvalidPassword).With(err).RecordCtx(ctx)
	}

	if claims.Sub.ProjectID == nil {
		err = uc.users.ResetPassword(ctx, claims.Sub.Subject, hashedPassword)
		if err != nil {
			return err
		}
		_, err = uc.sessions.MarkRevokedByFilter(ctx, contracts.Filter{
			EntityID:     claims.Sub.Subject,
			IdentityType: contracts.ClientIdentity,
		})
		if err != nil {
			return err
		}
	} else {
		err = uc.projectUsers.ResetPassword(ctx, claims.Sub.Subject, hashedPassword)
		if err != nil {
			return err
		}
		_, err = uc.sessions.MarkRevokedByFilter(ctx, contracts.Filter{
			EntityID:     claims.Sub.Subject,
			IdentityType: contracts.ProjectIdentity,
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

type ProjectLogoutInput struct {
	ProjectID          uuid.UUID
	RefreshTokenCookie *http.Cookie
}

func (uc *CommandService) LogoutProjectUser(ctx context.Context, in ProjectLogoutInput) error {
	ctx, span := uc.tracer.Start(ctx, "AuthService.LogoutProjectUser")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("logout.success", err == nil))
		}
	}()

	var refreshToken *contracts.RefreshClaims
	refreshToken, err = uc.tokens.VerifyRefreshToken(ctx, in.RefreshTokenCookie.Value)
	if err != nil {
		return err
	}

	var sess *contracts.Session
	sess, err = uc.sessions.GetByFamilyID(ctx, refreshToken.Sub.FamilyID)
	if err != nil {
		return err
	}

	identity, err := uc.sessions.GetIdentityByIDAndType(ctx, sess.IdentityID, contracts.ProjectIdentity)
	if err != nil {
		return err
	}

	_, err = uc.sessions.MarkRevokedByID(ctx, identity.EntityID, sess.SessionID, contracts.ProjectIdentity)
	if err != nil {
		return err
	}

	return nil
}

func (uc *CommandService) GetJWKS(ctx context.Context) (map[string]any, error) {
	ekeys, err := uc.keys.ListGoAuthPublicKeys(ctx)
	if err != nil {
		telemetry.Log().Error("Failed listing GoAuth public keys", zap.Error(err))
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
		"keys": jwkKeys,
	}, nil
}
