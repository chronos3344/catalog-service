package pcategory

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
		*_DB
	}

	_DB = rcpostgres.Client
)

func NewRepoFromPostgres(_ context.Context, d *rcpostgres.Client) (repository.Category, error) {
	return &repoPg{d}, nil
}

func (r *repoPg) Create(ctx context.Context, category entity.Category) (entity.Category, error) {
	_, err := r.NewInsert().
		Model(&category).
		Returning("*").
		Exec(ctx)

	if err != nil {
		return entity.Category{}, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (r *repoPg) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Category, error) {
	var category entity.Category

	err := r.NewSelect().
		Model(&category).
		Where("guid = ?", guid).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Category{}, fmt.Errorf("category not found")
		}
		return entity.Category{}, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

func (r *repoPg) GetByName(ctx context.Context, name string) (entity.Category, error) {
	var category entity.Category

	err := r.NewSelect().
		Model(&category).
		Where("name = ?", name).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Category{}, nil
		}
		return entity.Category{}, fmt.Errorf("failed to get category by name: %w", err)
	}

	return category, nil
}

func (r *repoPg) Update(ctx context.Context, category entity.Category) (entity.Category, error) {
	category.UpdatedAt = bun.Now()

	_, err := r.NewUpdate().
		Model(&category).
		WherePK().
		Returning("*").
		Exec(ctx)

	if err != nil {
		return entity.Category{}, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	category := &entity.Category{GUID: guid}

	_, err := r.NewDelete().
		Model(category).
		WherePK().
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (r *repoPg) List(ctx context.Context, page, limit int) ([]entity.Category, int, error) {
	var categories []entity.Category

	offset := (page - 1) * limit

	count, err := r.NewSelect().
		Model(&categories).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, count, nil
}
