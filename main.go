package main

import (
	"fmt"
	"log"

	"github.com/chronos3344/catalog-service/internal/app/config"
)

func main() {
	s := "Vladimir"
	fmt.Printf("Hello and welcome, %s!/n", s)

	// Простая загрузка
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Или принудительная загрузка (завершит программу при ошибке)
	// cfg := config.MustLoad()

	fmt.Printf("Server running on %s:%s/n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Database: %s@%s:%d/%s/n",
		cfg.DB.User,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name)
}
