package repository

import (
	"context"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/google/uuid"
)

type Category interface {
	Create(ctx context.Context, category entity.Category) (entity.Category, error)
	GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Category, error)
	GetByName(ctx context.Context, name string) (entity.Category, error)
	List(ctx context.Context) ([]entity.Category, error)
	Update(ctx context.Context, category entity.Category) (entity.Category, error)
	Delete(ctx context.Context, guid uuid.UUID) error
}

type Product interface {
	Create(ctx context.Context, product entity.Product) (entity.Product, error)
	GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error)
	GetByName(ctx context.Context, name string) (entity.Product, error)
	List(ctx context.Context, filter entity.RequestProductList) ([]entity.Product, error)
	Update(ctx context.Context, product entity.Product) (entity.Product, error)
	Delete(ctx context.Context, guid uuid.UUID) error
}
