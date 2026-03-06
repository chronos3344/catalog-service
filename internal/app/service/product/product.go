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

func (s *srv) Create(ctx context.Context, product entity.Product) (entity.ResponseProductCreate, error) {
	_, err := s.repoCategory.GetByGUID(ctx, product.CategoryGUID)
	if err != nil {
		return entity.ResponseProductCreate{}, err
	}

	existing, err := s.repoProduct.GetByNameAndCategory(ctx, product.Name, product.CategoryGUID)
	if err == nil && existing.GUID != uuid.Nil {
		return entity.ResponseProductCreate{}, entity.ErrProductAlreadyExists
	}
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return entity.ResponseProductCreate{}, err
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
	// Шаг 1: Получаем существующий продукт
	existingProduct, err := s.getProductForUpdate(ctx, guid)
	if err != nil {
		return entity.ResponseProductUpdate{}, err
	}

	// Шаг 2: Проверяем уникальность имени
	if err := s.checkNameUniqueness(ctx, guid, existingProduct.Name, req.Name); err != nil {
		return entity.ResponseProductUpdate{}, err
	}

	// Шаг 3: Проверяем существование категории
	if err := s.checkCategoryExists(ctx, existingProduct.CategoryGUID, req.CategoryGUID); err != nil {
		return entity.ResponseProductUpdate{}, err
	}

	// Шаг 4: Применяем обновления
	productToUpdate := s.buildProductForUpdate(guid, existingProduct, req)

	// Шаг 5: Сохраняем
	updated, err := s.repoProduct.Update(ctx, productToUpdate)
	if err != nil {
		return entity.ResponseProductUpdate{}, err
	}

	// Шаг 6: Формируем ответ
	return s.buildUpdateResponse(updated), nil
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
func (s *srv) checkNameUniqueness(ctx context.Context, guid uuid.UUID, oldName string, newName *string) error {
	if newName == nil || *newName == oldName {
		return nil
	}

	existing, err := s.repoProduct.GetByName(ctx, *newName)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return nil // имя уникально
		}
		return err
	}

	if existing.GUID != guid && existing.GUID != uuid.Nil {
		return entity.ErrProductAlreadyExists
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

// buildUpdateResponse формирует ответ
func (s *srv) buildUpdateResponse(updated entity.Product) entity.ResponseProductUpdate {
	return entity.ResponseProductUpdate{
		GUID:         updated.GUID,
		Name:         updated.Name,
		Price:        updated.Price,
		CategoryGUID: updated.CategoryGUID,
		Description:  updated.Description,
	}
}

func (s *srv) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := s.repoProduct.GetByGUID(ctx, guid)
	if err != nil {
		return err
	}

	return s.repoProduct.Delete(ctx, guid)
}
