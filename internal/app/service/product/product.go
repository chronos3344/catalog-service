package mproduct

import (
	"context"
	"errors"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/repository"
	"github.com/google/uuid"
)

type (
	Srv struct {
		*_service
	}
	_service struct {
		repoProduct  repository.Product
		repoCategory repository.Category
	}
)

func NewService(repoProduct repository.Product, repoCategory repository.Category) *Srv {
	return &Srv{
		_service: &_service{
			repoProduct:  repoProduct,
			repoCategory: repoCategory,
		},
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

	// Создаем продукт
	err = s.repoProduct.Create(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}

	// Возвращаем созданный продукт (с заполненными полями ID, GUID, CreatedAt, UpdatedAt)
	return product, nil
}

func (s *Srv) Get(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	_, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return entity.Product{}, err
	}

	return s.repoProduct.GetByGUID(ctx, guid)
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

	products, err := s.repoProduct.List(ctx, &product.Name, &product.CategoryGUID)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return entity.Product{}, err
	}

	for _, p := range products {
		if p.GUID != guid {
			return entity.Product{}, entity.ErrProductAlreadyExists
		}
	}

	err = s.repoProduct.Update(ctx, product)
	return product, err
}

func (s *Srv) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return err
	}

	return s.repoProduct.Delete(ctx, guid)
}
