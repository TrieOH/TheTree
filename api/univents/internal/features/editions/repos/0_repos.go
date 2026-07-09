package repos

import (
	"univents/internal/database/sqlc"

	"lib/database"
	"univents/contracts"
	"univents/internal/shared/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type editionsRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.EditionsRepository = (*editionsRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.EditionsRepository {
	return &editionsRepo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("activities"),
	}
}

func mapEditionFromDB(src *sqlc.Edition) *contracts.Edition {
	return &contracts.Edition{
		ID:                            src.ID,
		GoauthScopeID:                 src.GoauthScopeID,
		EventID:                       src.EventID,
		Type:                          contracts.EditionType(src.Type),
		EditionName:                   src.EditionName,
		Tagline:                       src.Tagline,
		Description:                   src.Description,
		Status:                        contracts.EditionStatus(src.Status),
		MonetaryType:                  contracts.EditionMonetaryType(src.MonetaryType),
		RegistrationOpensAt:           src.RegistrationOpensAt,
		RegistrationClosesAt:          src.RegistrationClosesAt,
		StartsAt:                      src.StartsAt,
		EndsAt:                        src.EndsAt,
		Timezone:                      src.Timezone,
		LocationName:                  src.LocationName,
		LocationAddress:               src.LocationAddress,
		LogoUrl:                       src.LogoUrl,
		BannerUrl:                     src.BannerUrl,
		ContactEmail:                  src.ContactEmail,
		ContactPhone:                  src.ContactPhone,
		OrganizerName:                 src.OrganizerName,
		TriePaymentsCredentialID:      src.TriePaymentsCredentialID,
		TriePaymentsProvider:          src.TriePaymentsProvider,
		TriePaymentsProviderPublicKey: src.TriePaymentsProviderPublicKey,
		CreatedBy:                     src.CreatedBy,
		CreatedAt:                     src.CreatedAt,
		UpdatedAt:                     src.UpdatedAt,
		DeletedAt:                     src.DeletedAt,
	}
}
