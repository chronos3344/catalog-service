//package main
//
//import (
//	"log"
//
//	"github.com/chronos3344/catalog-service/internal/app/config"
//	_ "github.com/chronos3344/catalog-service/internal/app/config/section"
//	rhealth "github.com/chronos3344/catalog-service/internal/app/handler/health"
//	rprocessor "github.com/chronos3344/catalog-service/internal/app/processor/http"
//)
//
//func main() {
//	cfg, err := config.Load()
//	if err != nil {
//		log.Fatal("Failed to load config:", err)
//	}
//
//	healthHandler := rhealth.NewHandler()
//
//	httpServer := rprocessor.NewHttp(healthHandler, cfg.ProcessorWebServer)
//
//	log.Fatal("HTTP server error:", httpServer.Serve())
//}

package main

import (
	"log"

	"github.com/chronos3344/catalog-service/internal/app/config"
)

func main() {
	config.Load()

	cfg := config.Root

	log.Printf("Server will start on port: %d", cfg.Processor.WebServer.ListenPort)
	log.Printf("Database: %s@%s/%s",
		cfg.Repository.Postgres.Username,
		cfg.Repository.Postgres.Address,
		cfg.Repository.Postgres.Name)
	log.Printf("Environment: %s, LogLevel: %s",
		cfg.Monitor.Environment,
		cfg.Monitor.LogLevel)
}
