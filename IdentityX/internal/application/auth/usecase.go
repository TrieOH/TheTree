package auth

import (
	"GoAuth/internal/adapters/email"
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/tokens"
	"GoAuth/internal/application/validation"
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/key"
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/domain/user"
	"GoAuth/internal/infrastructure"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	usecaseTracer = otel.Tracer("auth_usecase")
)

type UseCase struct {
	deps          Deps
	keys          inbounds.KeysService
	schema        inbounds.SchemaService
	tokenIssuer   inbounds.TokenIssuer
	tokenVerifier inbounds.TokenVerifier
	mailRenderer  outbounds.EmailRenderer
	mailSender    outbounds.Mailer
	tx            inbounds.TxRunner
}

type Deps struct {
	Users          outbounds.UserRepository
	Sessions       outbounds.SessionRepository
	Schemas        outbounds.SchemaRepository
	Versions       outbounds.SchemaVersionRepository
	Fields         outbounds.SchemaFieldsRepository
	Projects       outbounds.ProjectRepository
	ProjectUsers   outbounds.ProjectUserRepository
	Keys           outbounds.KeysRepository
	TokenReuseList outbounds.TokenReuseListRepository
	Cache          outbounds.CacheService
}

var _ inbounds.AuthService = (*UseCase)(nil)

func New(
	repos Deps,
	infra infrastructure.Infra,
	keys inbounds.KeysService,
	schema inbounds.SchemaService,
	tokenBundle tokens.TokenBundle,
	mailBundle email.MailBundle,
) inbounds.AuthService {
	return &UseCase{
		deps:          repos,
		keys:          keys,
		schema:        schema,
		tokenIssuer:   tokenBundle.Issuer,
		tokenVerifier: tokenBundle.Verifier,
		mailRenderer:  mailBundle.Renderer,
		mailSender:    mailBundle.Mailer,
		tx:            infra.Tx,
	}
}

// Register handles the business logic for creating a new user.
// It validates the input, hashes the password, and then attempts to create the user in the database.
// It returns an error if the email is already in use or if there is a problem with the database.
func (uc *UseCase) Register(ctx context.Context, in inbounds.RegisterUserInput) error {
	var err error
	ctx, span := usecaseTracer.Start(ctx, "AuthService.Register")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("register.success", err == nil))
		}
	}()

	var verificationEmail outbounds.Email
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

func (uc *UseCase) registerInternal(ctx context.Context, in inbounds.RegisterUserInput) (outbounds.Email, error) {
	var err error
	ctx, span := usecaseTracer.Start(ctx, "AuthService.registerInternal")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("register.success", err == nil))
		}
	}()

	users := uc.deps.Users
	sessions := uc.deps.Sessions
	keys := uc.keys
	issuer := uc.tokenIssuer

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	if len(in.Password) > 72 {
		return outbounds.Email{}, fail.New(apierr.AuthInvalidPassword).RecordCtx(ctx)
	}

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return outbounds.Email{}, fail.New(apierr.RequestInvalidPassword).With(err).RecordCtx(ctx)
	}

	var u *user.User
	u, err = users.Register(ctx, in.Email, string(hashedPassword))
	if err != nil {
		return outbounds.Email{}, err
	}

	span.SetAttributes(
		attribute.String("user.id", u.ID.String()),
		attribute.Int64("user.created_at", u.CreatedAt.Unix()),
		attribute.String("user.type", u.UserType),
	)

	var identity *session.Identity
	identity, err = sessions.CreateIdentity(ctx, session.ClientIdentity, u.ID)
	if err != nil {
		return outbounds.Email{}, err
	}

	span.SetAttributes(
		attribute.String("user.identity.id", identity.ID.String()),
		attribute.String("user.identity.type", string(identity.IdentityType)),
	)

	var SigningKid string
	SigningKid, err = keys.GetActiveGoAuthSigningKID(ctx)
	if err != nil {
		return outbounds.Email{}, err
	}

	var verificationPayload []byte
	verificationPayload, err = issuer.NewVerificationToken(inbounds.NewVerificationTokenInput{
		KID:       SigningKid,
		Subject:   u.ID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return outbounds.Email{}, err
	}

	var verificationSig []byte
	verificationSig, err = keys.SignGoAuth(ctx, verificationPayload)
	if err != nil {
		return outbounds.Email{}, err
	}

	verificationTokenStr := issuer.AssembleJWT(verificationPayload, verificationSig)

	var verificationEmail outbounds.Email
	verificationEmail, err = uc.mailRenderer.Verification(ctx, outbounds.VerificationEmailData{
		UserID: u.ID,
		Email:  u.Email,
		Token:  verificationTokenStr,
		Locale: "en",
	})
	if err != nil {
		return outbounds.Email{}, err
	}

	return verificationEmail, nil
}

