package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/chronos3344/catalog-service/internal/app/config/section"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
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
	u.Host = cfg.Address // Ссылка на конфиг с Адресом
	u.User = url.UserPassword(cfg.Username, cfg.Password)
	u.Path = "/" + cfg.Name // Ссылка на конфиг с названием базы данных

	args := make(url.Values)
	args.Set("sslmode", "disable") // Создаем query param с флагом отключенного SSL (Протокол безопасности),
	// так как мы его не настроили.
	u.RawQuery = args.Encode() // И шифруем аргументы в наш URL.

	// В результате получим URL для подключения в нужном формате с отключенным SSL:
	// postgres://username:password@host:port/dbname?sslmode=disable

	dsn := u.String()

	fmt.Printf("PostgreSQL connection timeouts - Read: %v, Write: %v\n",
		cfg.ReadTimeout, cfg.WriteTimeout)

	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
		pgdriver.WithReadTimeout(cfg.ReadTimeout.Duration),
		pgdriver.WithWriteTimeout(cfg.WriteTimeout.Duration),
	))

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	rawBunDB := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	var cancelFunc func()
	ctx, cancelFunc = context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()

	if err := rawBunDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping connection: %w", err)

	}
	client := &Client{
		rawBunDB: rawBunDB,
		cfg:      cfg,
	}
	client._bunDB = rawBunDB
	return client, nil
}
