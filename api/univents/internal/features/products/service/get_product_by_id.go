package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	return o.Products.GetProductByID(ctx, id)
}
