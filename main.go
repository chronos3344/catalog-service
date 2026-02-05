package main

import (
	"context"
	"log"

	"github.com/chronos3344/catalog-service/internal/app/config"
	rhealth "github.com/chronos3344/catalog-service/internal/app/handler/health"
	rprocessor "github.com/chronos3344/catalog-service/internal/app/processor/http"
	"github.com/chronos3344/catalog-service/internal/app/repository/conn/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	_, err = rcpostgres.NewConn(ctx, cfg.Repository.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	healthHandler := rhealth.NewHandler()

	server := rprocessor.NewHttp(healthHandler, cfg.Processor.WebServer)

	log.Printf("Starting catalog-service on port %d...", cfg.Processor.WebServer.ListenPort)
	if err := server.Serve(); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}

}
