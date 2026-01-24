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
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
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
	tokenIssuer   inbounds.TokenIssuer
	tokenVerifier inbounds.TokenVerifier
	mailRenderer  outbounds.EmailRenderer
	mailSender    outbounds.Mailer
	tx            inbounds.TxRunner
}

type Deps struct {
	Users        outbounds.UserRepository
	Sessions     outbounds.SessionRepository
	Schemas      outbounds.SchemaRepository
	Versions     outbounds.SchemaVersionRepository
	Fields       outbounds.SchemaFieldsRepository
	Projects     outbounds.ProjectRepository
	ProjectUsers outbounds.ProjectUserRepository
	Keys         outbounds.KeysRepository
}

var _ inbounds.AuthService = (*UseCase)(nil)

func New(
	repos Deps,
	infra infrastructure.Infra,
	keys inbounds.KeysService,
	tokenBundle tokens.TokenBundle,
	mailBundle email.MailBundle,
) inbounds.AuthService {
	return &UseCase{
		deps:          repos,
		keys:          keys,
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
		return apierr.FromService(span, err)
	}

	err = uc.mailSender.Send(ctx, verificationEmail)
	if err != nil {
		return apierr.FromService(span, err)
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
		return outbounds.Email{}, apierr.ErrPasswordTooLong{}
	}

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return outbounds.Email{}, inbounds.ErrHashingPassword{Cause: err}
	}

	var u *user.User
	u, err = users.Register(ctx, in.Email, string(hashedPassword))
	if apierr.IsConflict(err) {
		return outbounds.Email{}, inbounds.ErrEmailAlreadyInUse{Cause: errors.New("email already in use")}
	} else if err != nil {
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
	verificationEmail, err = uc.mailRenderer.Verification(outbounds.VerificationEmailData{
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
	if apierr.IsNotFound(err) {
		return nil, apierr.FromService(span, inbounds.ErrInvalidCredentials{Cause: nil})
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
		return nil, apierr.FromService(span, inbounds.ErrInvalidCredentials{Cause: err})
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
		return nil, apierr.FromService(span, inbounds.ErrGeneratingUUID{Cause: err, Location: "auth/login"})
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
		return nil, apierr.FromService(span, err)
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
		return nil, apierr.FromService(span, err)
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
	principal, err = RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return apierr.FromService(span, err)
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
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
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
		return nil, apierr.FromService(span, err)
	}

	var oldJTI uuid.UUID
	oldJTI, err = validation.RequireRefreshJTI(&refreshToken.ID)
	if err != nil {
		return nil, err
	}

	var uid uuid.UUID
	uid, err = uuid.NewV7()
	if err != nil {
		return nil, apierr.FromService(span, inbounds.ErrGeneratingUUID{Cause: err, Location: "auth/refreshInternal"})
	}

	var newRefreshJTI = uid
	var refreshExp = time.Now().Add(7 * 24 * time.Hour)

	span.SetAttributes(attribute.String("old_token.id", oldJTI.String()))
	span.SetAttributes(attribute.String("new_token.id", newRefreshJTI.String()))

	var sess *session.Session
	sess, err = sessions.GetByFamilyID(ctx, refreshToken.Sub.FamilyID)
	if err != nil {
		return nil, apierr.FromService(span, inbounds.ErrSessionNotFound{})
	}

	now := time.Now()
	if sess.ExpiresAt.Before(now) || sess.RevokedAt != nil {
		// FIXME Record suspicious behaviour on audit when it is implemented
		return nil, apierr.FromService(span, inbounds.ErrSessionNotFound{})
	}

	// should revoke the session because of replay attacks
	// FIXME Add suspicious behaviour to audit when it is implemented
	if sess.TokenID != oldJTI {
		err = sessions.MarkRevokedByFamilyID(ctx, sess.FamilyID)
		if err != nil {
			apierr.RecordDomainError(span, err)
		}
		return nil, apierr.FromService(span, inbounds.ErrTokenReuseNotAllowed{TokenType: "refresh"})
	}

	sess, err = sessions.RotateToken(ctx, refreshToken.Sub.FamilyID, newRefreshJTI, oldJTI, refreshExp)
	if apierr.IsNotFound(err) {
		// sql.ErrNoRows → raced / reused / revoked
		return nil, apierr.FromService(span, inbounds.ErrSessionNotFound{})
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
		return nil, apierr.FromService(span, inbounds.ErrGeneratingUUID{Cause: err, Location: "auth/finishClientRefresh"})
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
		return nil, apierr.FromService(span, err)
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
		return nil, apierr.FromService(span, err)
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
		return nil, apierr.FromService(span, inbounds.ErrGeneratingUUID{Cause: err, Location: "auth/finishProjectUserRefresh"})
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
		return nil, apierr.FromService(span, err)
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
		return nil, apierr.FromService(span, err)
	}

	var refreshSig []byte
	refreshSig, err = keys.SignProject(ctx, *sess.ProjectID, refreshPayload)
	if err != nil {
		return nil, err
	}

	refreshTokenStr := issuer.AssembleJWT(refreshPayload, refreshSig)

	return &inbounds.UserTokensOutput{
		AccessTokenString:  accessTokenStr,
		RefreshTokenString: refreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
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
		return apierr.FromService(span, err)
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
		return outbounds.Email{}, apierr.FromService(span, inbounds.ErrEmptyFlowID{})
	}

	if in.SchemaType == "" {
		return outbounds.Email{}, apierr.FromService(span, inbounds.ErrEmptySchemaType{})
	}

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.FlowID = strings.TrimSpace(strings.ToLower(in.FlowID))
	in.SchemaType = strings.TrimSpace(strings.ToLower(in.SchemaType))

	if !schema.IsValidSchemaType(in.SchemaType) {
		return outbounds.Email{}, apierr.FromService(span, inbounds.ErrInvalidSchemaType{})
	}

	// FlowIDs cannot be the same as schema types so if this matches we error out
	if schema.IsValidSchemaType(in.FlowID) {
		return outbounds.Email{}, apierr.FromService(span, inbounds.ErrInvalidFlowID{Why: "flow id can't be the same as a schema type"})
	}

	if schema.Type(in.SchemaType) == schema.Core && schema.IsFlowIDReserved(in.FlowID) && in.CustomFields != nil {
		return outbounds.Email{}, apierr.FromService(span, inbounds.ErrCustomFieldsNotAllowed{})
	}

	empty := json.RawMessage(`{}`)
	customFields := &empty

	// Validate and construct metadata for non-core or non-reserved flows
	isCoreWithReservedFlow := schema.Type(in.SchemaType) == schema.Core && schema.IsFlowIDReserved(in.FlowID)
	if !isCoreWithReservedFlow {
		validatedMetadata, err := uc.validateAndConstructMetadata(ctx, span, in.ProjectID, schema.Type(in.SchemaType), in.FlowID, in.CustomFields)
		if err != nil {
			return outbounds.Email{}, err
		}
		if validatedMetadata != nil {
			customFields = validatedMetadata
		}
	}

	if len(in.Password) > 72 {
		return outbounds.Email{}, apierr.FromService(span, apierr.ErrPasswordTooLong{})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return outbounds.Email{}, apierr.FromService(span, inbounds.ErrHashingPassword{Cause: err})
	}

	var usr *project_users.ProjectUser
	usr, err = projectUsers.Register(ctx, project_users.ProjectUser{
		ProjectID:    in.ProjectID,
		Email:        in.Email,
		PasswordHash: string(hashedPassword),
		Metadata:     customFields,
	})
	if apierr.IsConflict(err) {
		return outbounds.Email{}, apierr.FromService(span, inbounds.ErrEmailAlreadyInUse{Cause: errors.New("email already in use")})
	} else if err != nil {
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
	verificationEmail, err = uc.mailRenderer.Verification(outbounds.VerificationEmailData{
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
	if apierr.IsNotFound(err) {
		return nil, apierr.FromService(span, inbounds.ErrInvalidCredentials{Cause: nil})
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.PasswordHash), []byte(in.Password))
	if err != nil {
		return nil, apierr.FromService(span, inbounds.ErrInvalidCredentials{Cause: err})
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
		return nil, apierr.FromService(span, inbounds.ErrGeneratingUUID{Cause: err, Location: "auth/LoginProjectUser"})
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
		return nil, apierr.FromService(span, err)
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
		return nil, apierr.FromService(span, err)
	}

	var refreshSig []byte
	refreshSig, err = keys.SignProject(ctx, in.ProjectID, refreshPayload)
	if err != nil {
		return nil, err
	}

	refreshTokenStr := issuer.AssembleJWT(refreshPayload, refreshSig)

	return &inbounds.UserTokensOutput{
		AccessTokenString:  accessTokenStr,
		RefreshTokenString: refreshTokenStr,
		AccessExpiresAt:    accessExpiresAt,
		RefreshExpiresAt:   refreshExpiresAt,
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
	principal, err = RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return apierr.FromService(span, err)
	}

	var claims *auth.VerificationClaims
	claims, err = uc.tokenVerifier.VerifyVerificationToken(ctx, token)
	if err != nil {
		return apierr.FromService(span, err)
	}

	if claims.Sub.Subject != principal.UserID {
		return apierr.FromService(span, inbounds.ErrTokenUserMismatch{TokenType: "verification"})
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
	principal, err = RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return apierr.FromService(span, err)
	}

	if principal.IsVerified == true {
		return apierr.FromService(span, inbounds.ErrUserAlreadyVerified{})
	}

	if principal.ProjectID != nil {
		u, err := projectUsers.GetByIDInternal(ctx, principal.UserID, *principal.ProjectID)
		if err != nil {
			return err
		}
		if u.IsVerified == true {
			return apierr.FromService(span, inbounds.ErrUserAlreadyVerified{})
		}
	} else {
		u, err := users.GetUserByID(ctx, principal.UserID)
		if err != nil {
			return err
		}
		if u.IsVerified == true {
			return apierr.FromService(span, inbounds.ErrUserAlreadyVerified{})
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
		return apierr.FromService(span, err)
	}

	var verificationSig []byte
	verificationSig, err = keys.SignGoAuth(ctx, verificationPayload)
	if err != nil {
		return err
	}

	verificationTokenStr := issuer.AssembleJWT(verificationPayload, verificationSig)

	var verificationEmail outbounds.Email
	verificationEmail, err = uc.mailRenderer.Verification(outbounds.VerificationEmailData{
		UserID: principal.UserID,
		Email:  principal.Email,
		Token:  verificationTokenStr,
		Locale: "en",
	})
	if err != nil {
		return apierr.FromService(span, err)
	}

	err = uc.mailSender.Send(ctx, verificationEmail)
	if err != nil {
		return apierr.FromService(span, err)
	}

	return nil
}

func (uc *UseCase) GetJWKS(ctx context.Context) (map[string]any, error) {
	keys, err := uc.deps.Keys.ListGoAuthPublicKeys(ctx)
	if err != nil {
		logs.L().Error("Failed listing GoAuth public keys", zap.Error(err))
		return nil, inbounds.ErrFailedToRetrieveJWKS{Cause: err}
	}

	jwkKeys := make([]any, len(keys))
	for i, k := range keys {
		jwkKeys[i] = key.PublicKeyToJWK(k)
	}

	return map[string]any{
		"keys": jwkKeys,
	}, nil
}
