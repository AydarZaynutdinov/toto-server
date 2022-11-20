package repository

import (
	"context"
	"toto-server/internal/entity"
)

type ISkuConfig interface {
	Get(ctx context.Context, params QueryParameters) (entity.SkuConfig, error)
}
