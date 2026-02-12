package service

import (
	"context"
	
	"github.com/google/uuid"

	"github.com/chronos3344/catalog-service/internal/app/entity"
)

type Category interface {
	Create(ctx context.Context, name string) (entity.ResponseCategoryCreate, error)
	Get(ctx context.Context, guid uuid.UUID) (entity.ResponseCategoryGet, error)
	List(ctx context.Context) (entity.ResponseCategoryList, error)
	Update(ctx context.Context, guid uuid.UUID, name string) (entity.ResponseCategoryUpdate, error)
	Delete(ctx context.Context, guid uuid.UUID) error
}

type Product interface {
	Create(ctx context.Context, name string, price float64, categoryGUID uuid.UUID, description *string) (entity.ResponseProductCreate, error)
	Get(ctx context.Context, guid uuid.UUID) (entity.ResponseProductGet, error)
	List(ctx context.Context, filter entity.RequestProductList) (entity.ResponseProductList, error)
	Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.ResponseProductUpdate, error)
	Delete(ctx context.Context, guid uuid.UUID) error
}
