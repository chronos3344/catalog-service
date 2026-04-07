package main

import (
	"context"
	"github.com/uptrace/bun"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chronos3344/catalog-service/internal/app/config"
	hcategory "github.com/chronos3344/catalog-service/internal/app/handler/category"
	rhealth "github.com/chronos3344/catalog-service/internal/app/handler/health"
	hproduct "github.com/chronos3344/catalog-service/internal/app/handler/product"
	"github.com/chronos3344/catalog-service/internal/app/processor/http"
	pcategory "github.com/chronos3344/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/chronos3344/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/chronos3344/catalog-service/internal/app/repository/product"
	mcategory "github.com/chronos3344/catalog-service/internal/app/service/category"
	mproduct "github.com/chronos3344/catalog-service/internal/app/service/product"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config: %v")
	}

	ctx := context.Background()

	// Создаем подключение к PostgreSQL
	pgClient, err := rcpostgres.NewConn(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL: %v")
	}
	defer func(db *bun.DB) {
		err := db.Close()
		if err != nil {

		}
	}(pgClient.GetRawBunDB())

	// Применение миграций
	oldVer, newVer, err := pgClient.Migrate(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations: %v")
	}

	if oldVer != newVer {
		log.Printf("Database migrated from version %d to %d", oldVer, newVer)
	} else {
		log.Printf("Database is up to date, version: %d", newVer)
	}

	// Repositories
	categoryRepo := pcategory.NewRepoFromPostgres(pgClient)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create category repository: %v")
	}

	productRepo := pproduct.NewRepoFromPostgres(pgClient)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create product repository: %v")
	}

	// Services
	categoryService := mcategory.NewService(categoryRepo, productRepo)
	productService := mproduct.NewService(productRepo, categoryRepo)

	// Handlers
	categoryHandler := hcategory.NewHandler(categoryService)
	productHandler := hproduct.NewHandler(productService)
	healthHandler := rhealth.NewHandler()

	// Server - создаем роутер и регистрируем все обработчики
	server := rprocessor.NewHttp(healthHandler, categoryHandler, productHandler, cfg.Processor.WebServer)

	// Graceful shutdown
	go func() {
		log.Printf("Starting catalog-service on port %d...", cfg.Processor.WebServer.ListenPort)
		if err := server.Serve(); err != nil {
			log.Fatal().Err(err).Msg("Failed to start HTTP server: %v")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Print("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatal().Err(err).Msg("Server shutdown:")
	}
}
