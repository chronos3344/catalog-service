package section

import (
	"github.com/chronos3344/catalog-service/internal/app/util"
)

// RepositoryPostgres конфигурация PostgreSQL
type RepositoryPostgres struct {
	Address      string        `env:"APP_REPOSITORY_POSTGRES_ADDRESS" validate:"required"`
	Name         string        `env:"APP_REPOSITORY_POSTGRES_NAME" validate:"required"`
	Username     string        `env:"APP_REPOSITORY_POSTGRES_USERNAME" validate:"required"`
	Password     string        `env:"APP_REPOSITORY_POSTGRES_PASSWORD" validate:"required"`
	ReadTimeout  util.Duration `env:"APP_REPOSITORY_POSTGRES_READ_TIMEOUT" validate:"required"`
	WriteTimeout util.Duration `env:"APP_REPOSITORY_POSTGRES_WRITE_TIMEOUT" validate:"required"`
}

// Repository секция конфигурации репозиториев
type Repository struct {
	Postgres RepositoryPostgres
}
