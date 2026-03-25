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

func (r *repoPg) Create(ctx context.Context, product entity.Product) error {
	_, err := r._DB.NewInsert().Model(&product).Exec(ctx)
	return err
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

func (r *repoPg) List(ctx context.Context, name *string, categoryGUID *uuid.UUID) ([]entity.Product, error) {
	var products []entity.Product
	query := r._DB.NewSelect().Model(&products).Order("created_at DESC").Where("categoryGUID = ?", categoryGUID)

	if name != nil {
		query = query.Where("name = ?", name)
	}

	err := query.Scan(ctx)

	return products, err
}

func (r *repoPg) Update(ctx context.Context, product entity.Product) error {
	res, err := r._DB.NewUpdate().Model(&product).WherePK().OmitZero().Exec(ctx)
	return rcpostgres.UpdateErr(res, err)
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := r._DB.NewDelete().Model(&entity.Product{}).Where("guid = ?", guid).Exec(ctx)
	return rcpostgres.DeleteErr(err)
}
