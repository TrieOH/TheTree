package queries

import (
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type Queries struct {
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
) *Queries {
	return &Queries{
		namespaces: namespaces,
		forms:      forms,
		tx:         tx,
		tracer:     tracer,
	}
}
