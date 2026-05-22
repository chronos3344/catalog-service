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
	var product entity.Product

	err := s.repoProduct.InsideTx(ctx, func(ctx context.Context) error {
		_, err := s.repoCategory.GetByGUID(ctx, req.CategoryGUID)
		if err != nil {
			return err
		}

		existing, err := s.repoProduct.List(ctx, &req.Name, &req.CategoryGUID)
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			return entity.ErrAlreadyExists
		}

		now := time.Now()
		product = entity.Product{
			GUID:         uuid.New(),
			Name:         req.Name,
			Description:  req.Description,
			Price:        req.Price,
			CategoryGUID: req.CategoryGUID,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		return s.repoProduct.Create(ctx, product)
	})
	return product, err
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
	var product entity.Product
	var err error

	err = s.repoProduct.InsideTx(ctx, func(ctx context.Context) error {
		product, err = s.repoProduct.GetByGUID(ctx, guid)
		if err != nil {
			return err
		}

		if req.CategoryGUID != nil && *req.CategoryGUID != product.CategoryGUID {
			_, err := s.repoCategory.GetByGUID(ctx, *req.CategoryGUID)
			if err != nil {
				return err
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
			return err
		}
		for _, p := range existing {
			if p.GUID != guid {
				return entity.ErrAlreadyExists
			}
		}

		product.UpdatedAt = time.Now()

		return s.repoProduct.Update(ctx, product)
	})
	return product, err
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	err := s.repoProduct.InsideTx(ctx, func(ctx context.Context) error {
		_, err := s.repoProduct.GetByGUID(ctx, guid)
		if err != nil {
			return err
		}
		return s.repoProduct.Delete(ctx, guid)
	})
	if err != nil {
		return err
	}
	return nil
}
