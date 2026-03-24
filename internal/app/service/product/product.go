package mproduct

import (
	"context"
	"errors"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/repository"
	"github.com/google/uuid"
)

type srv struct {
	repoProduct  repository.Product
	repoCategory repository.Category
}

func NewService(repoProduct repository.Product, repoCategory repository.Category) *srv {
	return &srv{
		repoProduct:  repoProduct,
		repoCategory: repoCategory,
	}
}

func (s *srv) Create(ctx context.Context, product entity.Product) (entity.Product, error) {
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

func (s *srv) Get(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	product, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return entity.Product{}, err
	}

	return entity.Product{
		GUID:         product.GUID,
		Name:         product.Name,
		Price:        product.Price,
		CategoryGUID: product.CategoryGUID,
		Description:  product.Description,
	}, nil
}

func (s *srv) List(ctx context.Context, filter entity.RequestProductList) ([]entity.Product, error) {
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

func (s *srv) Update(ctx context.Context, guid uuid.UUID, req entity.RequestProductUpdate) (entity.Product, error) {
	// Шаг 1: Получаем существующий продукт
	existingProduct, err := s.getProductForUpdate(ctx, guid)
	if err != nil {
		return entity.Product{}, err
	}

	// Шаг 2: Проверяем уникальность имени
	if err := s.checkNameUniqueness(ctx, guid, existingProduct.Name, req.Name, existingProduct.CategoryGUID, req.CategoryGUID); err != nil {
		return entity.Product{}, err
	}

	// Шаг 3: Проверяем существование категории
	if err := s.checkCategoryExists(ctx, existingProduct.CategoryGUID, req.CategoryGUID); err != nil {
		return entity.Product{}, err
	}

	// Шаг 4: Применяем обновления
	productToUpdate := s.buildProductForUpdate(guid, existingProduct, req)

	// Шаг 5: Сохраняем
	err = s.repoProduct.Update(ctx, productToUpdate)
	if err != nil {
		return entity.Product{}, err
	}

	// Возвращаем обновленный продукт
	return productToUpdate, nil
}

// getProductForUpdate получает продукт для обновления
func (s *srv) getProductForUpdate(ctx context.Context, guid uuid.UUID) (*entity.Product, error) {
	product, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

// checkNameUniqueness проверяет уникальность имени при изменении
func (s *srv) checkNameUniqueness(ctx context.Context, guid uuid.UUID, oldName string, newName *string, oldCategoryGUID uuid.UUID, newCategoryGUID *uuid.UUID) error {
	if newName == nil || *newName == oldName {
		return nil
	}

	// Определяем категорию, в которой проверяем уникальность
	categoryGUID := oldCategoryGUID
	if newCategoryGUID != nil {
		categoryGUID = *newCategoryGUID
	}

	existingList, err := s.repoProduct.List(ctx, newName, &categoryGUID)
	if err != nil {
		return err
	}

	// Если нашли продукт с таким же именем в той же категории, и это не текущий продукт — ошибка
	for _, p := range existingList {
		if p.GUID != guid {
			return entity.ErrProductAlreadyExists
		}
	}
	return nil
}

// checkCategoryExists проверяет существование категории при изменении
func (s *srv) checkCategoryExists(ctx context.Context, oldCategoryGUID uuid.UUID, newCategoryGUID *uuid.UUID) error {
	if newCategoryGUID == nil || *newCategoryGUID == oldCategoryGUID {
		return nil
	}

	_, err := s.repoCategory.GetByGUID(ctx, *newCategoryGUID)
	return err
}

// buildProductForUpdate создает объект продукта для обновления
func (s *srv) buildProductForUpdate(guid uuid.UUID, existing *entity.Product, req entity.RequestProductUpdate) entity.Product {
	product := entity.Product{
		GUID:         guid,
		Name:         existing.Name,
		Price:        existing.Price,
		CategoryGUID: existing.CategoryGUID,
		Description:  existing.Description,
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.CategoryGUID != nil {
		product.CategoryGUID = *req.CategoryGUID
	}
	if req.Description != nil {
		product.Description = req.Description
	}
	return product
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return err
	}

	return s.repoProduct.Delete(ctx, guid)
}
