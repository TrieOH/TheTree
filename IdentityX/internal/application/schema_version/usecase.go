package schema_version

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/application/auth"
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/domain/version"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/MintzyG/fail"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	usecaseTracer = otel.Tracer("GoAuth.SchemaVersionService")
)

type UseCase struct {
	deps Deps
	tx   inbounds.TxRunner
}

type Deps struct {
	Schemas  outbounds.SchemaRepository
	Versions outbounds.SchemaVersionRepository
	Fields   outbounds.SchemaFieldsRepository
	Projects outbounds.ProjectRepository
}

var _ inbounds.SchemaVersionService = (*UseCase)(nil)

func New(
	deps Deps,
	tx inbounds.TxRunner,
) inbounds.SchemaVersionService {
	return &UseCase{
		deps: deps,
		tx:   tx,
	}
}

func (uc *UseCase) Draft(ctx context.Context, in inbounds.SchemaVersionServiceInput) (*inbounds.SchemaVersionOutput, error) {
	var out *inbounds.SchemaVersionOutput
	err := uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		var err error
		out, err = uc.draftInternal(ctx, in)
		return err
	})

	return out, err
}

func (uc *UseCase) draftInternal(ctx context.Context, in inbounds.SchemaVersionServiceInput) (*inbounds.SchemaVersionOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaVersionService.Draft")
	defer span.End()

	var err error
	defer func() {
		span.SetAttributes(attribute.Bool("draft.success", err == nil))
	}()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var principal *authz.Principal
	principal, err = auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(apierr.ProjectNotOwnedByPrincipal).WithArgs("cannot draft a schema version for a project you don't own")
	}

	var belongs bool
	belongs, err = schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return nil, err
	}

	if !belongs {
		return nil, fail.New(apierr.SchemaNotOwnedByPrincipal).WithArgs("cannot draft a schema version for a schema you don't own")
	}

	var latest *version.Version
	latest, err = versions.GetLatestForUpdate(ctx, in.SchemaID)

	if err != nil && !fail.Is(err, apierr.SQLNotFound) {
		return nil, err
	}

	if fail.Is(err, apierr.SQLNotFound) {
		newVersion := &version.Version{
			SchemaID:      in.SchemaID,
			VersionNumber: 1,
		}

		newVersion, err = versions.Draft(ctx, *newVersion)
		if err != nil {
			return nil, err
		}

		if err = schemas.SetVersion(ctx, schema.Schema{
			ID:               in.SchemaID,
			ProjectID:        in.ProjectID,
			CurrentVersionID: &newVersion.ID,
		}); err != nil {
			return nil, err
		}

		return inbounds.SchemaVersionToOutput(newVersion), nil
	}

	if latest.Status != version.StatusPublished {
		return nil, fail.New(apierr.SchemaVersionDraftOnNonPublished)
	}

	var newVersionDraft *version.Version
	newVersionDraft, err = versions.CopyOnDraft(ctx, latest.ID)
	if err != nil {
		return nil, err
	}

	err = fields.CloneFromTo(ctx, latest.ID, newVersionDraft.ID)
	if err != nil {
		return nil, err
	}

	return inbounds.SchemaVersionToOutput(newVersionDraft), nil
}

func (uc *UseCase) Publish(ctx context.Context, in inbounds.SchemaVersionServiceInput) error {
	ctx, span := usecaseTracer.Start(ctx, "SchemaVersionService.Publish")
	defer span.End()

	var err error
	defer func() {
		span.SetAttributes(attribute.Bool("publish.success", err == nil))
	}()

	projects := uc.deps.Projects
	schemas := uc.deps.Schemas
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	var principal *authz.Principal
	principal, err = auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return apierr.FromService(span, err)
	}

	var isOwner bool
	isOwner, err = projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return err
	}

	if !isOwner {
		return fail.New(apierr.ProjectNotOwnedByPrincipal).WithArgs("cannot publish a schema version for a project you don't own")
	}

	var belongs bool
	belongs, err = schemas.BelongsToProject(ctx, schema.Schema{
		ProjectID: in.ProjectID,
		ID:        in.SchemaID,
	})
	if err != nil {
		return err
	}

	if !belongs {
		return fail.New(apierr.SchemaNotOwnedByPrincipal).WithArgs("cannot publish a schema version for a schema you don't own")
	}

	var latest *version.Version
	latest, err = versions.GetLatest(ctx, in.SchemaID)
	if err != nil && !fail.Is(err, apierr.SQLNotFound) {
		return err
	}

	if err != nil && fail.Is(err, apierr.SQLNotFound) {
		return fail.New(apierr.SchemaVersionDraftDoesntExist)
	}

	if latest.Status != version.StatusDraft {
		if latest.Status == version.StatusPublished {
			err = fail.New(apierr.SchemaVersionTryingToPublishPublished)
		} else if latest.Status == version.StatusArchived {
			err = fail.New(apierr.SchemaVersionTryingToPublishArchived)
		} else {
			err = fail.New(apierr.SchemaVersionNoValidStatus)
		}
		return err
	}

	var hasFields bool
	hasFields, err = versions.HasFields(ctx, latest.ID)
	if err != nil {
		return err
	}

	if !hasFields {
		return fail.New(apierr.SchemaVersionPublishWithNoFields)
	}

	if latest.BasedOnVersionID == nil {
		if err = versions.Publish(ctx, version.Version{
			SchemaID: in.SchemaID,
			ID:       latest.ID,
		}); err != nil {
			return err
		}

		if err = schemas.SetVersion(ctx, schema.Schema{
			ID:               in.SchemaID,
			ProjectID:        in.ProjectID,
			CurrentVersionID: &latest.ID,
		}); err != nil {
			return err
		}

		return nil
	}

	diff, err := fields.DiffVersionsFullState(ctx, *latest.BasedOnVersionID, latest.ID)
	if err != nil {
		return err
	}

	diff.Annotate(span)

	if !diff.HasAnyChanges() {
		return fail.New(apierr.SchemaVersionNoChanges)
	}

	if err = versions.Publish(ctx, version.Version{
		SchemaID: in.SchemaID,
		ID:       latest.ID,
	}); err != nil {
		return err
	}

	if err = schemas.SetVersion(ctx, schema.Schema{
		ID:               in.SchemaID,
		ProjectID:        in.ProjectID,
		CurrentVersionID: &latest.ID,
	}); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) GetCurrent(ctx context.Context, in inbounds.SchemaVersionServiceInput) (*inbounds.SchemaVersionOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaVersionService.GetCurrent",
		trace.WithAttributes(
			attribute.String("project_id", in.ProjectID.String()),
			attribute.String("schema_id", in.SchemaID.String()),
		),
	)
	defer span.End()

	projects := uc.deps.Projects
	versions := uc.deps.Versions

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	isOwner, err := projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(apierr.ProjectNotOwnedByPrincipal).WithArgs("cannot get the current schema version for a project you don't own")
	}

	current, err := versions.GetCurrent(ctx, in.SchemaID)
	if err != nil {
		return nil, err
	}

	return inbounds.SchemaVersionToOutput(current), nil
}

