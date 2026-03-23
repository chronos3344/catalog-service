package mcategory

import (
	"context"
	"errors"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/repository"
	"github.com/chronos3344/catalog-service/internal/app/service"
	"github.com/google/uuid"
)

type srv struct {
	repoCategory repository.Category
}

func NewService(repoCategory repository.Category) service.Category {
	return &srv{
		repoCategory: repoCategory,
	}
}

func (s *srv) Create(ctx context.Context, name string) (entity.Category, error) {
	// Проверяем существование категории с таким именем
	categories, err := s.repoCategory.List(ctx, &name)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return entity.Category{}, err
	}

	for _, cat := range categories {
		if cat.Name == name {
			return entity.Category{}, entity.ErrCategoryAlreadyExists
		}
	}

	category := entity.Category{
		Name: name,
	}

	err = s.repoCategory.Create(ctx, category)
	if err != nil {
		return entity.Category{}, err
	}

	return category, nil
}

func (s *srv) Get(ctx context.Context, guid uuid.UUID) (entity.Category, error) {
	category, err := s.repoCategory.GetByGUID(ctx, guid)
	if err != nil {
		return entity.Category{}, err
	}

	return category, nil
}

func (s *srv) List(ctx context.Context) ([]entity.Category, error) {
	categories, err := s.repoCategory.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *srv) Update(ctx context.Context, guid uuid.UUID, name string) (entity.Category, error) {
	// Получаем существующую категорию
	category, err := s.repoCategory.GetByGUID(ctx, guid)
	if err != nil {
		return category, err
	}

	// Проверяем уникальность нового имени
	categories, err := s.repoCategory.List(ctx, &name)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return entity.Category{}, err
	}

	for _, cat := range categories {
		if cat.Name == name && cat.GUID != guid {
			return entity.Category{}, entity.ErrCategoryAlreadyExists
		}
	}

	// Обновляем имя
	category.Name = name

	// Сохраняем изменения
	err = s.repoCategory.Update(ctx, category)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	return s.repoCategory.Delete(ctx, guid)
}
