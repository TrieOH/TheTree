package queries

import (
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type QueryService struct {
	namespaces ports.NamespaceRepo
	forms      ports.FormsRepo
	tx         database.TxRunner
	tracer     trace.Tracer
}

func NewQueries(
	namespaces ports.NamespaceRepo,
	forms ports.FormsRepo,
	tx database.TxRunner,
	tracer trace.Tracer,
) *QueryService {
	return &QueryService{
		namespaces: namespaces,
		forms:      forms,
		tx:         tx,
		tracer:     tracer,
	}
}
