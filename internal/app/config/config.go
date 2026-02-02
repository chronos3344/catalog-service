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

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env file not found")
	}

	var cfg Config
	if err := envconfig.Process("APP", &cfg); err != nil {
		return nil, err
	}

	Root = cfg
	log.Println("Config loaded successfully")
	return &cfg, nil

}
