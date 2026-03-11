package queries

import (
	"TriePayments/internal/core/domain"
	"TriePayments/internal/plataform/database"

	"github.com/TrieOH/goauth-sdk-go"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	workspaces   domain.WorkspaceRepo
	marketplaces domain.MarketplaceConfigRepo
	gaClient     *goauth.Client
	tx           database.TxRunner
	tracer       trace.Tracer
}

func New(
	workspaces domain.WorkspaceRepo,
	marketplaces domain.MarketplaceConfigRepo,
	gaClient *goauth.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		workspaces:   workspaces,
		marketplaces: marketplaces,
		gaClient:     gaClient,
		tx:           tx,
		tracer:       tracer,
	}
}
