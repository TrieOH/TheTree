package projects

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/contracts"
	crypto2 "IdentityX/internal/shared/crypto"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	users    ports.UserRepository
	projects ports.ProjectRepository
	keys     ports.KeysRepository
	logger   *zap.Logger
	tracer   trace.Tracer
	txRunner database.TxRunner
}

func NewCommandService(
	users ports.UserRepository,
	projects ports.ProjectRepository,
	keys ports.KeysRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	txRunner database.TxRunner,
) *CommandService {
	return &CommandService{
		users:    users,
		projects: projects,
		keys:     keys,
		logger:   logger,
		tracer:   tracer,
		txRunner: txRunner,
	}
}

type ProjectServiceInput struct {
	ProjectID   uuid.UUID
	ProjectName string
	Domain      string
	Metadata    json.RawMessage
}

// Create handles the business logic for creating a new project.
// It requires a valid principal in the context, generates a new key pair for the project,
// and then creates the project in the database.
func (uc *CommandService) Create(ctx context.Context, in ProjectServiceInput) (project *contracts.Project, err error) {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.Create")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("create.success", err == nil))
		}
	}()

	err = uc.txRunner.WithinTx(ctx, func(ctx context.Context) error {
		project, err = uc.createInternal(ctx, in)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

func (uc *CommandService) createInternal(ctx context.Context, in ProjectServiceInput) (*contracts.Project, error) {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.createInternal")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, fail.New(errx.ProjectErrorGeneratingKeys).With(err).RecordCtx(ctx)
	}

	createdProject, err := uc.projects.Create(ctx, contracts.Project{
		ProjectName: in.ProjectName,
		Domain:      in.Domain,
		OwnerID:     principal.UserID,
		Metadata:    in.Metadata,
		IsActive:    true,
	})
	if err != nil {
		return nil, err
	}

	pub, priv, err := crypto2.GenerateEd25519()
	if err != nil {
		return nil, fail.New(errx.ProjectErrorGeneratingKeys).With(err).RecordCtx(ctx)
	}
	defer zero(priv)

	encryptedPriv, err := crypto2.Encrypt(priv)
	if err != nil {
		return nil, fail.New(errx.ProjectErrorGeneratingKeys).With(err).RecordCtx(ctx)
	}

	kid := fmt.Sprintf(
		"project:%s:%s",
		createdProject.ID.String(),
		ulid.Make().String(),
	)

	expiresAt := time.Now().Add(viper.GetDuration("GOAUTH_KEY_LIFETIME"))
	_, err = uc.keys.CreateKeyPair(ctx, contracts.Pair{
		KID:             kid,
		ProjectID:       &createdProject.ID,
		KeyType:         contracts.TypeProject,
		Algorithm:       contracts.AlgEdDSA,
		PublicKey:       pub,
		PrivateKey:      encryptedPriv,
		Usage:           contracts.UsageSign,
		Status:          contracts.StatusActive,
		ExpiresAt:       expiresAt,
		VerifyExpiresAt: expiresAt.Add(viper.GetDuration("GOAUTH_KEY_LIFETIME")),
	})
	if err != nil {
		return nil, fail.New(errx.ProjectErrorGeneratingKeys).With(err).RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("project.id", createdProject.ID.String()),
		attribute.String("project.owner_id", createdProject.OwnerID.String()),
		attribute.String("project.name", createdProject.ProjectName),
	)

	return createdProject, nil
}

// GetByID handles the business logic for retrieving a project by its ID.
// It requires a valid principal in the context and that the principal is the owner of the project.
func (uc *CommandService) GetByID(ctx context.Context, projectID uuid.UUID) (*contracts.Project, error) {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.GetByID",
		trace.WithAttributes(attribute.String("project.id", projectID.String())),
	)
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	if principal.ProjectID != nil && *principal.ProjectID != projectID {
		return nil, fail.New(errx.ProjectNotFound).RecordCtx(ctx)
	}

	proj, err := uc.projects.GetByIDExternal(ctx, projectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("project.owner_id", proj.OwnerID.String()),
		attribute.String("project.name", proj.ProjectName),
	)

	return proj, nil
}

// List handles the business logic for listing all projects for the authenticated user.
func (uc *CommandService) List(ctx context.Context) ([]contracts.Project, error) {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.List")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	projects, err := uc.projects.List(ctx, principal.UserID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("projects.count", len(projects)))

	return projects, nil
}

// GetJWKS handles the business logic for retrieving the JWKS for a project.
// It retrieves the public key for the project and converts it to a JWK set.
func (uc *CommandService) GetJWKS(ctx context.Context, projectID uuid.UUID) (map[string]any, error) {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.GetJWKS")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	if principal.ProjectID != nil && *principal.ProjectID != projectID {
		return nil, fail.New(errx.ProjectNotFound).RecordCtx(ctx)
	}

	if principal.ProjectID == nil {
		isOwner, err := uc.projects.IsOwnerOf(ctx, projectID, principal.UserID)
		if err != nil {
			return nil, err
		}
		if !isOwner {
			return nil, fail.New(errx.ProjectNotFound).RecordCtx(ctx)
		}
	}

	keys, err := uc.keys.ListProjectPublicKeys(ctx, projectID)
	if err != nil {
		return map[string]any{"security": []any{}}, err
	}

	jwkKeys := make([]any, 0, len(keys))
	var jwk map[string]any
	for _, k := range keys {
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

// Update handles the business logic for updating a project.
// It requires a valid principal in the context and that the principal is the owner of the project.
// It retrieves the project, updates the fields, and then saves the changes to the database.
func (uc *CommandService) Update(ctx context.Context, in ProjectServiceInput) (*contracts.Project, error) {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.Update")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.String("project.id", in.ProjectID.String()))

	newProject, err := uc.projects.GetByIDExternal(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if in.ProjectName != "" {
		newProject.ProjectName = in.ProjectName
	}
	newProject.Metadata = in.Metadata

	if in.Domain != "" {
		newProject.Domain = in.Domain
	}

	updatedProject, err := uc.projects.Update(ctx, contracts.Project{
		ID:          newProject.ID,
		ProjectName: newProject.ProjectName,
		Domain:      newProject.Domain,
		Metadata:    newProject.Metadata,
	},
		principal.UserID,
	)
	if err != nil {
		return nil, err
	}

	return updatedProject, nil
}

// Delete handles the business logic for deleting a project.
// It requires a valid principal in the context and that the principal is the owner of the project.
func (uc *CommandService) Delete(ctx context.Context, projectID uuid.UUID) error {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.Delete",
		trace.WithAttributes(attribute.String("project.id", projectID.String())),
	)
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	err = uc.projects.Delete(ctx, projectID, principal.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *CommandService) ListUsers(ctx context.Context, projectID uuid.UUID) ([]contracts.User, error) {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.ListUsers",
		trace.WithAttributes(attribute.String("project.id", projectID.String())),
	)
	defer span.End()

	users, err := uc.users.ListFromProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(attribute.Int("users.count", len(users)))

	return users, nil
}

func (uc *CommandService) GetUser(ctx context.Context, projectID, userID uuid.UUID) (*contracts.User, error) {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.GetUser",
		trace.WithAttributes(
			attribute.String("project.id", projectID.String()),
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	user, err := uc.users.GetByIDFromProject(ctx, userID, projectID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
