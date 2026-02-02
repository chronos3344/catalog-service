package rcpostgres

import (
	"context"
	"net/url"

	"github.com/uptrace/bun"

	"github.com/chronos3344/catalog-service/internal/app/config/section"
)

type (
	Client struct {
		_bunDB
		rawBunDB *bun.DB // Для служебных целей

		cfg section.RepositoryPostgres
	}

	_bunDB = bun.IDB
)

func (c *Client) GetRawBunDB() *bun.DB {
	return c.rawBunDB
}

func NewConn(ctx context.Context, cfg section.RepositoryPostgres) (*Client, error) {
	var u url.URL
	u.Scheme = "postgres"
	u.Host = cfg.Address  // Ссылка на конфиг с Адресом
	u.User = url.UserPassword(cfg.Username, cfg.Password)
	u.Path = "/" + cfg.Name // Ссылка на конфиг с названием базы данных

	var args make(url.Values)
	args.Set("sslmode", "disable") // Создаем query param с флагом отключенного SSL (Протокол безопасности),
									// так как мы его не настроили.
	u.RawQuery = args.Encode()     // И шифруем аргументы в наш URL.

	// В результате получим URL для подключения в нужном формате с отключенным SSL:
	// postgres://username:password@host:port/dbname?sslmode=disable
}