// Login handles the business logic for logging in a user.
// It finds the user by email, compares the password, and if successful,
// creates a new session and returns a new set of access and refresh tokens.
func (uc *UseCase) Login(ctx context.Context, in inbounds.LoginUserInput) (tokens *inbounds.UserTokensOutput, err error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	ctx, span := usecaseTracer.Start(ctx, "AuthService.Login")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("login.success", err == nil))
		}
	}()

	users := uc.deps.Users
	sessions := uc.deps.Sessions
	keys := uc.keys
	issuer := uc.tokenIssuer

	var u *user.User
	u, err = users.GetUserByEmail(ctx, in.Email)
	if fail.Is(err, apierr.SQLNotFound) {
		return nil, fail.New(apierr.AuthInvalidCredentials).RecordCtx(ctx)
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
		return nil, fail.New(apierr.AuthInvalidCredentials).Trace(err.Error()).RecordCtx(ctx)
	}

	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)

	var identity *session.Identity
	identity, err = sessions.GetIdentityByEntityIDAndType(ctx, u.ID, session.ClientIdentity)
	if err != nil {
		return nil, err
	}

	var sess *session.Session
	sess, err = sessions.Create(ctx, session.Session{
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
		return nil, fail.New(apierr.SYSUUIDV7GenerationError).With(err).WithArgs("auth/login").RecordCtx(ctx)
	}

	var SigningKid string
	SigningKid, err = keys.GetActiveGoAuthSigningKID(ctx)
	if err != nil {
		return nil, err
	}

	accessExpiresAt := time.Now().Add(15 * time.Minute)
	var accessPayload []byte
	accessPayload, err = issuer.NewAccessToken(inbounds.NewAccessTokenInput{
		KID:       SigningKid,
		User:      *u,
		IP:        in.IP,
		Agent:     in.Agent,
		AccessJTI: accessJTI.String(),
		SessionID: sess.SessionID,
		ExpiresAt: accessExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	var accessSig []byte
	accessSig, err = keys.SignGoAuth(ctx, accessPayload)
	if err != nil {
		return nil, err
	}

	accessTokenStr := issuer.AssembleJWT(accessPayload, accessSig)

	var refreshPayload []byte
	refreshPayload, err = issuer.NewRefreshToken(inbounds.NewRefreshTokenInput{
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
	refreshSig, err = keys.SignGoAuth(ctx, refreshPayload)
	if err != nil {
		return nil, err
	}

	refreshTokenStr := issuer.AssembleJWT(refreshPayload, refreshSig)

	return &inbounds.UserTokensOutput{
		AccessTokenString:  accessTokenStr,
		RefreshTokenString: refreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
		IsUpToDate:         true,
	}, nil
}

// Logout handles the business logic for logging out a user.
// It retrieves the principal from the context, deletes the session, and revokes the refresh token.
func (uc *UseCase) Logout(ctx context.Context) error {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.Logout")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("logout.success", err == nil))
		}
	}()

	sessions := uc.deps.Sessions

	var principal *authz.Principal
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var identityType session.IdentityType
	if principal.ProjectID == nil {
		identityType = session.ClientIdentity
	} else {
		identityType = session.ProjectIdentity
	}

	var sess *session.Session
	sess, err = sessions.MarkRevokedByID(ctx, principal.UserID, principal.SessionID, identityType)
	if err != nil {
		return err
	}

	span.SetAttributes(attribute.String("session.id", sess.SessionID.String()))

	return nil
}

// Refresh handles the business logic for refreshing a user's tokens.
// It parses the refresh token, checks if it's revoked, and if not,
// determines whether to refresh the tokens for a client or a project user.
func (uc *UseCase) Refresh(ctx context.Context, in inbounds.RefreshInput) (*inbounds.UserTokensOutput, error) {
	txOptions := inbounds.TxOptions{
		Isolation: pgx.ReadCommitted,
		ReadOnly:  pgx.ReadWrite,
	}

	var out *inbounds.UserTokensOutput
	err := uc.tx.WithinTxWithOptions(ctx, txOptions, func(ctx context.Context) error {
		var err error
		out, err = uc.refreshInternal(ctx, in)
		return err
	})

	return out, err
}

func (uc *UseCase) refreshInternal(ctx context.Context, in inbounds.RefreshInput) (*inbounds.UserTokensOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.Refresh")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("refresh.success", err == nil))
		}
	}()

	sessions := uc.deps.Sessions

	var refreshToken *auth.RefreshClaims
	refreshToken, err = uc.tokenVerifier.VerifyRefreshToken(ctx, in.RefreshCookie.Value)
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
		return nil, fail.New(apierr.SYSUUIDV7GenerationError).With(err).WithArgs("auth/refreshInternal").RecordCtx(ctx)
	}

	var newRefreshJTI = uid
	var refreshExp = time.Now().Add(7 * 24 * time.Hour)

	span.SetAttributes(attribute.String("old_token.id", oldJTI.String()))
	span.SetAttributes(attribute.String("new_token.id", newRefreshJTI.String()))

	var sess *session.Session
	sess, err = sessions.GetByFamilyID(ctx, refreshToken.Sub.FamilyID)
	if err != nil {
		return nil, fail.New(apierr.SessionNotFound).RecordCtx(ctx)
	}

	now := time.Now()
	if sess.ExpiresAt.Before(now) || sess.RevokedAt != nil {
		// FIXME Record suspicious behaviour on audit when it is implemented
		return nil, fail.New(apierr.SessionNotFound).RecordCtx(ctx)
	}

	// should revoke the session because of replay attacks
	// FIXME Add suspicious behaviour to audit when it is implemented
	if sess.TokenID != oldJTI {
		_ = sessions.MarkRevokedByFamilyID(ctx, sess.FamilyID)
		return nil, fail.New(apierr.TokenReuseIdentified).WithArgs("refresh").RecordCtx(ctx)
	}

	sess, err = sessions.RotateToken(ctx, refreshToken.Sub.FamilyID, newRefreshJTI, oldJTI, refreshExp)
	if fail.Is(err, apierr.SQLNotFound) {
		// sql.ErrNoRows → raced / reused / revoked
		return nil, fail.New(apierr.SessionNotFound).RecordCtx(ctx)
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

func (uc *UseCase) finishClientRefresh(
	ctx context.Context,
	sess *session.Session,
	in inbounds.RefreshInput,
	refreshJTI uuid.UUID,
	refreshExpiresAt time.Time,
) (tokens *inbounds.UserTokensOutput, err error) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.finishClientRefresh")
	defer span.End()

	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("finishClientRefresh.success", err == nil))
		}
	}()

	users := uc.deps.Users
	sessions := uc.deps.Sessions
	keys := uc.keys
	issuer := uc.tokenIssuer

	var identity *session.Identity
	identity, err = sessions.GetIdentityByIDAndType(ctx, sess.IdentityID, session.ClientIdentity)
	if err != nil {
		return nil, err
	}

	var u *user.User
	u, err = users.GetUserByID(ctx, identity.EntityID)
	if err != nil {
		return nil, err
	}

	var newAccessJTI uuid.UUID
	newAccessJTI, err = uuid.NewV7()
	if err != nil {
		return nil, fail.New(apierr.SYSUUIDV7GenerationError).With(err).WithArgs("auth/finishClientRefresh").RecordCtx(ctx)
	}

	SigningKid, err := keys.GetActiveGoAuthSigningKID(ctx)
	if err != nil {
		return nil, err
	}

	accessExpiresAt := time.Now().Add(15 * time.Minute)
	var accessPayload []byte
	accessPayload, err = issuer.NewAccessToken(inbounds.NewAccessTokenInput{
		KID:       SigningKid,
		User:      *u,
		IP:        in.IP,
		Agent:     in.Agent,
		AccessJTI: newAccessJTI.String(),
		SessionID: sess.SessionID,
		ExpiresAt: accessExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	var accessSig []byte
	accessSig, err = keys.SignGoAuth(ctx, accessPayload)
	if err != nil {
		return nil, err
	}

	accessTokenStr := issuer.AssembleJWT(
		accessPayload,
		accessSig,
	)

	var refreshPayload []byte
	refreshPayload, err = issuer.NewRefreshToken(inbounds.NewRefreshTokenInput{
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
	refreshSig, err = keys.SignGoAuth(ctx, refreshPayload)
	if err != nil {
		return nil, err
	}

	refreshTokenStr := issuer.AssembleJWT(
		refreshPayload,
		refreshSig,
	)

	return &inbounds.UserTokensOutput{
		AccessTokenString:  accessTokenStr,
		RefreshTokenString: refreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
		IsUpToDate:         true,
	}, nil
}

func (uc *UseCase) finishProjectUserRefresh(
	ctx context.Context,
	sess *session.Session,
	in inbounds.RefreshInput,
	refreshJTI uuid.UUID,
	refreshExpiresAt time.Time,
) (*inbounds.UserTokensOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.finishProjectUserRefresh")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("finishProjectUserRefresh.success", err == nil))
		}
	}()

	projectUsers := uc.deps.ProjectUsers
	sessions := uc.deps.Sessions
	keys := uc.keys
	issuer := uc.tokenIssuer

	var identity *session.Identity
	identity, err = sessions.GetIdentityByIDAndType(ctx, sess.IdentityID, session.ProjectIdentity)
	if err != nil {
		return nil, err
	}

	projectUser, err := projectUsers.GetByIDInternal(ctx, identity.EntityID, *sess.ProjectID)
	if err != nil {
		return nil, err
	}

	var newAccessJTI uuid.UUID
	newAccessJTI, err = uuid.NewV7()
	if err != nil {
		return nil, fail.New(apierr.SYSUUIDV7GenerationError).With(err).WithArgs("auth/finishProjectUserRefresh").RecordCtx(ctx)
	}

	SigningKid, err := keys.GetActiveProjectSigningKID(ctx, *sess.ProjectID)
	if err != nil {
		return nil, err
	}

	var accessPayload []byte
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	accessPayload, err = issuer.NewProjectAccessToken(inbounds.NewProjectAccessTokenInput{
		KID:       SigningKid,
		User:      *projectUser,
		IP:        in.IP,
		Agent:     in.Agent,
		AccessJTI: newAccessJTI.String(),
		SessionID: sess.SessionID,
		ExpiresAt: accessExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	var accessSig []byte
	accessSig, err = keys.SignProject(ctx, *sess.ProjectID, accessPayload)
	if err != nil {
		return nil, err
	}

	accessTokenStr := issuer.AssembleJWT(accessPayload, accessSig)

	var refreshPayload []byte
	refreshPayload, err = issuer.NewRefreshToken(inbounds.NewRefreshTokenInput{
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
	refreshSig, err = keys.SignProject(ctx, *sess.ProjectID, refreshPayload)
	if err != nil {
		return nil, err
	}

	refreshTokenStr := issuer.AssembleJWT(refreshPayload, refreshSig)

	isUpToDate, err := uc.schema.CheckSchemaCompatibility(ctx, projectUser.ID, *sess.ProjectID)
	if err != nil {
		logs.L().Error("Failed to check schema compatibility during refresh", zap.Error(err))
		isUpToDate = false
	}

	return &inbounds.UserTokensOutput{
		AccessTokenString:  accessTokenStr,
		RefreshTokenString: refreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
		IsUpToDate:         isUpToDate,
	}, nil
}

// RegisterProjectUser handles the business logic for creating a new project user.
// It validates the input, hashes the password, and then attempts to create the user in the database.
func (uc *UseCase) RegisterProjectUser(ctx context.Context, in inbounds.ProjectRegisterInput) error {
	var err error
	ctx, span := usecaseTracer.Start(ctx, "AuthService.RegisterProjectUser")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("register.success", err == nil))
		}
	}()

	var verificationEmail outbounds.Email
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

func (uc *UseCase) registerProjectUserInternal(ctx context.Context, in inbounds.ProjectRegisterInput) (outbounds.Email, error) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.registerProjectUserInternal",
		trace.WithAttributes(attribute.String("project.id", in.ProjectID.String())),
	)
	defer span.End()

	projectUsers := uc.deps.ProjectUsers
	sessions := uc.deps.Sessions
	keys := uc.keys
	issuer := uc.tokenIssuer

	if in.FlowID == "" {
		return outbounds.Email{}, fail.New(apierr.SchemaEmptyFlowID).RecordCtx(ctx)
	}

	if in.SchemaType == "" {
		return outbounds.Email{}, fail.New(apierr.SchemaEmptySchemaType).RecordCtx(ctx)
	}

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.FlowID = strings.TrimSpace(strings.ToLower(in.FlowID))
	in.SchemaType = strings.TrimSpace(strings.ToLower(in.SchemaType))

	if !schema.IsValidSchemaType(in.SchemaType) {
		return outbounds.Email{}, fail.New(apierr.SchemaInvalidSchemaType).RecordCtx(ctx)
	}

	// FlowIDs cannot be the same as schema types so if this matches we error out
	if schema.IsValidSchemaType(in.FlowID) {
		return outbounds.Email{}, fail.New(apierr.SchemaInvalidFlowID).WithArgs("flow id can't be the same as a schema type").RecordCtx(ctx)
	}

	if schema.Type(in.SchemaType) == schema.Core && schema.IsFlowIDReserved(in.FlowID) && in.CustomFields != nil {
		return outbounds.Email{}, fail.New(apierr.SchemaMetadataNotAllowed).RecordCtx(ctx)
	}

	empty := json.RawMessage(`{}`)
	customFields := &empty

	// Validate and construct metadata for non-core or non-reserved flows
	isCoreWithReservedFlow := schema.Type(in.SchemaType) == schema.Core && schema.IsFlowIDReserved(in.FlowID)
	if !isCoreWithReservedFlow {
		validatedMetadata, err := uc.schema.ValidateAndConstructMetadata(ctx, in.ProjectID, schema.Type(in.SchemaType), in.FlowID, in.CustomFields)
		if err != nil {
			return outbounds.Email{}, err
		}
		if validatedMetadata != nil {
			customFields = validatedMetadata
		}
	}

	if len(in.Password) > 72 {
		return outbounds.Email{}, fail.New(apierr.AuthInvalidPassword).RecordCtx(ctx)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return outbounds.Email{}, fail.New(apierr.RequestInvalidPassword).With(err).RecordCtx(ctx)
	}

	var usr *project_users.ProjectUser
	usr, err = projectUsers.Register(ctx, project_users.ProjectUser{
		ProjectID:    in.ProjectID,
		Email:        in.Email,
		PasswordHash: string(hashedPassword),
		Metadata:     customFields,
	})
	if err != nil {
		return outbounds.Email{}, err
	}

	span.SetAttributes(
		attribute.String("user.id", usr.ID.String()),
		attribute.Int64("user.created_at", usr.CreatedAt.Unix()),
		attribute.String("user.type", usr.UserType),
	)

	var identity *session.Identity
	identity, err = sessions.CreateIdentity(ctx, session.ProjectIdentity, usr.ID)
	if err != nil {
		return outbounds.Email{}, err
	}

	span.SetAttributes(
		attribute.String("user.identity.id", identity.ID.String()),
		attribute.String("user.identity.type", string(identity.IdentityType)),
	)

	SigningKid, err := keys.GetActiveGoAuthSigningKID(ctx)
	if err != nil {
		return outbounds.Email{}, err
	}

	var verificationPayload []byte
	verificationPayload, err = issuer.NewVerificationToken(inbounds.NewVerificationTokenInput{
		KID:       SigningKid,
		Subject:   usr.ID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return outbounds.Email{}, err
	}

	var verificationSig []byte
	verificationSig, err = keys.SignGoAuth(ctx, verificationPayload)
	if err != nil {
		return outbounds.Email{}, err
	}

	verificationTokenStr := issuer.AssembleJWT(
		verificationPayload,
		verificationSig,
	)

	var verificationEmail outbounds.Email
	verificationEmail, err = uc.mailRenderer.Verification(ctx, outbounds.VerificationEmailData{
		UserID: usr.ID,
		Email:  usr.Email,
		Token:  verificationTokenStr,
		Locale: "en",
	})
	if err != nil {
		return outbounds.Email{}, err
	}

	return verificationEmail, nil
}

// LoginProjectUser handles the business logic for logging in a project user.
// It finds the user by email, compares the password, and if successful,
// creates a new session and returns a new set of access and refresh tokens.
func (uc *UseCase) LoginProjectUser(
	ctx context.Context,
	in inbounds.ProjectLoginInput,
) (
	tokens *inbounds.UserTokensOutput,
	err error,
) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.LoginProjectUser",
		trace.WithAttributes(attribute.String("project.id", in.ProjectID.String())),
	)
	defer span.End()

	projectUsers := uc.deps.ProjectUsers
	sessions := uc.deps.Sessions
	keys := uc.keys
	issuer := uc.tokenIssuer

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	var usr *project_users.ProjectUser
	usr, err = projectUsers.GetByEmailInternal(ctx, in.ProjectID, in.Email)
	if fail.Is(err, apierr.SQLNotFound) {
		return nil, fail.New(apierr.AuthInvalidCredentials).RecordCtx(ctx)
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(in.Password))
	if err != nil {
		return nil, fail.New(apierr.AuthInvalidCredentials).Trace(err.Error()).RecordCtx(ctx)
	}

	var identity *session.Identity
	identity, err = sessions.GetIdentityByEntityIDAndType(ctx, usr.ID, session.ProjectIdentity)
	if err != nil {
		return nil, err
	}

	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)
	var sess *session.Session
	sess, err = sessions.Create(ctx, session.Session{
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

	var accessJTI uuid.UUID
	accessJTI, err = uuid.NewV7()
	if err != nil {
		return nil, fail.New(apierr.SYSUUIDV7GenerationError).With(err).WithArgs("auth/LoginProjectUser").RecordCtx(ctx)
	}

	var SigningKid string
	SigningKid, err = keys.GetActiveProjectSigningKID(ctx, in.ProjectID)
	if err != nil {
		return nil, err
	}

	accessExpiresAt := time.Now().Add(15 * time.Minute)
	var accessPayload []byte
	accessPayload, err = issuer.NewProjectAccessToken(inbounds.NewProjectAccessTokenInput{
		KID:       SigningKid,
		User:      *usr,
		IP:        in.IP,
		Agent:     in.Agent,
		AccessJTI: accessJTI.String(),
		SessionID: sess.SessionID,
		ExpiresAt: accessExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	var accessSig []byte
	accessSig, err = keys.SignProject(ctx, in.ProjectID, accessPayload)
	if err != nil {
		return nil, err
	}

	accessTokenStr := issuer.AssembleJWT(accessPayload, accessSig)

	var refreshPayload []byte
	refreshPayload, err = uc.tokenIssuer.NewRefreshToken(inbounds.NewRefreshTokenInput{
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
	refreshSig, err = keys.SignProject(ctx, in.ProjectID, refreshPayload)
	if err != nil {
		return nil, err
	}

	refreshTokenStr := issuer.AssembleJWT(refreshPayload, refreshSig)

	isUpToDate, err := uc.schema.CheckSchemaCompatibility(ctx, usr.ID, in.ProjectID)
	if err != nil {
		logs.L().Error("Failed to check schema compatibility during login", zap.Error(err))
		isUpToDate = false
	}

	return &inbounds.UserTokensOutput{
		AccessTokenString:  accessTokenStr,
		RefreshTokenString: refreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
		IsUpToDate:         isUpToDate,
	}, nil
}

func (uc *UseCase) Verify(ctx context.Context, token string) (err error) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.Verify")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("verify.success", err == nil))
		}
	}()

	users := uc.deps.Users
	projectUsers := uc.deps.ProjectUsers

	var principal *authz.Principal
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	var claims *auth.VerificationClaims
	claims, err = uc.tokenVerifier.VerifyVerificationToken(ctx, token)
	if err != nil {
		return err
	}

	if claims.Sub.Subject != principal.UserID {
		return fail.New(apierr.TokenUserMismatch).WithArgs("verification").RecordCtx(ctx)
	}

	var wasAlreadyVerified bool
	if principal.ProjectID == nil {
		span.SetAttributes(attribute.String("user.type", "client"))
		wasAlreadyVerified, err = users.Verify(ctx, claims.Sub.Subject)
		if err != nil {
			return err
		}
	} else {
		span.SetAttributes(attribute.String("user.type", "project"))
		span.SetAttributes(attribute.String("user.project_id", principal.ProjectID.String()))
		wasAlreadyVerified, err = projectUsers.Verify(ctx, claims.Sub.Subject)
		if err != nil {
			return err
		}
	}

	span.SetAttributes(attribute.Bool("user.was_already_verified", wasAlreadyVerified))

	return nil
}

func (uc *UseCase) ResendVerificationEmail(ctx context.Context) (err error) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.ResendVerificationEmail")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("resend_verification.success", err == nil))
		}
	}()

	users := uc.deps.Users
	projectUsers := uc.deps.ProjectUsers
	keys := uc.keys
	issuer := uc.tokenIssuer

	var principal *authz.Principal
	principal, err = authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	if principal.IsVerified == true {
		return fail.New(apierr.AuthAlreadyVerified).RecordCtx(ctx)
	}

	if principal.ProjectID != nil {
		u, err := projectUsers.GetByIDInternal(ctx, principal.UserID, *principal.ProjectID)
		if err != nil {
			return err
		}
		if u.IsVerified == true {
			return fail.New(apierr.AuthAlreadyVerified).RecordCtx(ctx)
		}
	} else {
		u, err := users.GetUserByID(ctx, principal.UserID)
		if err != nil {
			return err
		}
		if u.IsVerified == true {
			return fail.New(apierr.AuthAlreadyVerified).RecordCtx(ctx)
		}
	}

	SigningKid, err := keys.GetActiveGoAuthSigningKID(ctx)
	if err != nil {
		return err
	}

	var verificationPayload []byte
	verificationPayload, err = issuer.NewVerificationToken(inbounds.NewVerificationTokenInput{
		KID:       SigningKid,
		Subject:   principal.UserID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return err
	}

	var verificationSig []byte
	verificationSig, err = keys.SignGoAuth(ctx, verificationPayload)
	if err != nil {
		return err
	}

	verificationTokenStr := issuer.AssembleJWT(verificationPayload, verificationSig)

	var verificationEmail outbounds.Email
	verificationEmail, err = uc.mailRenderer.Verification(ctx, outbounds.VerificationEmailData{
		UserID: principal.UserID,
		Email:  principal.Email,
		Token:  verificationTokenStr,
		Locale: "en",
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

func (uc *UseCase) ForgotPassword(ctx context.Context, in inbounds.ForgotPasswordInput) (err error) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.ForgotPassword")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("forgot_password.success", err == nil))
	}()

	var u *user.User
	var pu *project_users.ProjectUser

	u, err = uc.deps.Users.GetUserByEmail(ctx, in.Email)
	if err != nil && !apierr.IsNotFoundNew(err) {
		return err
	}

	if err != nil && apierr.IsNotFoundNew(err) {
		// Global user not found
		if in.ProjectID == nil {
			return nil // silent success
		}

		pu, err = uc.deps.ProjectUsers.GetByEmailInternal(ctx, *in.ProjectID, in.Email)
		if err != nil {
			if apierr.IsNotFoundNew(err) {
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
	resetPayload, err = uc.tokenIssuer.NewResetPasswordToken(inbounds.NewResetPasswordInput{
		KID:       SigningKid,
		Subject:   subjectID,
		ProjectID: in.ProjectID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return err
	}

	var resetSig []byte
	resetSig, err = uc.keys.SignGoAuth(ctx, resetPayload)
	if err != nil {
		return err
	}

	resetPassTokenStr := uc.tokenIssuer.AssembleJWT(resetPayload, resetSig)

	var e outbounds.Email
	e, err = uc.mailRenderer.PasswordReset(ctx, outbounds.PasswordResetEmailData{
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

func (uc *UseCase) ResetPassword(ctx context.Context, in inbounds.ResetPasswordInput) (err error) {
	ctx, span := usecaseTracer.Start(ctx, "ResetPassword")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("reset_password.success", err == nil))
	}()

	err = uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		return uc.resetPasswordInternal(ctx, in)
	})

	return err
}

func (uc *UseCase) resetPasswordInternal(ctx context.Context, in inbounds.ResetPasswordInput) (err error) {
	ctx, span := usecaseTracer.Start(ctx, "ResetPasswordInternal")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("reset_password.success", err == nil))
	}()

	var claims *auth.ResetPasswordClaims
	claims, err = uc.tokenVerifier.VerifyResetPasswordToken(ctx, in.Token)
	if err != nil {
		return err
	}

	var jti uuid.UUID
	jti, err = uuid.Parse(claims.ID)
	if err != nil {
		return fail.New(apierr.RequestParseUUIDError).RecordCtx(ctx)
	}

	var exists bool
	exists, err = uc.deps.TokenReuseList.Exists(ctx, jti, claims.Sub.Subject)
	if err != nil {
		return err
	}
	if exists {
		// FIXME when the audit is implemented add this to the audit
		return fail.New(apierr.AuthTokenAlreadyUsed).RecordCtx(ctx)
	}

	if len(in.NewPassword) > 72 {
		return fail.New(apierr.AuthInvalidPassword).RecordCtx(ctx)
	}

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fail.New(apierr.RequestInvalidPassword).With(err).RecordCtx(ctx)
	}

	if claims.Sub.ProjectID == nil {
		err = uc.deps.Users.ResetPassword(ctx, claims.Sub.Subject, hashedPassword)
		if err != nil {
			return err
		}
		_, err = uc.deps.Sessions.MarkRevokedByFilter(ctx, session.Filter{
			EntityID:     claims.Sub.Subject,
			IdentityType: session.ClientIdentity,
		})
		if err != nil {
			return err
		}
	} else {
		err = uc.deps.ProjectUsers.ResetPassword(ctx, claims.Sub.Subject, hashedPassword)
		if err != nil {
			return err
		}
		_, err = uc.deps.Sessions.MarkRevokedByFilter(ctx, session.Filter{
			EntityID:     claims.Sub.Subject,
			IdentityType: session.ProjectIdentity,
		})
		if err != nil {
			return err
		}
	}

	err = uc.deps.TokenReuseList.Append(ctx, jti, claims.Sub.Subject, claims.ExpiresAt.Time)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GetJWKS(ctx context.Context) (map[string]any, error) {
	keys, err := uc.deps.Keys.ListGoAuthPublicKeys(ctx)
	if err != nil {
		logs.L().Error("Failed listing GoAuth public keys", zap.Error(err))
		return nil, fail.New(apierr.SYSJWKSRetrievalFailed).With(err).RecordCtx(ctx)
	}

	jwkKeys := make([]any, len(keys))
	for i, k := range keys {
		jwkKeys[i] = key.PublicKeyToJWK(k)
	}

	return map[string]any{
		"keys": jwkKeys,
	}, nil
}
