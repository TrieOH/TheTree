package project

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/crypto"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/key"
	"GoAuth/internal/domain/project"
	"GoAuth/internal/domain/scopes"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"
	"fmt"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	usecaseTracer = otel.Tracer("auth_usecase")
)

type UseCase struct {
	projects outbounds.ProjectRepository
	scopes   outbounds.ScopeRepository
	keys     outbounds.KeysRepository
	tx       inbounds.TxRunner
}

var _ inbounds.ProjectService = (*UseCase)(nil)

func New(
	projects outbounds.ProjectRepository,
	scopes outbounds.ScopeRepository,
	keys outbounds.KeysRepository,
	tx inbounds.TxRunner,
) inbounds.ProjectService {
	return &UseCase{
		projects: projects,
		scopes:   scopes,
		keys:     keys,
		tx:       tx,
	}
}

// Create handles the business logic for creating a new project.
// It requires a valid principal in the context, generates a new key pair for the project,
// and then creates the project in the database.
func (uc *UseCase) Create(ctx context.Context, in inbounds.ProjectServiceInput) (project *inbounds.OutputProject, err error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.Create")
	defer span.End()
	defer func() {
		if span != nil {
			span.SetAttributes(attribute.Bool("create.success", err == nil))
		}
	}()

	err = uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		project, err = uc.createInternal(ctx, in)
		return err
	})
	if err != nil {
		return nil, err
	}

	return
}

func (uc *UseCase) createInternal(ctx context.Context, in inbounds.ProjectServiceInput) (*inbounds.OutputProject, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.createInternal")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, fail.New(apierr.ProjectErrorGeneratingKeys).With(err).RecordCtx(ctx)
	}

	createdProject, err := uc.projects.Create(ctx, project.Project{
		ProjectName: in.ProjectName,
		OwnerID:     principal.UserID,
		Metadata:    in.Metadata,
		IsActive:    true,
	})
	if err != nil {
		return nil, err
	}

	pub, priv, err := crypto.GenerateEd25519()
	if err != nil {
		return nil, fail.New(apierr.ProjectErrorGeneratingKeys).With(err).RecordCtx(ctx)
	}
	defer zero(priv)

	encryptedPriv, err := crypto.Encrypt(priv)
	if err != nil {
		return nil, fail.New(apierr.ProjectErrorGeneratingKeys).With(err).RecordCtx(ctx)
	}

	kid := fmt.Sprintf(
		"project:%s:%s",
		createdProject.ID.String(),
		ulid.Make().String(),
	)

	_, err = uc.keys.CreateKeyPair(ctx, key.Pair{
		KID:        kid,
		ProjectID:  &createdProject.ID,
		KeyType:    key.TypeProject,
		Algorithm:  key.AlgEdDSA,
		PublicKey:  pub,
		PrivateKey: encryptedPriv,
		Usage:      key.UsageSign,
		Status:     key.StatusActive,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		return nil, fail.New(apierr.ProjectErrorGeneratingKeys).With(err).RecordCtx(ctx)
	}

	span.SetAttributes(
		attribute.String("project.id", createdProject.ID.String()),
		attribute.String("project.owner_id", createdProject.OwnerID.String()),
		attribute.String("project.name", createdProject.ProjectName),
	)

	createdScope, err := uc.scopes.Create(ctx, scopes.Scope{
		Type:      scopes.ScopeTypeProjectRoot,
		ProjectID: &createdProject.ID,
	})
	if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("project.scope.id", createdScope.ID.String()),
		attribute.String("project.scope.type", string(createdScope.Type)),
	)

	return inbounds.OutputProjectFromProject(createdProject), nil
}

// GetByID handles the business logic for retrieving a project by its ID.
// It requires a valid principal in the context and that the principal is the owner of the project.
func (uc *UseCase) GetByID(ctx context.Context, projectID uuid.UUID) (*inbounds.OutputProject, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.GetByID",
		trace.WithAttributes(attribute.String("project.id", projectID.String())),
	)
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	if principal.ProjectID != nil && *principal.ProjectID != projectID {
		return nil, fail.New(apierr.ProjectNotFound).RecordCtx(ctx)
	}

	proj, err := uc.projects.GetByIDExternal(ctx, projectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("project.owner_id", proj.OwnerID.String()),
		attribute.String("project.name", proj.ProjectName),
	)

	return inbounds.OutputProjectFromProject(proj), nil
}

// List handles the business logic for listing all projects for the authenticated user.
func (uc *UseCase) List(ctx context.Context) ([]inbounds.OutputProject, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.List")
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

	return inbounds.OutputProjectSliceFromProjectSlice(projects), nil
}

// GetJWKS handles the business logic for retrieving the JWKS for a project.
// It retrieves the public key for the project and converts it to a JWK set.
func (uc *UseCase) GetJWKS(ctx context.Context, projectID uuid.UUID) (map[string]any, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.GetJWKS")
	defer span.End()

	keys, err := uc.keys.ListProjectPublicKeys(ctx, projectID)
	if err != nil {
		return map[string]any{"keys": []any{}}, err
	}

	jwkKeys := make([]any, len(keys))
	for i, k := range keys {
		jwkKeys[i] = key.PublicKeyToJWK(k)
	}

	return map[string]any{
		"keys": jwkKeys,
	}, nil
}

// Update handles the business logic for updating a project.
// It requires a valid principal in the context and that the principal is the owner of the project.
// It retrieves the project, updates the fields, and then saves the changes to the database.
func (uc *UseCase) Update(ctx context.Context, in inbounds.ProjectServiceInput) (*inbounds.OutputProject, error) {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.Update")
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

	updatedProject, err := uc.projects.Update(ctx, project.Project{
		ID:          newProject.ID,
		ProjectName: newProject.ProjectName,
		Metadata:    newProject.Metadata,
	},
		principal.UserID,
	)
	if err != nil {
		return nil, err
	}

	return inbounds.OutputProjectFromProject(updatedProject), nil
}

// Delete handles the business logic for deleting a project.
// It requires a valid principal in the context and that the principal is the owner of the project.
func (uc *UseCase) Delete(ctx context.Context, projectID uuid.UUID) error {
	ctx, span := usecaseTracer.Start(ctx, "ProjectService.Delete",
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

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
