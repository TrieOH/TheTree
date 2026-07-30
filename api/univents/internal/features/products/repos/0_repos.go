package repos

import (
	sqlc2 "univents/internal/sqlc"
	"univents/models"
	"univents/ports"

	"lib/database"
)

type Repo struct {
	q   *sqlc2.Queries
	dbe database.ErrorHandler
}

var _ ports.ProductRepo = (*Repo)(nil)

func NewRepo(q *sqlc2.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("product"),
	}
}

func mapProduct(src sqlc2.Product) models.Product {
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

func mapVariant(src sqlc2.ProductVariant) models.ProductVariant {
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
