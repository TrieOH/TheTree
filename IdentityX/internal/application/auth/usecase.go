package auth

import (
	"GoAuth/internal/adapters/email"
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/tokens"
	"GoAuth/internal/application/validation"
	"GoAuth/internal/domain/auth"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/project_users"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/session"
	"GoAuth/internal/domain/user"
	"GoAuth/internal/infrastructure"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"
	"crypto/ed25519"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"GoAuth/internal/utils"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

var (
	usecaseTracer = otel.Tracer("auth_usecase")
)

type UseCase struct {
	deps          Deps
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
}

var _ inbounds.AuthService = (*UseCase)(nil)

func New(
	repos Deps,
	infra infrastructure.Infra,
	tokenBundle tokens.TokenBundle,
	mailBundle email.MailBundle,
) inbounds.AuthService {
	return &UseCase{
		deps:          repos,
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

	var verificationToken string
	verificationToken, err = uc.tokenIssuer.NewVerificationToken(inbounds.NewVerificationTokenInput{
		Subject:    u.ID,
		ExpiresAt:  time.Now().Add(15 * time.Minute),
		PrivateKey: utils.GoAuthPrivateKey,
	})
	if err != nil {
		return outbounds.Email{}, err
	}

	var verificationEmail outbounds.Email
	verificationEmail, err = uc.mailRenderer.Verification(outbounds.VerificationEmailData{
		UserID: u.ID,
		Email:  u.Email,
		Token:  verificationToken,
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
func (uc *UseCase) Login(ctx context.Context, in inbounds.LoginUserInput) (*inbounds.UserTokensOutput, error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))

	var err error
	ctx, span := usecaseTracer.Start(ctx, "AuthService.Login")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("login.success", err == nil))
		}
	}()

	users := uc.deps.Users
	sessions := uc.deps.Sessions

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

	var accessToken string
	var accessJTI uuid.UUID
	accessJTI, err = uuid.NewV7()
	if err != nil {
		return nil, apierr.FromService(span, inbounds.ErrGeneratingUUID{Cause: err, Location: "auth/login"})
	}
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	accessToken, err = uc.tokenIssuer.NewAccessToken(inbounds.NewAccessTokenInput{
		User:       *u,
		PrivateKey: utils.GoAuthPrivateKey,
		IP:         in.IP,
		Agent:      in.Agent,
		AccessJTI:  accessJTI.String(),
		KeyID:      "goauth:v1",
		SessionID:  sess.SessionID,
		ExpiresAt:  accessExpiresAt,
	})
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	var refreshToken string
	refreshToken, err = uc.tokenIssuer.NewRefreshToken(inbounds.NewRefreshTokenInput{
		KeyID:      "goauth:v1",
		PrivateKey: utils.GoAuthPrivateKey,
		AccessJTI:  accessJTI,
		RefreshJTI: sess.TokenID,
		ExpiresAt:  refreshExpiresAt,
		FamilyID:   sess.FamilyID,
	})
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	return &inbounds.UserTokensOutput{
		AccessTokenString:  accessToken,
		RefreshTokenString: refreshToken,
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
) (*inbounds.UserTokensOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.finishClientRefresh")
	defer span.End()

	var err error
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("finishClientRefresh.success", err == nil))
		}
	}()

	users := uc.deps.Users
	sessions := uc.deps.Sessions

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

	var uid uuid.UUID
	uid, err = uuid.NewV7()
	if err != nil {
		return nil, apierr.FromService(span, inbounds.ErrGeneratingUUID{Cause: err, Location: "auth/finishClientRefresh"})
	}

	newAccessJTI := uid
	var accessTokenStr string
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	accessTokenStr, err = uc.tokenIssuer.NewAccessToken(inbounds.NewAccessTokenInput{
		User:       *u,
		PrivateKey: utils.GoAuthPrivateKey,
		IP:         in.IP,
		Agent:      in.Agent,
		AccessJTI:  newAccessJTI.String(),
		KeyID:      "goauth:v1",
		SessionID:  sess.SessionID,
		ExpiresAt:  accessExpiresAt,
	})
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	var refreshTokenStr string
	refreshTokenStr, err = uc.tokenIssuer.NewRefreshToken(inbounds.NewRefreshTokenInput{
		KeyID:      "goauth:v1",
		PrivateKey: utils.GoAuthPrivateKey,
		AccessJTI:  newAccessJTI,
		RefreshJTI: refreshJTI,
		ExpiresAt:  refreshExpiresAt,
		FamilyID:   sess.FamilyID,
	})
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

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
	projects := uc.deps.Projects
	sessions := uc.deps.Sessions

	var identity *session.Identity
	identity, err = sessions.GetIdentityByIDAndType(ctx, sess.IdentityID, session.ProjectIdentity)
	if err != nil {
		return nil, err
	}

	projectUser, err := projectUsers.GetByIDInternal(ctx, identity.EntityID, *sess.ProjectID)
	if err != nil {
		return nil, err
	}

	var privKey string
	privKey, err = projects.GetPrivateKeyByIDInternal(ctx, *sess.ProjectID)
	if err != nil {
		return nil, err
	}

	var decodedKey ed25519.PrivateKey
	decodedKey, err = utils.ParseEd25519PrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	var uid uuid.UUID
	uid, err = uuid.NewV7()
	if err != nil {
		return nil, apierr.FromService(span, inbounds.ErrGeneratingUUID{Cause: err, Location: "auth/finishProjectUserRefresh"})
	}

	newAccessJTI := uid
	var keyID = "project:" + sess.ProjectID.String() + ":v1"
	var accessTokenStr string
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	accessTokenStr, err = uc.tokenIssuer.NewProjectAccessToken(inbounds.NewProjectAccessTokenInput{
		User:       *projectUser,
		IP:         in.IP,
		Agent:      in.Agent,
		AccessJTI:  newAccessJTI.String(),
		KeyID:      keyID,
		SessionID:  sess.SessionID,
		ExpiresAt:  accessExpiresAt,
		PrivateKey: decodedKey,
	})
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	var refreshTokenStr string
	refreshTokenStr, err = uc.tokenIssuer.NewRefreshToken(inbounds.NewRefreshTokenInput{
		KeyID:      keyID,
		PrivateKey: decodedKey,
		AccessJTI:  newAccessJTI,
		RefreshJTI: refreshJTI,
		ExpiresAt:  refreshExpiresAt,
		FamilyID:   sess.FamilyID,
	})
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

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

	var verificationToken string
	verificationToken, err = uc.tokenIssuer.NewVerificationToken(inbounds.NewVerificationTokenInput{
		Subject:    usr.ID,
		ExpiresAt:  time.Now().Add(15 * time.Minute),
		PrivateKey: utils.GoAuthPrivateKey,
	})
	if err != nil {
		return outbounds.Email{}, err
	}

	var verificationEmail outbounds.Email
	verificationEmail, err = uc.mailRenderer.Verification(outbounds.VerificationEmailData{
		UserID: usr.ID,
		Email:  usr.Email,
		Token:  verificationToken,
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
	*inbounds.UserTokensOutput,
	error,
) {
	ctx, span := usecaseTracer.Start(ctx, "AuthService.LoginProjectUser",
		trace.WithAttributes(attribute.String("project.id", in.ProjectID.String())),
	)
	defer span.End()

	projectUsers := uc.deps.ProjectUsers
	projects := uc.deps.Projects
	sessions := uc.deps.Sessions

	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	usr, err := projectUsers.GetByEmailInternal(ctx, in.ProjectID, in.Email)
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
	sess, err := sessions.Create(ctx, session.Session{
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

	var privKey string
	privKey, err = projects.GetPrivateKeyByIDInternal(ctx, *sess.ProjectID)
	if err != nil {
		return nil, err
	}

	var decodedKey ed25519.PrivateKey
	decodedKey, err = utils.ParseEd25519PrivateKey(privKey)
	if err != nil {
		return nil, err
	}

	var uid uuid.UUID
	uid, err = uuid.NewV7()
	if err != nil {
		return nil, apierr.FromService(span, inbounds.ErrGeneratingUUID{Cause: err, Location: "auth/LoginProjectUser"})
	}

	keyID := "project:" + sess.ProjectID.String() + ":v1"
	accessJTI := uid
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	accessToken, err := uc.tokenIssuer.NewProjectAccessToken(inbounds.NewProjectAccessTokenInput{
		User:       *usr,
		IP:         in.IP,
		Agent:      in.Agent,
		AccessJTI:  accessJTI.String(),
		KeyID:      keyID,
		SessionID:  sess.SessionID,
		ExpiresAt:  accessExpiresAt,
		PrivateKey: decodedKey,
	})
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	refreshToken, err := uc.tokenIssuer.NewRefreshToken(inbounds.NewRefreshTokenInput{
		KeyID:      keyID,
		PrivateKey: decodedKey,
		AccessJTI:  accessJTI,
		RefreshJTI: sess.TokenID,
		ExpiresAt:  refreshExpiresAt,
		FamilyID:   sess.FamilyID,
	})
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	return &inbounds.UserTokensOutput{
		AccessTokenString:  accessToken,
		RefreshTokenString: refreshToken,
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
	claims, err = uc.tokenVerifier.VerifyVerificationToken(token)

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
