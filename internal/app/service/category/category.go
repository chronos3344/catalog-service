package mcategory

import (
	"context"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/repository"
	"github.com/google/uuid"
)

type service struct {
	repoCategory repository.Category
}

func NewService(repoCategory repository.Category) service.Category {
	return &service{
		repoCategory: repoCategory,
	}
}

func (s *service) Create(ctx context.Context, name string) (entity.ResponseCategoryCreate, error) {
	existing, err := s.repoCategory.GetByName(ctx, name)
	if err == nil && existing.GUID != uuid.Nil {
		return entity.ResponseCategoryCreate{}, repository.ErrConflict
	}

	category := entity.Category{
		Name: name,
	}

	created, err := s.repoCategory.Create(ctx, category)
	if err != nil {
		return entity.ResponseCategoryCreate{}, err
	}

	return entity.ResponseCategoryCreate{
		GUID: created.GUID,
		Name: created.Name,
	}, nil
}

func (s *service) Get(ctx context.Context, guid uuid.UUID) (entity.ResponseCategoryGet, error) {
	category, err := s.repoCategory.GetByGUID(ctx, guid)
	if err != nil {
		return entity.ResponseCategoryGet{}, err
	}

	return entity.ResponseCategoryGet{
		GUID: category.GUID,
		Name: category.Name,
	}, nil
}

func (s *service) List(ctx context.Context) (entity.ResponseCategoryList, error) {
	categories, err := s.repoCategory.List(ctx)
	if err != nil {
		return nil, err
	}

	var response entity.ResponseCategoryList
	for _, cat := range categories {
		response = append(response, entity.ResponseCategoryGet{
			GUID: cat.GUID,
			Name: cat.Name,
		})
	}
	return response, nil
}

func (s *service) Update(ctx context.Context, guid uuid.UUID, name string) (entity.ResponseCategoryUpdate, error) {
	category, err := s.repoCategory.GetByGUID(ctx, guid)
	if err != nil {
		return entity.ResponseCategoryUpdate{}, err
	}

	existing, err := s.repoCategory.GetByName(ctx, name)
	if err == nil && existing.GUID != guid && existing.GUID != uuid.Nil {
		return entity.ResponseCategoryUpdate{}, repository.ErrConflict
	}

	category.Name = name

	updated, err := s.repoCategory.Update(ctx, category)
	if err != nil {
		return entity.ResponseCategoryUpdate{}, err
	}

	return entity.ResponseCategoryUpdate{
		GUID: updated.GUID,
		Name: updated.Name,
	}, nil
}

func (s *service) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := s.repoCategory.GetByGUID(ctx, guid)
	if err != nil {
		return err
	}

	return s.repoCategory.Delete(ctx, guid)
}
