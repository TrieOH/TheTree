package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID                   uuid.UUID  `json:"id"`
	EditionID            uuid.UUID  `json:"edition_id"`
	VendorCode           string     `json:"vendor_code"`
	RequiresRegistration bool       `json:"requires_registration"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at"`
}

type ProductVariant struct {
	ID          uuid.UUID  `json:"id"`
	EditionID   uuid.UUID  `json:"edition_id"`
	ProductID   uuid.UUID  `json:"product_id"`
	VendorCode  string     `json:"vendor_code"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Price       int64      `json:"price"`
	Stock       *int       `json:"stock"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// ── Create Initial Product (product + first variant) ──────────────────────

type CreateInitialProductRequest struct {
	RequiresRegistration bool    `json:"requires_registration"`
	VendorCode           string  `json:"vendor_code"           validate:"required,min=2,max=255"`
	VariantVendorCode    string  `json:"variant_vendor_code"   validate:"required,min=2,max=255"`
	Name                 string  `json:"name"                  validate:"required,min=2"`
	Description          *string `json:"description"`
	Price                int64   `json:"price"                 validate:"gte=0"`
	Stock                *int    `json:"stock"                 validate:"omitempty,gte=0"`
}

func (r CreateInitialProductRequest) ToInput(editionID uuid.UUID) CreateInitialProductInput {
	return CreateInitialProductInput{
		EditionID:            editionID,
		RequiresRegistration: r.RequiresRegistration,
		VendorCode:           r.VendorCode,
		VariantVendorCode:    r.VariantVendorCode,
		Name:                 r.Name,
		Description:          r.Description,
		Price:                r.Price,
		Stock:                r.Stock,
	}
}

type CreateInitialProductInput struct {
	EditionID            uuid.UUID `json:"edition_id"`
	VendorCode           string    `json:"vendor_code"`
	VariantVendorCode    string    `json:"variant_vendor_code"`
	RequiresRegistration bool      `json:"requires_registration"`
	Name                 string    `json:"name"`
	Description          *string   `json:"description"`
	Price                int64     `json:"price"`
	Stock                *int      `json:"stock"`
}

// ── Create Product Variant ────────────────────────────────────────────────

type CreateProductVariantRequest struct {
	VendorCode  string  `json:"vendor_code" validate:"required,min=2,max=255"`
	Name        string  `json:"name"        validate:"required,min=2"`
	Description *string `json:"description"`
	Price       int64   `json:"price"       validate:"gte=0"`
	Stock       *int    `json:"stock"       validate:"omitempty,gte=0"`
}

func (r CreateProductVariantRequest) ToInput(productID uuid.UUID) CreateProductVariantInput {
	return CreateProductVariantInput{
		ProductID:   productID,
		VendorCode:  r.VendorCode,
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
		Stock:       r.Stock,
	}
}

type CreateProductVariantInput struct {
	ProductID   uuid.UUID `json:"product_id"`
	VendorCode  string    `json:"vendor_code"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Price       int64     `json:"price"`
	Stock       *int      `json:"stock"`
}

// ── Patch Product ─────────────────────────────────────────────────────────

type PatchProductRequest struct {
	VendorCode           string `json:"vendor_code"           validate:"required,min=2,max=255"`
	RequiresRegistration bool   `json:"requires_registration"`
}

func (r PatchProductRequest) ToInput(productID uuid.UUID) PatchProductInput {
	return PatchProductInput{
		ProductID:            productID,
		VendorCode:           r.VendorCode,
		RequiresRegistration: r.RequiresRegistration,
	}
}

type PatchProductInput struct {
	ProductID            uuid.UUID
	VendorCode           string
	RequiresRegistration bool
}

// ── Patch Product Variant ─────────────────────────────────────────────────

type PatchProductVariantRequest struct {
	VendorCode  string  `json:"vendor_code" validate:"required,min=2,max=255"`
	Name        string  `json:"name"        validate:"required,min=2"`
	Description *string `json:"description"`
	Price       int64   `json:"price"       validate:"gte=0"`
	Stock       *int    `json:"stock"       validate:"omitempty,gte=0"`
}

func (r PatchProductVariantRequest) ToInput(variantID uuid.UUID) PatchProductVariantInput {
	return PatchProductVariantInput{
		VariantID:   variantID,
		VendorCode:  r.VendorCode,
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
		Stock:       r.Stock,
	}
}

type PatchProductVariantInput struct {
	VariantID   uuid.UUID
	VendorCode  string
	Name        string
	Description *string
	Price       int64
	Stock       *int
}
