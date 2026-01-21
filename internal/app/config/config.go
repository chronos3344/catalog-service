package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config - основная структура конфигурации
type Config struct {
	DB     DatabaseConfig
	Server ServerConfig
	Log    LogConfig
}

// DatabaseConfig - конфигурация базы данных
type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     int    `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"test_user"`
	Password string `envconfig:"DB_PASSWORD" required:"test_password"`
	Name     string `envconfig:"DB_NAME" default:"test_db"`
}

// ServerConfig - конфигурация сервера
type ServerConfig struct {
	Host         string `envconfig:"SERVER_HOST" default:"0.0.0.0"`
	Port         string `envconfig:"SERVER_PORT" default:"8080"`
	ReadTimeout  string `envconfig:"SERVER_READ_TIMEOUT" default:"30s"`
	WriteTimeout string `envconfig:"SERVER_WRITE_TIMEOUT" default:"30s"`
}

// LogConfig - конфигурация логирования
type LogConfig struct {
	Level  string `envconfig:"LOG_LEVEL" default:"info"`
	Format string `envconfig:"LOG_FORMAT" default:"json"`
}

// Load - загружает конфигурацию из .env файла и переменных окружения
func Load() (Config, error) {
	// Загружаем переменные из .env файла (если он существует)
	// Игнорируем ошибку, если файл не найден
	_ = godotenv.Load()

	var cfg Config

	// Парсим переменные окружения в структуру
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// MustLoad - загружает конфигурацию или завершает работу при ошибке
func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return cfg
}
