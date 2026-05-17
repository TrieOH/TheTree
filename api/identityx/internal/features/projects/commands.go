package projects

import (
	"IdentityX/contracts"
	"IdentityX/internal/shared/authz"
	"IdentityX/internal/shared/feature_deps"
	"IdentityX/internal/shared/ports"
	"context"
	"fmt"
	"lib/crypto"
	"lib/database"
	"lib/errx"
	"time"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	encryptionKey []byte
	keyLifetime   time.Duration
	projects      ports.ProjectRepository
	keys          ports.KeysRepository
	logger        *zap.Logger
	tracer        trace.Tracer
	tx            database.TxRunner
}

func NewCommandService(deps feature_deps.ProjectCommandDeps) *CommandService {
	return errx.MustProvide(&CommandService{
		keyLifetime:   deps.KeyLifetime,
		encryptionKey: deps.EncryptionKey,
		projects:      deps.Projects,
		keys:          deps.Keys,
		logger:        deps.Logger,
		tracer:        deps.Tracer,
		tx:            deps.Tx,
	})
}

// Create handles the business logic for creating a new project.
// It requires a valid principal in the context, generates a new key pair for the project,
// and then creates the project in the database.
func (uc *CommandService) Create(ctx context.Context, in contracts.CreateProjectInput) (project *contracts.Project, err error) {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
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

func (uc *CommandService) createInternal(ctx context.Context, in contracts.CreateProjectInput) (*contracts.Project, error) {
	ctx, span := uc.tracer.Start(ctx, "ProjectService.createInternal")
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	var createdProject *contracts.Project
	if createdProject, err = uc.projects.Create(ctx, contracts.Project{
		ProjectName: in.ProjectName,
		Domain:      in.Domain,
		OwnerID:     principal.UserID,
		Metadata:    in.Metadata,
		IsActive:    true,
	}); err != nil {
		return nil, err
	}

	pub, priv, err := crypto.GenerateEd25519()
	if err != nil {
		return nil, fun.Errf("error generating project keys: %s", err).Internal()
	}
	defer zero(priv)

	encryptedPriv, err := crypto.Encrypt(priv, uc.encryptionKey)
	if err != nil {
		return nil, fun.Errf("error encrypting project keys: %s", err).Internal()
	}

	kid := fmt.Sprintf("project:%s:%s", createdProject.ID.String(), ulid.Make().String())

	expiresAt := time.Now().Add(uc.keyLifetime)
	if _, err = uc.keys.CreateKeyPair(ctx, contracts.Pair{
		KID:             kid,
		ProjectID:       &createdProject.ID,
		KeyType:         contracts.TypeProject,
		Algorithm:       contracts.AlgEdDSA,
		PublicKey:       pub,
		PrivateKey:      encryptedPriv,
		Usage:           contracts.UsageSign,
		Status:          contracts.StatusActive,
		ExpiresAt:       expiresAt,
		VerifyExpiresAt: expiresAt.Add(uc.keyLifetime),
	}); err != nil {
		return nil, err
	}

	span.SetAttributes(
		attribute.String("project.id", createdProject.ID.String()),
		attribute.String("project.owner_id", createdProject.OwnerID.String()),
		attribute.String("project.name", createdProject.ProjectName),
	)

	return createdProject, nil
}

// Update handles the business logic for updating a project.
// It requires a valid principal in the context and that the principal is the owner of the project.
// It retrieves the project, updates the fields, and then saves the changes to the database.
func (uc *CommandService) Update(ctx context.Context, in contracts.UpdateProjectInput) (*contracts.Project, error) {
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

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
