package main

import (
	"log"

	"github.com/chronos3344/catalog-service/internal/app/config"
	rhealth "github.com/chronos3344/catalog-service/internal/app/handler/health"
	rprocessor "github.com/chronos3344/catalog-service/internal/app/processor/http"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаем health handler
	healthHandler := rhealth.NewHandler()
	// Создаем HTTP сервер
	server := rprocessor.NewHttp(healthHandler, cfg.Processor.WebServer)
	// Запускаем сервер
	log.Printf("Starting catalog-service on port %d...", cfg.Processor.WebServer.ListenPort)
	if err := server.Serve(); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
