package pcategory

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
	_DB = rcpostgres.Client // Наш клиент, который мы получаем, когда инициализируем соединение с БД
)

func NewRepoFromPostgres(_ context.Context, d *rcpostgres.Client) (repository.Category, error) {
	return &repoPg{_DB: d}, nil
}

func (r *repoPg) Create(ctx context.Context, category entity.Category) (entity.Category, error) {
	_, err := r._DB.NewInsert().Model(&category).Exec(ctx)
	return category, err
}

func (r *repoPg) GetByGUID(ctx context.Context, guid uuid.UUID) (entity.Category, error) {
	var category entity.Category
	err := r._DB.NewSelect().Model(&category).Where("guid = ?", guid).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Category{}, entity.ErrNotFound
		}
		return entity.Category{}, err
	}
	return category, nil
}

func (r *repoPg) GetByName(ctx context.Context, name string) (entity.Category, error) {
	var category entity.Category
	err := r._DB.NewSelect().Model(&category).Where("name = ?", name).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Category{}, entity.ErrNotFound
		}
		return entity.Category{}, err
	}
	return category, nil
}

func (r *repoPg) List(ctx context.Context) ([]entity.Category, error) {
	var categories []entity.Category
	err := r._DB.NewSelect().Model(&categories).Scan(ctx)
	return categories, err
}

func (r *repoPg) Update(ctx context.Context, category entity.Category) (entity.Category, error) {
	_, err := r._DB.NewUpdate().Model(&category).WherePK().Exec(ctx)
	return category, err
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := r._DB.NewDelete().Model(&entity.Category{}).Where("guid = ?", guid).Exec(ctx)
	return err
}
