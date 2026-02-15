package subcontext

import (
	"GoAuth/internal/domain/authz"
	"GoAuth/internal/errx"
	"GoAuth/internal/ports/inbounds"
	"GoAuth/internal/ports/outbounds"
	"context"
	"encoding/json"
	"strings"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	usecaseTracer = otel.Tracer("GoAuth.SubContextService")
)

type UseCase struct {
	deps Deps
	tx   inbounds.TxRunner
}

type Deps struct {
	Projects     outbounds.ProjectRepository
	ProjectUsers outbounds.ProjectUserRepository
}

var _ inbounds.SubContextService = (*UseCase)(nil)

func New(
	deps Deps,
	tx inbounds.TxRunner,
) inbounds.SubContextService {
	return &UseCase{
		deps: deps,
		tx:   tx,
	}
}

func (uc *UseCase) Add(ctx context.Context, projectID, userID uuid.UUID, data json.RawMessage) error {
	ctx, span := usecaseTracer.Start(ctx, "SubContextService.Add",
		trace.WithAttributes(
			attribute.String("project.id", projectID.String()),
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	// Only clients (project owners) can manage sub-context
	if principal.ProjectID != nil {
		return fail.New(errx.AuthNotClient).RecordCtx(ctx)
	}

	// Verify the client owns the project
	_, err = uc.deps.Projects.GetByIDExternal(ctx, projectID, principal.UserID)
	if err != nil {
		return err
	}

	// Verify user exists in project
	_, err = uc.deps.ProjectUsers.GetByIDInternal(ctx, userID, projectID)
	if err != nil {
		return err
	}

	// Get current user metadata
	user, err := uc.deps.ProjectUsers.GetByIDInternal(ctx, userID, projectID)
	if err != nil {
		return err
	}

	// Parse current metadata
	var metadataMap map[string]any
	if user.Metadata != nil && len(*user.Metadata) > 0 {
		if err := json.Unmarshal(*user.Metadata, &metadataMap); err != nil {
			return fail.New(errx.ProjectUserMetadataDecodeFailed).With(err).RecordCtx(ctx)
		}
	} else {
		metadataMap = make(map[string]any)
	}

	// Parse new sub-context data
	var newSubContext map[string]any
	if err := json.Unmarshal(data, &newSubContext); err != nil {
		return fail.New(errx.RequestInvalidSubContextJSON).With(err).RecordCtx(ctx)
	}

	// Get existing sub-context or create new
	var existingSubContext map[string]any
	if existing, ok := metadataMap["sub-context"].(map[string]any); ok {
		existingSubContext = existing
	} else {
		existingSubContext = make(map[string]any)
	}

	// Merge new data into existing sub-context
	for k, v := range newSubContext {
		existingSubContext[k] = v
	}

	// Update metadata with merged sub-context
	metadataMap["sub-context"] = existingSubContext

	// Marshal updated metadata
	updatedMetadata, err := json.Marshal(metadataMap)
	if err != nil {
		return fail.New(errx.ProjectUserErrorEncodingMetadata).With(err).RecordCtx(ctx)
	}

	rawMetadata := json.RawMessage(updatedMetadata)
	if err := uc.deps.ProjectUsers.UpdateMetadata(ctx, userID, projectID, &rawMetadata); err != nil {
		return err
	}

	span.SetAttributes(attribute.Bool("sub_context.added", true))
	return nil
}

func (uc *UseCase) Remove(ctx context.Context, projectID, userID uuid.UUID, keys []string) error {
	ctx, span := usecaseTracer.Start(ctx, "SubContextService.Remove",
		trace.WithAttributes(
			attribute.String("project.id", projectID.String()),
			attribute.String("user.id", userID.String()),
			attribute.Int("keys.count", len(keys)),
		),
	)
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return err
	}

	// Only clients (project owners) can manage sub-context
	if principal.ProjectID != nil {
		return fail.New(errx.AuthNotClient).RecordCtx(ctx)
	}

	// Verify the client owns the project
	_, err = uc.deps.Projects.GetByIDExternal(ctx, projectID, principal.UserID)
	if err != nil {
		return err
	}

	// Get current user metadata
	user, err := uc.deps.ProjectUsers.GetByIDInternal(ctx, userID, projectID)
	if err != nil {
		return err
	}

	// Parse current metadata
	var metadataMap map[string]any
	if user.Metadata != nil && len(*user.Metadata) > 0 {
		if err := json.Unmarshal(*user.Metadata, &metadataMap); err != nil {
			return fail.New(errx.ProjectUserMetadataDecodeFailed).With(err).RecordCtx(ctx)
		}
	} else {
		// No metadata, nothing to remove
		return nil
	}

	// Get existing sub-context
	subContext, ok := metadataMap["sub-context"].(map[string]any)
	if !ok || subContext == nil {
		// No sub-context, nothing to remove
		return nil
	}

	// Remove specified keys (supports nested paths like "permissions.can_delete")
	for _, key := range keys {
		removeKeyFromMap(subContext, key)
	}

	// Update metadata
	metadataMap["sub-context"] = subContext

	// Marshal updated metadata
	updatedMetadata, err := json.Marshal(metadataMap)
	if err != nil {
		return fail.New(errx.ProjectUserErrorEncodingMetadata).With(err).RecordCtx(ctx)
	}

	rawMetadata := json.RawMessage(updatedMetadata)
	if err := uc.deps.ProjectUsers.UpdateMetadata(ctx, userID, projectID, &rawMetadata); err != nil {
		return err
	}

	span.SetAttributes(attribute.Bool("sub_context.removed", true))
	return nil
}

func removeKeyFromMap(m map[string]any, key string) {
	parts := strings.Split(key, ".")
	if len(parts) == 1 {
		delete(m, key)
		return
	}

	// Navigate nested structure
	current := m
	for i := 0; i < len(parts)-1; i++ {
		next, ok := current[parts[i]].(map[string]any)
		if !ok {
			return // Path doesn't exist
		}
		current = next
	}

	delete(current, parts[len(parts)-1])
}

func (uc *UseCase) Get(ctx context.Context, projectID, userID uuid.UUID) (json.RawMessage, error) {
	ctx, span := usecaseTracer.Start(ctx, "SubContextService.Get",
		trace.WithAttributes(
			attribute.String("project.id", projectID.String()),
			attribute.String("user.id", userID.String()),
		),
	)
	defer span.End()

	principal, err := authz.RequirePrincipalAndAnnotate(ctx, span)
	if err != nil {
		return nil, err
	}

	// Only clients (project owners) can view sub-context
	if principal.ProjectID != nil {
		return nil, fail.New(errx.AuthNotClient).RecordCtx(ctx)
	}

	// Verify the client owns the project
	_, err = uc.deps.Projects.GetByIDExternal(ctx, projectID, principal.UserID)
	if err != nil {
		return nil, err
	}

	// Get user
	user, err := uc.deps.ProjectUsers.GetByIDInternal(ctx, userID, projectID)
	if err != nil {
		return nil, err
	}

	// Parse metadata to extract sub-context
	var metadataMap map[string]any
	if user.Metadata != nil && len(*user.Metadata) > 0 {
		if err := json.Unmarshal(*user.Metadata, &metadataMap); err != nil {
			return nil, fail.New(errx.ProjectUserMetadataDecodeFailed).With(err).RecordCtx(ctx)
		}
	}

	// Extract sub-context
	subContext, ok := metadataMap["sub-context"].(map[string]any)
	if !ok || subContext == nil {
		// Return empty object if no sub-context exists
		return json.RawMessage("{}"), nil
	}

	// Marshal sub-context back to JSON
	result, err := json.Marshal(subContext)
	if err != nil {
		return nil, fail.New(errx.ProjectUserErrorEncodingMetadata).With(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Bool("sub_context.found", true))
	return json.RawMessage(result), nil
}
