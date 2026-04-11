package commands

import (
	"univents/internal/commerce/domain"
	coreDomain "univents/internal/core/domain"
	"univents/internal/plataform/database"
	"univents/internal/shared/sockets"

	paymentsSDK "github.com/TrieOH/TriePaymentsSDK"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/authzed/authzed-go/v1"
	"github.com/hibiken/asynq"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/trace"
)

type CommandService struct {
	editions  coreDomain.EditionsRepository
	products  domain.ProductsRepository
	purchases domain.PurchaseRepository
	payments  *paymentsSDK.Client
	sessions  domain.PurchaseSessionStore
	ws        *sockets.Registry
	inventory domain.InventoryPublisher
	minio     *minio.Client
	asynq     *asynq.Client
	inspector *asynq.Inspector
	gaClient  *goauth.Client
	tracer    trace.Tracer
	az        *authzed.Client
	tx        database.TxRunner
}

func New(
	editions coreDomain.EditionsRepository,
	products domain.ProductsRepository,
	purchases domain.PurchaseRepository,
	payments *paymentsSDK.Client,
	session domain.PurchaseSessionStore,
	ws *sockets.Registry,
	inventory domain.InventoryPublisher,
	minio *minio.Client,
	asynq *asynq.Client,
	inspector *asynq.Inspector,
	gaClient *goauth.Client,
	tracer trace.Tracer,
	az *authzed.Client,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		editions:  editions,
		products:  products,
		purchases: purchases,
		payments:  payments,
		sessions:  session,
		ws:        ws,
		inventory: inventory,
		minio:     minio,
		asynq:     asynq,
		inspector: inspector,
		gaClient:  gaClient,
		tracer:    tracer,
		az:        az,
		tx:        tx,
	}
}
