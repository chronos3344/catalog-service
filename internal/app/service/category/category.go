package mcategory

import (
	"context"
	"errors"
	"time"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/repository"
	"github.com/chronos3344/catalog-service/internal/app/service"
	"github.com/google/uuid"
)

type srv struct {
	repoCategory repository.Category
	repoProduct  repository.Product
}

func NewService(repoCategory repository.Category, repoProduct repository.Product) service.Category {
	return &srv{
		repoCategory: repoCategory,
		repoProduct:  repoProduct,
	}
}

func (s *srv) Create(ctx context.Context, name string) (entity.Category, error) {
	var category entity.Category

	err := s.repoCategory.InsideTx(ctx, func(ctx context.Context) error {
		categories, err := s.repoCategory.List(ctx, &name)
		if err != nil {
			return err
		}

		if len(categories) > 0 {
			return entity.ErrAlreadyExists
		}

		category = entity.Category{
			GUID:      uuid.New(),
			Name:      name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return s.repoCategory.Create(ctx, category)
	})

	return category, err
}

func (s *srv) Get(ctx context.Context, guid uuid.UUID) (entity.Category, error) {
	return s.repoCategory.GetByGUID(ctx, guid)
}

func (s *srv) List(ctx context.Context) ([]entity.Category, error) {
	categories, err := s.repoCategory.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *srv) Update(ctx context.Context, guid uuid.UUID, name string) (entity.Category, error) {
	var category entity.Category
	var err error

	err = s.repoCategory.InsideTx(ctx, func(ctx context.Context) error {
		category, err = s.repoCategory.GetByGUID(ctx, guid)
		if err != nil {
			return err
		}

		categories, err := s.repoCategory.List(ctx, &name)
		if err != nil && !errors.Is(err, entity.ErrNotFound) {
			return err
		}

		if len(categories) > 0 && categories[0].GUID != guid {
			return entity.ErrAlreadyExists
		}

		category.Name = name
		category.UpdatedAt = time.Now()

		err = s.repoCategory.Update(ctx, category)
		if err != nil {
			return err
		}

		return nil
	})
	return category, err
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	err := s.repoCategory.InsideTx(ctx, func(ctx context.Context) error {
		_, err := s.repoCategory.GetByGUID(ctx, guid)
		if err != nil {
			return err
		}

		products, err := s.repoProduct.List(ctx, nil, &guid)
		if err != nil {
			return err
		}
		if len(products) > 0 {
			return entity.ErrCategoryHasProducts
		}

		return s.repoCategory.Delete(ctx, guid)
	})
	if err != nil {
		return err
	}
	return nil
}
