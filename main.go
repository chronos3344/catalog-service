package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

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
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// DB connection
	dsn := "postgres://user:password@localhost:5432/catalog?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

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

	// Repositories
	categoryRepo, _ := pcategory.NewRepoFromPostgres(ctx, &db)
	productRepo, _ := pproduct.NewRepoFromPostgres(ctx, &db)

	// Services
	categoryService := mcategory.NewService(categoryRepo)
	productService := mproduct.NewService(productRepo, categoryRepo)

	// Handlers
	categoryHandler := hcategory.NewHandler(categoryService)
	productHandler := hproduct.NewHandler(productService)

	// Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: rprocessor.NewHttp(nil, categoryHandler, productHandler),
	}

	go func() {
		log.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server shutdown:", err)
	}
}

//func main() {
//	cfg, err := config.Load()
//	if err != nil {
//		log.Fatalf("Failed to load config: %v", err)
//	}
//
//	ctx := context.Background()
//
