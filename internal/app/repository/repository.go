package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/chronos3344/catalog-service/internal/app/entity"

)

type (
	Category interface {
		Create(ctx context.Context, req entity.RequestCategoryCreate) (entity.ResponseCategory, error)
		GetByGUID(ctx context.Context, guid uuid.UUID) (entity.ResponseCategory, error)
		Update(ctx context.Context, guid uuid.UUID, req entity.RequestCategoryUpdate) (entity.ResponseCategory, error)
		Delete(ctx context.Context, guid uuid.UUID) error
		List(ctx context.Context, page, limit int) (entity.ResponseCategoryList, error)
	}

	Product interface {
		Create(ctx context.Context, req entity.RequestProductCreate) (entity.ResponseProduct, error)
		GetByGUID(ctx context.Context, guid uuid.UUID) (entity.ResponseProduct, error)
		Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.ResponseProduct, error)
		Delete(ctx context.Context, guid uuid.UUID) error
		List(ctx context.Context, page, limit int, categoryGUID *uuid.UUID) (entity.ResponseProductList, error)
	}
)