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
	Create(ctx context.Context, name string, price float64, categoryGUID uuid.UUID, description *string) (entity.Product, error)
	GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error)
	GetByName(ctx context.Context, name string) (entity.Product, error)
	List(ctx context.Context, categoryGUID *uuid.UUID, minPrice *float64, maxPrice *float64) ([]entity.Product, error)
	Update(ctx context.Context, guid uuid.UUID, name *string, price *float64, categoryGUID *uuid.UUID, description *string) (entity.ResponseProductUpdate, error)
	Delete(ctx context.Context, guid uuid.UUID) error
}
