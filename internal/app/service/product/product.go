package mproduct

import (
	"context"
	"errors"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/repository"
	"github.com/chronos3344/catalog-service/internal/app/service"
	"github.com/google/uuid"
)

type Srv struct {
	repoProduct  repository.Product
	repoCategory repository.Category
}

func NewService(repoProduct repository.Product, repoCategory repository.Category) service.Product {
	return &Srv{
		repoProduct:  repoProduct,
		repoCategory: repoCategory,
	}
}

func (s *Srv) Create(ctx context.Context, product entity.Product) (entity.Product, error) {
	_, err := s.repoCategory.GetByGUID(ctx, product.CategoryGUID)
	if err != nil {
		return entity.Product{}, err
	}

	existingList, err := s.repoProduct.List(ctx, &product.Name, &product.CategoryGUID)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return entity.Product{}, err
	}
	if len(existingList) > 0 {
		return entity.Product{}, entity.ErrProductAlreadyExists
	}

	err = s.repoProduct.Create(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (s *Srv) Get(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	product, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (s *Srv) List(ctx context.Context, filter entity.RequestProductList) ([]entity.Product, error) {
	if filter.CategoryGUID != nil {
		_, err := s.repoCategory.GetByGUID(ctx, *filter.CategoryGUID)
		if err != nil {
			return nil, err
		}
	}

	products, err := s.repoProduct.List(ctx, filter.Name, filter.CategoryGUID)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *Srv) Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.Product, error) {
	product, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return entity.Product{}, err
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.CategoryGUID != nil {
		product.CategoryGUID = *req.CategoryGUID
	}

	if req.Price != nil {
		product.Price = *req.Price
	}

	if req.Description != nil {
		product.Description = req.Description
	}

	products, err := s.repoProduct.List(ctx, &product.Name, &product.CategoryGUID)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return entity.Product{}, err
	}

	if len(products) > 0 && products[0].GUID != guid {
		return entity.Product{}, entity.ErrProductAlreadyExists
	}

	err = s.repoProduct.Update(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}
	return product, nil
}

func (s *Srv) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return err
	}

	return s.repoProduct.Delete(ctx, guid)
}
