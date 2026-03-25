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

func (r *repoPg) Create(ctx context.Context, category entity.Category) error {
	_, err := r._DB.NewInsert().Model(&category).Exec(ctx)
	return err
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

func (r *repoPg) List(ctx context.Context, name *string) ([]entity.Category, error) {
	var categories []entity.Category
	query := r._DB.NewSelect().Model(&categories)

	if name != nil {
		query = query.Where("name = ?", name)
	}

	err := query.Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, err
}

func (r *repoPg) Update(ctx context.Context, category entity.Category) error {
	res, err := r._DB.NewUpdate().Model(&category).WherePK().OmitZero().Exec(ctx)
	return rcpostgres.UpdateErr(res, err)
}

func (r *repoPg) Delete(ctx context.Context, guid uuid.UUID) error {
	_, err := r._DB.NewDelete().Model(&entity.Category{}).Where("guid = ?", guid).Exec(ctx)
	return rcpostgres.DeleteErr(err)
}
