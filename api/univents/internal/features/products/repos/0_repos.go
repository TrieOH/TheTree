package repos

import (
	"univents/internal/database/sqlc"
	"univents/models"
	"univents/ports"

	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.ProductRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ProductRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("product"),
	}
}

func mapProduct(src sqlc.Product) models.Product {
	return models.Product{
		ID:                   src.ID,
		EditionID:            src.EditionID,
		VendorCode:           src.VendorCode,
		RequiresRegistration: src.RequiresRegistration,
		CreatedAt:            src.CreatedAt,
		UpdatedAt:            src.UpdatedAt,
		DeletedAt:            src.DeletedAt,
	}
}

func mapVariant(src sqlc.ProductVariant) models.ProductVariant {
	return models.ProductVariant{
		ID:          src.ID,
		EditionID:   src.EditionID,
		ProductID:   src.ProductID,
		VendorCode:  src.VendorCode,
		Name:        src.Name,
		Description: src.Description,
		Price:       src.Price,
		Stock:       src.Stock,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
		DeletedAt:   src.DeletedAt,
	}
}
