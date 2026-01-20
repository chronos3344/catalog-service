package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App        App
	Repository Repository
	Broker     Broker
	Processor  Processor
	Monitor    Monitor
}

type Repository struct {
	Postgres PostgresConfig
	Service1 Service1Config
	Service2 Service2Config
}

type PostgresConfig struct {
	Address        string        `envconfig:"APP_REPOSITORY_POSTGRES_ADDRESS" default:"localhost:5432"`
	Username       string        `envconfig:"APP_REPOSITORY_POSTGRES_USERNAME"`
	Password       string        `envconfig:"APP_REPOSITORY_POSTGRES_PASSWORD"`
	Name           string        `envconfig:"APP_REPOSITORY_POSTGRES_NAME"`
	MigrationTable string        `envconfig:"APP_REPOSITORY_POSTGRES_MIGRATION_TABLE" default:"schema_migrations"`
	ReadTimeout    time.Duration `envconfig:"APP_REPOSITORY_POSTGRES_READ_TIMEOUT" default:"30s"`
	WriteTimeout   time.Duration `envconfig:"APP_REPOSITORY_POSTGRES_WRITE_TIMEOUT" default:"30s"`
}

// Добавить эти структуры (минимум):
type Service1Config struct {
	// добавьте нужные поля
}

type Service2Config struct {
	// добавьте нужные поля
}

type App struct {
	// добавьте поля
}

type Broker struct {
	// добавьте поля
}

type Processor struct {
	// добавьте поля
}

type Monitor struct {
	// добавьте поля
}

func Load() (*Config, error) {
	envPath := os.Getenv("ENV_FILE_PATH")
	if envPath == "" {
		envPath = ".env"
	}

	if err := godotenv.Load(envPath); err != nil {
		fmt.Printf("Note: .env file not found at %s, using environment variables\n", envPath)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}

	return &cfg, nil
}

//---------------------------

//type Config struct {
//	App        section.App
//	Repository section.Repository
//	Broker     section.Broker
//	Processor  section.Processor
//	Monitor    section.Monitor
//}
//
//type Repository struct {
//	Postgres        RepositoryPostgres
//	Service1        RepositoryService1
//	Service2        RepositoryService2
//}
//
//type RepositoryPostgres struct {
//	Address        string        `envconfig:"APP_REPOSITORY_POSTGRES_ADDRESS"             default:"localhost:5429"`
//	Username       string        `envconfig:"APP_REPOSITORY_POSTGRES_USERNAME"`
//	Password       string        `envconfig:"APP_REPOSITORY_POSTGRES_PASSWORD"`
//	Name           string        `envconfig:"APP_REPOSITORY_POSTGRES_NAME"`
//	MigrationTable string        `envconfig:"APP_REPOSITORY_POSTGRES_MIGRATION_TABLE"     default:"schema_migrations"`
//	ReadTimeout    time.Duration `envconfig:"APP_REPOSITORY_POSTGRES_READ_TIMEOUT"        default:"30s"`
//	WriteTimeout   time.Duration `envconfig:"APP_REPOSITORY_POSTGRES_WRITE_TIMEOUT"       default:"30s"`
//
//// Load загружает конфигурацию из файла .env и переменных окружения
//func Load() (*Config, error) {
//	// Определяем путь к файлу .env (можно передать через переменную окружения)
//	envPath := os.Getenv("ENV_FILE_PATH")
//	if envPath == "" {
//	envPath = ".env"
//	}
//
//	// Загружаем переменные из .env файла
//	if err := godotenv.Load(envPath); err != nil {
//		// Не считаем ошибкой отсутствие .env файла, так как переменные могут быть установлены в окружении
//		fmt.Printf("Note: .env file not found at %s, using environment variables\n", envPath)
//	}
//
//	// Создаем пустую конфигурацию
//	var cfg Config
//
//	// Заполняем структуру конфигурации из переменных окружения
//	if err := envconfig.Process("", &cfg); err != nil {
//		return nil, fmt.Errorf("failed to process environment variables: %w", err)
//	}
//
//	return &cfg, nil
//}

//-----------------------
//type Config struct {
//	// Server содержит настройки сервера
//	Server ServerConfig
//	// Database содержит настройки базы данных
//	Database DatabaseConfig
//	// Logging содержит настройки логирования
//	Logging LoggingConfig
//	// Cache содержит настройки кэширования
//	Cache CacheConfig
//}
//
//// DatabaseConfig содержит настройки подключения к базе данных
//type DatabaseConfig struct {
//	Host     string `envconfig:"DB_HOST" default:"localhost"`
//	Port     int    `envconfig:"DB_PORT" default:"5432"`
//	User     string `envconfig:"DB_USER" default:"postgres"`
//	Password string `envconfig:"DB_PASSWORD" default:"postgres"`
//	Name     string `envconfig:"DB_NAME" default:"catalog"`
//	SSLMode  string `envconfig:"DB_SSL_MODE" default:"disable"`
//	MaxConns int    `envconfig:"DB_MAX_CONNS" default:"20"`
//}

//// Load загружает конфигурацию из файла .env и переменных окружения
//func Load() (Config, error) {
//	var cfg Config
//	// Загружаем переменные из файла .env, если он существует
//	// Игнорируем ошибку, если файл не найден
//	err := godotenv.Load()
//	if err != nil {
//		log.Printf("Warning: .env file not found, using environment variables")
//	}
//
//	// Парсим переменные окружения в структуру
//	err = envconfig.Process("", &cfg)
//	if err != nil {
//		return cfg, fmt.Errorf("failed to parse environment variables: %w", err)
//	}
//
//	return cfg, nil
//}
