package pcategory

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

func NewRepoFromPostgres(_ context.Context, d *rcpostgres.Client) (repository.Category, error) {
	return &repoPg{db: d.GetRawBunDB()}, nil
}

func (r *repoPg) Create(ctx context.Context, category entity.Category) (entity.Category, error) {
	_, err := r.db.NewInsert().Model(&category).Exec(ctx)
	return category, err
}

func (r *repoPg) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Category, error) {
	var category entity.Category
	err := r.db.NewSelect().Model(&category).Where("guid = ?", guid).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Category{}, repository.ErrNotFound
		}
		return entity.Category{}, err
	}
	return category, nil
}

func (r *repoPg) GetByName(ctx context.Context, name string) (entity.Category, error) {
	var category entity.Category
	err := r.db.NewSelect().Model(&category).Where("name = ?", name).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Category{}, repository.ErrNotFound
		}
		return entity.Category{}, err
	}
	return category, nil
}

func (r *repoPg) List(ctx context.Context) ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.NewSelect().Model(&categories).Scan(ctx)
	return categories, err
}

func (r *repoPg) Update(ctx context.Context, category entity.Category) (entity.Category, error) {
	_, err := r.db.NewUpdate().Model(&category).WherePK().Exec(ctx)
	return category, err
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&entity.Category{}).Where("guid = ?", guid).Exec(ctx)
	return err
}
