package pproduct

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/repository"
	rcpostgres "github.com/chronos3344/catalog-service/internal/app/repository/conn/postgres"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type (
	repoPg struct {
		*rcpostgres.Client
	}
)

func NewRepoFromPostgres(_ context.Context, d *rcpostgres.Client) (repository.Product, error) {
	return &repoPg{d}, nil
}

func (r *repoPg) Create(ctx context.Context, product entity.Product) (entity.Product, error) {
	_, err := r.NewInsert().
		Model(&product).
		Returning("*").
		Exec(ctx)

	if err != nil {
		return entity.Product{}, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (r *repoPg) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	var product entity.Product

	err := r.NewSelect().
		Model(&product).
		Relation("Category").
		Where("p.guid = ?", guid).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, fmt.Errorf("product not found")
		}
		return entity.Product{}, fmt.Errorf("failed to get product: %w", err)
	}

	return product, nil
}

func (r *repoPg) GetByNameAndCategory(ctx context.Context, name string, categoryGUID uuid.UUID) (entity.Product, error) {
	var product entity.Product

	err := r.NewSelect().
		Model(&product).
		Where("name = ? AND category_guid = ?", name, categoryGUID).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, nil
		}
		return entity.Product{}, fmt.Errorf("failed to get product by name and category: %w", err)
	}

	return product, nil
}

func (r *repoPg) Update(ctx context.Context, product entity.Product) (entity.Product, error) {
	product.UpdatedAt = bun.Now()

	_, err := r.NewUpdate().
		Model(&product).
		WherePK().
		Returning("*").
		Exec(ctx)

	if err != nil {
		return entity.Product{}, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	product := &entity.Product{GUID: guid}

	_, err := r.NewDelete().
		Model(product).
		WherePK().
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (r *repoPg) List(ctx context.Context, page, limit int, categoryGUID *uuid.UUID) ([]entity.Product, int, error) {
	var products []entity.Product

	offset := (page - 1) * limit

	query := r.NewSelect().
		Model(&products).
		Relation("Category").
		Order("p.created_at DESC").
		Limit(limit).
		Offset(offset)

	if categoryGUID != nil {
		query = query.Where("p.category_guid = ?", *categoryGUID)
	}

	count, err := query.ScanAndCount(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	return products, count, nil
}
