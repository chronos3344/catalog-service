package mproduct

import (
	"context"
	"time"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/repository"
	"github.com/chronos3344/catalog-service/internal/app/service"
	"github.com/google/uuid"
)

type srv struct {
	repoProduct  repository.Product
	repoCategory repository.Category
}

func NewService(repoProduct repository.Product, repoCategory repository.Category) service.Product {
	return &srv{
		repoProduct:  repoProduct,
		repoCategory: repoCategory,
	}
}

func (s *srv) Create(ctx context.Context, req entity.RequestProductCreate) (entity.Product, error) {
	_, err := s.repoCategory.GetByGUID(ctx, req.CategoryGUID)
	if err != nil {
		return entity.Product{}, err
	}

	existing, err := s.repoProduct.List(ctx, &req.Name, &req.CategoryGUID)
	if err != nil {
		return entity.Product{}, err
	}
	if len(existing) > 0 {
		return entity.Product{}, entity.ErrAlreadyExists
	}

	now := time.Now()
	product := entity.Product{
		GUID:         uuid.New(),
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		CategoryGUID: req.CategoryGUID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repoProduct.Create(ctx, product); err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (s *srv) Get(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	product, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (s *srv) List(ctx context.Context) ([]entity.Product, error) {
	return s.repoProduct.List(ctx, nil, nil)
}

func (s *srv) Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.Product, error) {
	product, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return entity.Product{}, err
	}

	if req.CategoryGUID != nil && *req.CategoryGUID != product.CategoryGUID {
		_, err := s.repoCategory.GetByGUID(ctx, *req.CategoryGUID)
		if err != nil {
			return entity.Product{}, err
		}
		product.CategoryGUID = *req.CategoryGUID
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Description != nil {
		product.Description = req.Description
	}

	existing, err := s.repoProduct.List(ctx, &product.Name, &product.CategoryGUID)
	if err != nil {
		return entity.Product{}, err
	}
	for _, p := range existing {
		if p.GUID != guid {
			return entity.Product{}, entity.ErrAlreadyExists
		}
	}

	product.UpdatedAt = time.Now()

	if err := s.repoProduct.Update(ctx, product); err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return err
	}

	return s.repoProduct.Delete(ctx, guid)
}
