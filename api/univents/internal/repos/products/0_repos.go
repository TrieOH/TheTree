package products

import (
	"univents/internal/sqlc"
	"univents/models"
	"univents/ports"

	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.ProductRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("product"),
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
		GalleryURLs: src.GalleryUrls,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
		DeletedAt:   src.DeletedAt,
	}
}
