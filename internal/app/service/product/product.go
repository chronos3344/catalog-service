package mproduct

import (
	"context"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/repository"
	"github.com/google/uuid"
)

type srv struct {
	repoProduct  repository.Product
	repoCategory repository.Category
}

func NewService(repoProduct *interface{}, repoCategory repository.Category) *srv {
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
	if err != nil {
		return entity.ResponseProductCreate{}, err
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

func (s *srv) List(ctx context.Context, categoryGUID *uuid.UUID, minPrice *float64, maxPrice *float64) (entity.ResponseProductList, error) {
	if categoryGUID != nil {
		_, err := s.repoCategory.GetByGUID(ctx, *categoryGUID)
		if err != nil {
			return nil, err
		}
	}

	product := entity.RequestProductList{
		CategoryGUID: categoryGUID,
		MinPrice:     minPrice,
		MaxPrice:     maxPrice,
	}

	products, err := s.repoProduct.List(ctx, product)
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

func (s *srv) Update(ctx context.Context, guid uuid.UUID, name *string, price *float64, categoryGUID *uuid.UUID, description *string) (entity.ResponseProductUpdate, error) {
	_, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return entity.ResponseProductUpdate{}, err
	}

	product := entity.ResponseProductUpdate{
		GUID: guid,
	}

	if name != nil {
		existing, err := s.repoProduct.GetByName(ctx, *name)
		if err == nil && existing.GUID != guid && existing.GUID != uuid.Nil {
			return entity.ResponseProductUpdate{}, entity.ErrProductAlreadyExists
		}
		if err != nil {
			return entity.ResponseProductUpdate{}, err
		}
		product.Name = *name
	}

	if price != nil {
		product.Price = *price
	}

	if categoryGUID != nil {
		_, err := s.repoCategory.GetByGUID(ctx, *categoryGUID)
		if err != nil {
			return entity.ResponseProductUpdate{}, err
		}
		product.CategoryGUID = *categoryGUID
	}

	if description != nil {
		product.Description = description
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
