package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chronos3344/catalog-service/internal/app/config"
	hcategory "github.com/chronos3344/catalog-service/internal/app/handler/category"
	rhealth "github.com/chronos3344/catalog-service/internal/app/handler/health"
	hproduct "github.com/chronos3344/catalog-service/internal/app/handler/product"
	rprocessor "github.com/chronos3344/catalog-service/internal/app/processor/http"
	pcategory "github.com/chronos3344/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/chronos3344/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/chronos3344/catalog-service/internal/app/repository/product"
	mcategory "github.com/chronos3344/catalog-service/internal/app/service/category"
	mproduct "github.com/chronos3344/catalog-service/internal/app/service/product"
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
	defer pgClient.GetRawBunDB().Close()

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

	// Repositories
	categoryRepo, err := pcategory.NewRepoFromPostgres(ctx, pgClient)
	if err != nil {
		log.Fatalf("Failed to create category repository: %v", err)
	}

	productRepo, err := pproduct.NewRepoFromPostgres(ctx, pgClient)
	if err != nil {
		log.Fatalf("Failed to create product repository: %v", err)
	}

	// Services
	categoryService := mcategory.NewService(categoryRepo)
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
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server shutdown:", err)
	}
}
