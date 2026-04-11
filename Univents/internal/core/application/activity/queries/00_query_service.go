package queries

import (
	"univents/internal/core/domain"
	"univents/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/authzed/authzed-go/v1"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	activities domain.ActivitiesRepository
	editions   domain.EditionsRepository
	gaClient   *goauth.Client
	tracer     trace.Tracer
	az         *authzed.Client
	tx         database.TxRunner
}

func New(
	activities domain.ActivitiesRepository,
	editions domain.EditionsRepository,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	az *authzed.Client,
	tx database.TxRunner,
) *QueryService {
	return &QueryService{
		activities: activities,
		editions:   editions,
		gaClient:   gaClient,
		tracer:     tracer,
		az:         az,
		tx:         tx,
	}
}
