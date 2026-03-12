package pproduct

import (
	"context"
	"database/sql"
	"errors"

	"github.com/chronos3344/catalog-service/internal/app/entity"
	"github.com/chronos3344/catalog-service/internal/app/repository"
	rcpostgres "github.com/chronos3344/catalog-service/internal/app/repository/conn/postgres"
	"github.com/google/uuid"
)

type (
	repoPg struct {
		*_DB
	}
	_DB = rcpostgres.Client
)

func NewRepoFromPostgres(_ context.Context, d *rcpostgres.Client) (repository.Product, error) {
	return &repoPg{_DB: d}, nil
}

func (r *repoPg) Create(ctx context.Context, product entity.Product) (entity.Product, error) {
	_, err := r._DB.NewInsert().Model(&product).Exec(ctx)
	return product, err
}

func (r *repoPg) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Product, error) {
	var product entity.Product
	err := r._DB.NewSelect().Model(&product).Where("guid = ?", guid).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Product{}, entity.ErrNotFound
		}
		return entity.Product{}, err
	}
	return product, nil
}

func (r *repoPg) List(ctx context.Context, filter entity.RequestProductList) ([]entity.Product, error) {
	query := r._DB.NewSelect().Model(&entity.Product{})

	if filter.Name != nil {
		query = query.Where("name = ?", *filter.Name)
	}
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
	err := query.Order("created_at DESC").Scan(ctx, &products)
	return products, err
}

func (r *repoPg) Update(ctx context.Context, product entity.Product) (entity.Product, error) {
	_, err := r._DB.NewUpdate().Model(&product).WherePK().Exec(ctx)
	return product, err
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := r._DB.NewDelete().Model(&entity.Product{}).Where("guid = ?", guid).Exec(ctx)
	return err
}
