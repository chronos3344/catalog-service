package rcpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/chronos3344/catalog-service/internal/app/config/section"
	"github.com/chronos3344/catalog-service/migration"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
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
	u.Host = cfg.Address
	u.User = url.UserPassword(cfg.Username, cfg.Password)
	u.Path = "/" + cfg.Name

	args := make(url.Values)
	args.Set("sslmode", "disable")
	u.RawQuery = args.Encode()

	dsn := u.String()
	fmt.Printf("PostgreSQL connection URL: %s\n", dsn)

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
		_bunDB:   rawBunDB,
	}

	return client, nil
}

func (c *Client) Migrate(ctx context.Context) (oldVer, newVer int64, err error) {
	migrations := migrate.NewMigrations()
	if err = migrations.Discover(migration.Postgres); err != nil {
		return 0, 0, fmt.Errorf("failed to discover migrations: %w", err)
	}

	opts := []migrate.MigratorOption{
		migrate.WithTableName(c.cfg.MigrationTable),
		migrate.WithLocksTableName(c.cfg.MigrationTable + "_lock"),
		migrate.WithMarkAppliedOnSuccess(true),
	}

	m := migrate.NewMigrator(c.rawBunDB, migrations, opts...)

	if err = m.Init(ctx); err != nil {
		return 0, 0, fmt.Errorf("failed to init migrations table: %w", err)
	}

	applied, err := m.AppliedMigrations(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(applied) > 0 {
		oldVer, _ = strconv.ParseInt(applied[0].Name, 10, 64)
	}

	mgg, err := m.Migrate(ctx)
	if err != nil {
		return oldVer, oldVer, fmt.Errorf("failed to apply migrations: %w", err)
	}

	newVer = oldVer
	for _, mg := range mgg.Migrations {
		ver, _ := strconv.ParseInt(mg.Name, 10, 64)
		if ver > newVer {
			newVer = ver
		}
	}

	return oldVer, newVer, nil
}

//func extractMigrationVersion(name string) (int64, error) {
//	var version int64
//	_, err := fmt.Sscanf(name, "%d_", &version)
//	if err != nil {
//		return 0, fmt.Errorf("invalid migration name format: %s", name)
//	}
//	return version, nil
//}
