package main

import (
	"log"

	"github.com/chronos3344/catalog-service/internal/app/config"
	_ "github.com/chronos3344/catalog-service/internal/app/config/section"
	rhealth "github.com/chronos3344/catalog-service/internal/app/handler/health"
	rprocessor "github.com/chronos3344/catalog-service/internal/app/processor/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	healthHandler := rhealth.NewHandler()

	httpServer := rprocessor.NewHttp(healthHandler, cfg.ProcessorWebServer)

	log.Fatal("HTTP server error:", httpServer.Serve())
}
