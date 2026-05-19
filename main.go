package main

import (
	"context"
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
	ctx := context.Background()

	config.Load(config.LoadArgs{
		Output:          os.Stdout,
		EnableSimpleLog: true,
	})

	cfg := config.Root

	pgClient, err := rcpostgres.NewClient(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	// Применение миграций
	oldVer, newVer, err := pgClient.Migrate(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}
	if oldVer != newVer {
		log.Info().Int64("old_version", oldVer).Int64("new_version", newVer).Msg("Database migrated")
	} else {
		log.Info().Int64("version", newVer).Msg("Database is up to date")
	}

	// Repositories
	categoryRepo := pcategory.NewRepoFromPostgres(pgClient)
	productRepo := pproduct.NewRepoFromPostgres(pgClient)

	// Services
	categoryService := mcategory.NewService(categoryRepo, productRepo)
	productService := mproduct.NewService(productRepo, categoryRepo)

	// Handlers
	categoryHandler := hcategory.NewHandler(categoryService)
	productHandler := hproduct.NewHandler(productService)
	healthHandler := rhealth.NewHandler()

	// Server
	server := rprocessor.NewHttp(healthHandler, categoryHandler, productHandler, cfg.Processor.WebServer)

	// Graceful shutdown
	go func() {
		log.Info().Msgf("Starting catalog-service on port %d...", cfg.Processor.WebServer.ListenPort)
		if err := server.Serve(); err != nil {
			log.Fatal().Err(err).Msg("Failed to start HTTP server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msgf("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Error().Err(err).Msg("Server shutdown error")
	}
}
