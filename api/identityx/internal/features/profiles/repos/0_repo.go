package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type profileRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.ProfileRepo = (*profileRepo)(nil)

func NewProfileRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ProfileRepo {
	return &profileRepo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("actor_profile"),
	}
}

func mapActorProfile(src sqlc.ActorProfile) models.ActorProfile {
	return models.ActorProfile{
		ActorID:   src.ActorID,
		Profile:   src.Profile,
		UpdatedAt: src.UpdatedAt,
	}
}
