package pproduct

import (
	"context"
	"database/sql"
	"errors"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/repository"
	rcpostgres "github.com/chronos3344/catalog-service/internal/app/repository/conn/postgres"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type repoPg struct {
	db *bun.DB
}

func NewRepoFromPostgres(_ context.Context, d *rcpostgres.Client) (repository.Product, error) {
	return &repoPg{db: d}, nil
}

func (r *repoPg) Create(ctx context.Context, product entity.Product) (entity.Product, error) {
	_, err := r.db.NewInsert().Model(&product).Exec(ctx)
	return product, err
}

func (r *repoPg) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	var product entity.Product
	err := r.db.NewSelect().Model(&product).Where("guid = ?", guid).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, repository.ErrNotFound
		}
		return entity.Product{}, err
	}
	return product, nil
}

func (r *repoPg) GetByName(ctx context.Context, name string) (entity.Product, error) {
	var product entity.Product
	err := r.db.NewSelect().Model(&product).Where("name = ?", name).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, repository.ErrNotFound
		}
		return entity.Product{}, err
	}
	return product, nil
}

func (r *repoPg) List(ctx context.Context, filter entity.RequestProductList) ([]entity.Product, error) {
	query := r.db.NewSelect().Model(&entity.Product{})

	if filter.CategoryGUID != nil {
		query = query.Where("category_guid = ?", *filter.CategoryGUID)
	}
	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}

	var products []entity.Product
	err := query.Scan(ctx)
	return products, err
}

func (r *repoPg) Update(ctx context.Context, product entity.Product) (entity.Product, error) {
	_, err := r.db.NewUpdate().Model(&product).WherePK().Exec(ctx)
	return product, err
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&entity.Product{}).Where("guid = ?", guid).Exec(ctx)
	return err
}
