package queries

import (
	"Informd/ports"
	"lib/database"

	v1 "github.com/authzed/authzed-go/v1"
	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	namespaces ports.NamespaceRepo
	az         *v1.Client
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewQueries(
	namespaces ports.NamespaceRepo,
	az *v1.Client,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		namespaces: namespaces,
		az:         az,
		tx:         tx,
		tracer:     tracer,
	}
}
