package service

import (
	"context"
	"univents/models"
)

func (o *Operations) GetCertByHash(ctx context.Context, hash string) (*models.Certification, error) {
	return o.Certs.GetByHash(ctx, hash)
}