func (uc *UseCase) GetLatest(ctx context.Context, in inbounds.SchemaVersionServiceInput) (*inbounds.SchemaVersionOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaVersionService.GetLatest",
		trace.WithAttributes(
			attribute.String("project_id", in.ProjectID.String()),
			attribute.String("schema_id", in.SchemaID.String()),
		),
	)
	defer span.End()

	projects := uc.deps.Projects
	versions := uc.deps.Versions

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	isOwner, err := projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(apierr.ProjectNotOwnedByPrincipal).WithArgs("cannot get the latest schema version for a project you don't own")
	}

	latest, err := versions.GetLatest(ctx, in.SchemaID)
	if err != nil {
		return nil, err
	}

	return inbounds.SchemaVersionToOutput(latest), nil
}

func (uc *UseCase) GetVerbose(ctx context.Context, in inbounds.SchemaVersionServiceInput) (*inbounds.VersionVerboseOutput, error) {
	ctx, span := usecaseTracer.Start(ctx, "SchemaVersionService.GetVerbose",
		trace.WithAttributes(
			attribute.String("project_id", in.ProjectID.String()),
			attribute.String("schema_id", in.SchemaID.String()),
			attribute.Int("version", in.VersionNumber),
		),
	)
	defer span.End()

	projects := uc.deps.Projects
	versions := uc.deps.Versions
	fields := uc.deps.Fields

	principal, err := auth.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, apierr.FromService(span, err)
	}

	isOwner, err := projects.IsOwnerOf(ctx, in.ProjectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	if !isOwner {
		return nil, fail.New(apierr.ProjectNotOwnedByPrincipal).WithArgs("cannot get the latest schema version for a project you don't own")
	}

	var foundVersion *version.Version
	if in.VersionID == nil {
		foundVersion, err = versions.GetByVersionNumber(ctx, in.SchemaID, in.VersionNumber)
		if err != nil {
			return nil, err
		}
	} else {
		foundVersion, err = versions.GetByID(ctx, *in.VersionID)
		if err != nil {
			return nil, err
		}
	}

	versionFields, err := fields.ListFromVersion(ctx, in.SchemaID, foundVersion.ID)
	if err != nil {
		return nil, err
	}

	out := inbounds.VersionVerboseOutput{
		SchemaVersionOutput: inbounds.SchemaVersionOutput{
			ID:               foundVersion.ID,
			SchemaID:         foundVersion.SchemaID,
			BasedOnVersionID: foundVersion.BasedOnVersionID,
			VersionNumber:    foundVersion.VersionNumber,
			Status:           foundVersion.Status,
			CreatedAt:        foundVersion.CreatedAt,
			UpdatedAt:        foundVersion.UpdatedAt,
		},
		Fields: nil,
	}

	for _, f := range versionFields {
		out.Fields = append(out.Fields, inbounds.OutputField{
			ObjectID:        f.ObjectID,
			ID:              f.ID,
			Key:             f.Key,
			SchemaID:        f.SchemaID,
			SchemaVersionID: f.SchemaVersionID,
			Type:            string(f.Type),
			Owner:           string(f.Owner),
			Title:           f.Title,
			Description:     f.Description,
			Placeholder:     f.Placeholder,
			Required:        f.Required,
			Mutable:         f.Mutable,
			DefaultValue:    f.DefaultValue,
			Position:        f.Position,
			CreatedAt:       f.CreatedAt,
			UpdatedAt:       f.UpdatedAt,
		})
	}

	return &out, nil
}
