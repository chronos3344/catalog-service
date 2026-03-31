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
	// Проверяем существование категории с таким именем
	categories, err := s.repoCategory.List(ctx, &name)
	if err != nil {
		return entity.Category{}, err
	}

	if len(categories) > 0 {
		return entity.Category{}, entity.ErrAlreadyExists
	}

	category := entity.Category{
		GUID:      uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repoCategory.Create(ctx, category)
	if err != nil {
		return entity.Category{}, err
	}

	return category, nil
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
	// Получаем существующую категорию
	category, err := s.repoCategory.GetByGUID(ctx, guid)
	if err != nil {
		return entity.Category{}, err
	}

	// Проверяем уникальность нового имени
	categories, err := s.repoCategory.List(ctx, &name)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return entity.Category{}, err
	}

	if len(categories) > 0 && categories[0].GUID != guid {
		return entity.Category{}, entity.ErrAlreadyExists
	}

	// Обновляем имя
	category.Name = name
	category.UpdatedAt = time.Now()

	// Сохраняем изменения
	err = s.repoCategory.Update(ctx, category)
	if err != nil {
		return entity.Category{}, err
	}

	return category, nil
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
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
}
