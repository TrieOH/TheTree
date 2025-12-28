package service

import (
	"GoAuth/internal/models"
	"GoAuth/internal/repo"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type AuthService struct {
	userRepo                 repo.UserRepo
	sessionRepo              repo.SessionRepo
	revokedRefreshTokensRepo repo.RevokedRefreshTokensRepo
	projectRepo              repo.ProjectRepo
	projectUserRepo          repo.ProjectUserRepo
}

func NewAuthService(
	userRepo repo.UserRepo,
	sessionRepo repo.SessionRepo,
	revokedRefreshTokensRepo repo.RevokedRefreshTokensRepo,
	projectRepo repo.ProjectRepo,
	projectUserRepo repo.ProjectUserRepo,
) *AuthService {
	return &AuthService{
		userRepo:                 userRepo,
		sessionRepo:              sessionRepo,
		revokedRefreshTokensRepo: revokedRefreshTokensRepo,
		projectRepo:              projectRepo,
		projectUserRepo:          projectUserRepo,
	}
}

var (
	GoAuthServiceTracer = otel.Tracer("goauth/service")
)

func annotateAccessClaims(span trace.Span, claims *models.AccessClaims) {
	span.SetAttributes(
		attribute.String("user.id", claims.Sub.ID.String()),
		attribute.String("user.session_id", claims.Sub.SessionID.String()),
		attribute.String("user.type", claims.Sub.UserType),
	)

	if claims.Sub.ProjectID != nil {
		span.SetAttributes(attribute.String("user.project_id", claims.Sub.ProjectID.String()))
	}
}
