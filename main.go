package main

import (
	"context"
	"log"

	"github.com/chronos3344/catalog-service/internal/app/config"
	rhealth "github.com/chronos3344/catalog-service/internal/app/handler/health"
	rprocessor "github.com/chronos3344/catalog-service/internal/app/processor/http"
	rcpostgres "github.com/chronos3344/catalog-service/internal/app/repository/conn/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	// Создаем подключение к PostgreSQL
	pgClient, err := rcpostgres.NewConn(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Применение миграций
	oldVer, newVer, err := pgClient.Migrate(ctx)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if oldVer != newVer {
		log.Printf("Database migrated from version %d to %d", oldVer, newVer)
	} else {
		log.Printf("Database is up to date, version: %d", newVer)
	}

	healthHandler := rhealth.NewHandler()
	server := rprocessor.NewHttp(healthHandler, cfg.Processor.WebServer)

	log.Printf("Starting catalog-service on port %d...", cfg.Processor.WebServer.ListenPort)
	if err := server.Serve(); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
