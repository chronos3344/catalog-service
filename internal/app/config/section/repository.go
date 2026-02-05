package section

import (
	"github.com/chronos3344/catalog-service/internal/app/util"
)

type Repository struct {
	Postgres RepositoryPostgres
}
type RepositoryPostgres struct {
	Address      string        `validate:"required"`
	Name         string        `validate:"required"`
	Username     string        `validate:"required"`
	Password     string        `validate:"required"`
	ReadTimeout  util.Duration `validate:"required" split_words:"true"`
	WriteTimeout util.Duration `validate:"required" split_words:"true"`
}
