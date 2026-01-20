package main

import (
	"fmt"
	"log"

	"github.com/chronos3344/catalog-service/internal/app/config"
)

func main() {
	s := "Vladimir"
	fmt.Printf("Hello and welcome, %s!\n", s)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	log.Printf("Server will start on port: %d", cfg.Server.Port)
	log.Printf("Database DSN: postgresql://%s:$s@%s/%s",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Address,
		cfg.Database.Name)

}
