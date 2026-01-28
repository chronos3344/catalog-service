package config

import (
	"log"

	"github.com/chronos3344/catalog-service/internal/app/config/section"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Repository section.Repository
	Processor  section.Processor
	Monitor    section.Monitor
}

var Root Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env file not found")
	}

	if err := envconfig.Process("APP", &Root); err != nil {
		log.Fatalf("Failed parse: %v", err)
	}
	log.Println("Config loaded")
}
