package mproduct

import (
	"context"

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

func (s *srv) Create(ctx context.Context, name string, price float64, categoryGUID uuid.UUID, description *string) (entity.ResponseProductCreate, error) {
	_, err := s.repoCategory.GetByGUID(ctx, categoryGUID)
	if err != nil {
		return entity.ResponseProductCreate{}, err
	}

	existing, err := s.repoProduct.GetByName(ctx, name)
	if err == nil && existing.GUID != uuid.Nil {
		return entity.ResponseProductCreate{}, entity.ErrProductAlreadyExists
	}

	product := entity.Product{
		Name:         name,
		Price:        price,
		CategoryGUID: categoryGUID,
		Description:  description,
	}

	created, err := s.repoProduct.Create(ctx, product)
	if err != nil {
		return entity.ResponseProductCreate{}, err
	}

	return entity.ResponseProductCreate{
		GUID:         created.GUID,
		Name:         created.Name,
		Price:        created.Price,
		CategoryGUID: created.CategoryGUID,
		Description:  created.Description,
	}, nil
}

func (s *srv) Get(ctx context.Context, guid uuid.UUID) (entity.ResponseProductGet, error) {
	product, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return entity.ResponseProductGet{}, err
	}

	return entity.ResponseProductGet{
		GUID:         product.GUID,
		Name:         product.Name,
		Price:        product.Price,
		CategoryGUID: product.CategoryGUID,
		Description:  product.Description,
	}, nil
}

func (s *srv) List(ctx context.Context, filter entity.RequestProductList) (entity.ResponseProductList, error) {
	if filter.CategoryGUID != nil {
		_, err := s.repoCategory.GetByGUID(ctx, *filter.CategoryGUID)
		if err != nil {
			return nil, err
		}
	}

	products, err := s.repoProduct.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	var response entity.ResponseProductList
	for _, p := range products {
		response = append(response, entity.ResponseProductGet{
			GUID:         p.GUID,
			Name:         p.Name,
			Price:        p.Price,
			CategoryGUID: p.CategoryGUID,
			Description:  p.Description,
		})
	}
	return response, nil
}

func (s *srv) Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.ResponseProductUpdate, error) {
	product, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return entity.ResponseProductUpdate{}, err
	}

	if req.Name != nil {
		existing, err := s.repoProduct.GetByName(ctx, *req.Name)
		if err == nil && existing.GUID != guid && existing.GUID != uuid.Nil {
			return entity.ResponseProductUpdate{}, entity.ErrProductAlreadyExists
		}
		product.Name = *req.Name
	}

	if req.Price != nil {
		product.Price = *req.Price
	}

	if req.CategoryGUID != nil {
		_, err := s.repoCategory.GetByGUID(ctx, *req.CategoryGUID)
		if err != nil {
			return entity.ResponseProductUpdate{}, err
		}
		product.CategoryGUID = *req.CategoryGUID
	}

	if req.Description != nil {
		product.Description = req.Description
	}

	updated, err := s.repoProduct.Update(ctx, product)
	if err != nil {
		return entity.ResponseProductUpdate{}, err
	}

	return entity.ResponseProductUpdate{
		GUID:         updated.GUID,
		Name:         updated.Name,
		Price:        updated.Price,
		CategoryGUID: updated.CategoryGUID,
		Description:  updated.Description,
	}, nil
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return err
	}

	return s.repoProduct.Delete(ctx, guid)
}
