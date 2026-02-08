package pcategory

import (
	"context"

	rcpostgres "catalog-service/internal/app/repository/conn/postgres"

	"github.com/chronos3344/catalog-service/internal/app/repository"
)

type (
	repoPg struct {
		*_DB
	}

	_DB = rcpostgres.Client // Наш клиент, который мы получаем, когда инициализируем соединение с БД
)

// Наша функция-билдер
func NewRepoFromPostgres(_ ctx context.Context, d *rcpostgres.Client) (repository.Category, error) {
	return &repoPg{d}, nil
}

// Далее нужно имплементировать наш интерфейс и реализовать методы
func (r *repoPg) Create(ctx context.Context, category entity.Category) (entity.Category, error) {
	panic("implement me!")

	// Здесь вы должны создать операцию создания в БД
}