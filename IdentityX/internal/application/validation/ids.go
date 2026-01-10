package validation

import (
	"GoAuth/internal/apierr"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

func ParseProjectID(span trace.Span, projectID string) (*uuid.UUID, error) {
	pid, err := uuid.Parse(projectID)
	if err != nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("invalid project id").WithID(apierr.ProjectInvalidID).WithCause(err)
		apierr.RecordDomainError(span, apiErr)
		return nil, apiErr
	}
	return &pid, nil
}

func ProjectIDNotNull(span trace.Span, projectID *string) error {
	if projectID == nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("project id is required").WithID(apierr.ProjectMissingID)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}
	return nil
}

func RequireProjectID(span trace.Span, projectID *string) (*uuid.UUID, error) {
	err := ProjectIDNotNull(span, projectID)
	if err != nil {
		return nil, err
	}
	return ParseProjectID(span, *projectID)
}

func ParseSchemaID(span trace.Span, schemaID string) (*uuid.UUID, error) {
	sid, err := uuid.Parse(schemaID)
	if err != nil {
		err = apierr.ErrInvalidInput.WithMsg("invalid schema id").WithID(apierr.SchemaInvalidID).WithCause(err)
		apierr.RecordDomainError(span, err)
		return nil, err
	}
	return &sid, nil
}

func SchemaIDNotNull(span trace.Span, schemaID *string) error {
	if schemaID == nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("schema id is required").WithID(apierr.SchemaMissingID)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}
	return nil
}

func RequireSchemaID(span trace.Span, schemaID *string) (*uuid.UUID, error) {
	err := SchemaIDNotNull(span, schemaID)
	if err != nil {
		return nil, err
	}
	return ParseSchemaID(span, *schemaID)
}

func ParseSessionID(span trace.Span, sessionID string) (*uuid.UUID, error) {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		err = apierr.ErrInvalidInput.WithMsg("invalid session id").WithID(apierr.SessionInvalidID).WithCause(err)
		apierr.RecordDomainError(span, err)
		return nil, err
	}
	return &sid, nil
}

func SessionIDNotNull(span trace.Span, sessionID *string) error {
	if sessionID == nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("session id is required").WithID(apierr.SessionMissingID)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}
	return nil
}

func RequireSessionID(span trace.Span, sessionID *string) (*uuid.UUID, error) {
	err := SessionIDNotNull(span, sessionID)
	if err != nil {
		return nil, err
	}
	return ParseSessionID(span, *sessionID)
}

func ParseRefreshJTI(span trace.Span, refreshJTI string) (*uuid.UUID, error) {
	jti, err := uuid.Parse(refreshJTI)
	if err != nil {
		tokenErr := apierr.ErrInvalidInput.WithMsg("invalid refresh token id").WithID(apierr.TokenRefreshInvalidID)
		apierr.RecordDomainError(span, tokenErr)
		return nil, tokenErr
	}
	return &jti, nil
}

func RefreshJTINotNull(span trace.Span, refreshJTI *string) error {
	if refreshJTI == nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("refresh token id is required").WithID(apierr.TokenRefreshIDMissing)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}
	return nil
}

func RequireRefreshJTI(span trace.Span, refreshJTI *string) (*uuid.UUID, error) {
	err := RefreshJTINotNull(span, refreshJTI)
	if err != nil {
		return nil, err
	}
	return ParseRefreshJTI(span, *refreshJTI)
}

func ParseAccessJTI(span trace.Span, refreshJTI string) (*uuid.UUID, error) {
	jti, err := uuid.Parse(refreshJTI)
	if err != nil {
		tokenErr := apierr.ErrInvalidInput.WithMsg("invalid access token id").WithID(apierr.TokenAccessInvalidID)
		apierr.RecordDomainError(span, tokenErr)
		return nil, tokenErr
	}
	return &jti, nil
}

func AccessJTINotNull(span trace.Span, refreshJTI *string) error {
	if refreshJTI == nil {
		apiErr := apierr.ErrInvalidInput.WithMsg("refresh token id is required").WithID(apierr.TokenAccessIDMissing)
		apierr.RecordDomainError(span, apiErr)
		return apiErr
	}
	return nil
}

func RequireAccessJTI(span trace.Span, refreshJTI *string) (*uuid.UUID, error) {
	err := AccessJTINotNull(span, refreshJTI)
	if err != nil {
		return nil, err
	}
	return ParseAccessJTI(span, *refreshJTI)
}
