package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chronos3344/catalog-service/internal/app/config"
	"github.com/chronos3344/catalog-service/internal/app/entity"
	hcategory "github.com/chronos3344/catalog-service/internal/app/handler/category"
	rhealth "github.com/chronos3344/catalog-service/internal/app/handler/health"
	hproduct "github.com/chronos3344/catalog-service/internal/app/handler/product"
	"github.com/chronos3344/catalog-service/internal/app/processor/http"
	pcategory "github.com/chronos3344/catalog-service/internal/app/repository/category"
	rcpostgres "github.com/chronos3344/catalog-service/internal/app/repository/conn/postgres"
	pproduct "github.com/chronos3344/catalog-service/internal/app/repository/product"
	mcategory "github.com/chronos3344/catalog-service/internal/app/service/category"
	mproduct "github.com/chronos3344/catalog-service/internal/app/service/product"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
	//categoryRepo:= pcategory.NewRepoFromPostgres(pgClient)
	//if err != nil {
	//	log.Fatalf("Failed to create category repository: %v", err)
	//}
	//
	//productRepo := pproduct.NewRepoFromPostgres(pgClient)
	//if err != nil {
	//	log.Fatalf("Failed to create product repository: %v", err)
	//}

	// === Временная проверка репозитория ===

	categoryRepo := pcategory.NewRepoFromPostgres(pgClient)
	productRepo := pproduct.NewRepoFromPostgres(pgClient)

	// 1. Создание категории
	cat := entity.Category{
		GUID:      uuid.Must(uuid.NewV4()),
		Name:      "Электроника",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := categoryRepo.Create(ctx, cat); err != nil {
		log.Fatal().Err(err).Msg("Category create failed")
	}
	log.Info().Str("guid", cat.GUID.String()).Msg("Category created")

	// 2. Получение категории
	found, err := categoryRepo.GetByGUID(ctx, cat.GUID)
	if err != nil {
		log.Fatal().Err(err).Msg("Category GetByGUID failed")
	}
	log.Info().Str("name", found.Name).Msg("Category found")

	// 3. Обновление категории
	found.Name = "Бытовая техника"
	found.UpdatedAt = time.Now()
	if err := categoryRepo.Update(ctx, found); err != nil {
		log.Fatal().Err(err).Msg("Category update failed")
	}
	log.Info().Msg("Category updated")

	// 4. List — все категории
	allCats, err := categoryRepo.List(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Category list failed")
	}
	log.Info().Int("count", len(allCats)).Msg("All categories")

	// 5. List — фильтр по имени
	filterName := "Бытовая техника"
	filtered, err := categoryRepo.List(ctx, &filterName)
	if err != nil {
		log.Fatal().Err(err).Msg("Category list by name failed")
	}
	log.Info().Int("count", len(filtered)).Msg("Filtered categories")

	// 6. Создание продукта
	desc := "Мощный пылесос"
	prod := entity.Product{
		GUID:         uuid.Must(uuid.NewV4()),
		Name:         "Пылесос Dyson V15",
		Description:  &desc,
		Price:        49999.99,
		CategoryGUID: cat.GUID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := productRepo.Create(ctx, prod); err != nil {
		log.Fatal().Err(err).Msg("Product create failed")
	}
	log.Info().Str("guid", prod.GUID.String()).Msg("Product created")

	// 7. List продуктов по категории
	products, err := productRepo.List(ctx, nil, &cat.GUID)
	if err != nil {
		log.Fatal().Err(err).Msg("Product list by category failed")
	}
	log.Info().Int("count", len(products)).Msg("Products in category")

	// 8. Удаление продукта, затем категории
	if err := productRepo.Delete(ctx, prod.GUID); err != nil {
		log.Fatal().Err(err).Msg("Product delete failed")
	}
	if err := categoryRepo.Delete(ctx, cat.GUID); err != nil {
		log.Fatal().Err(err).Msg("Category delete failed")
	}
	log.Info().Msg("Cleanup complete")

	// === Конец проверки ===

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
