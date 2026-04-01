package service

import (
	"context"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/google/uuid"
)

type Category interface {
	Create(ctx context.Context, name string) (entity.Category, error)
	Get(ctx context.Context, guid uuid.UUID) (entity.Category, error)
	List(ctx context.Context) ([]entity.Category, error)
	Update(ctx context.Context, guid uuid.UUID, name string) (entity.Category, error)
	Delete(ctx context.Context, guid uuid.UUID) error
}

type Product interface {
	Create(ctx context.Context, product entity.RequestProductCreate) (entity.Product, error)
	Get(ctx context.Context, guid uuid.UUID) (entity.Product, error)
	List(ctx context.Context) ([]entity.Product, error)
	Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.Product, error)
	Delete(ctx context.Context, guid uuid.UUID) error
}
