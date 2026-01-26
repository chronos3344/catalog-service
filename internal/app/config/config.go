//package config
//
//import (
//	"log"
//
//	"github.com/chronos3344/catalog-service/internal/app/config/section"
//	"github.com/joho/godotenv"
//	"github.com/kelseyhightower/envconfig"
//)
//
//// Config - основная структура конфигурации
//type Config struct {
//	DB                 DatabaseConfig
//	Server             ServerConfig
//	Log                LogConfig
//	ProcessorWebServer section.ProcessorWebServer
//}
//
//// DatabaseConfig - конфигурация базы данных
//type DatabaseConfig struct {
//	Host     string `envconfig:"APP_REPOSITORY_POSTGRES_HOST" default:"localhost"`
//	Port     int    `envconfig:"APP_REPOSITORY_POSTGRES_PORT" default:"5432"`
//	User     string `envconfig:"APP_REPOSITORY_POSTGRES_USER" default:"test_user"`
//	Password string `envconfig:"APP_REPOSITORY_POSTGRES_PASSWORD" required:"test_password"`
//	Name     string `envconfig:"APP_REPOSITORY_POSTGRES_NAME" default:"test_db"`
//}
//
//// ServerConfig - конфигурация сервера
//type ServerConfig struct {
//	Host         string `envconfig:"APP_PROCESSOR_WEB_SERVER_LISTEN_HOST" default:"0.0.0.0"`
//	Port         string `envconfig:"APP_PROCESSOR_WEB_SERVER_LISTEN_PORT" default:"8080"`
//	ReadTimeout  string `envconfig:"APP_PROCESSOR_WEB_SERVER_LISTEN_TIMEOUT" default:"30s"`
//	WriteTimeout string `envconfig:"APP_PROCESSOR_WEB_SERVER_LISTEN_TIMEOUT" default:"30s"`
//}
//
//// LogConfig - конфигурация логирования
//type LogConfig struct {
//	Level  string `envconfig:"APP_MONITOR_LOG_LEVEL" default:"info"`
//	Format string `envconfig:"APP_MONITOR_LOG_FORMAT" default:"json"`
//}
//
//// Load - загружает конфигурацию из .env файла и переменных окружения
//func Load() (Config, error) {
//	// Загружаем переменные из .env файла (если он существует)
//	// Игнорируем ошибку, если файл не найден
//	_ = godotenv.Load()
//
//	var cfg Config
//
//	// Парсим переменные окружения в структуру
//	err := envconfig.Process("", &cfg)
//	if err != nil {
//		return Config{}, err
//	}
//
//	return cfg, nil
//}

package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/chronos3344/catalog-service/internal/app/config/section"
	"github.com/chronos3344/catalog-service/internal/app/util"
	"github.com/joho/godotenv"
)

// Root глобальная конфигурация приложения
var Root = struct {
	Repository section.Repository
	Processor  section.Processor
	Monitor    section.Monitor
}{}

// Load загружает конфигурацию из .env файла и переменных окружения
func Load() {
	// Загружаем .env файл (не критично, если его нет)
	_ = godotenv.Load(".env")

	// Загружаем конфигурацию из структуры
	if err := loadConfig(&Root, "APP"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Config loaded successfully")
}

// loadConfig рекурсивно загружает конфигурацию из переменных окружения
func loadConfig(config interface{}, prefix string) error {
	val := reflect.ValueOf(config).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Если поле является структурой, рекурсивно обрабатываем её
		if field.Kind() == reflect.Struct {
			newPrefix := prefix
			if fieldType.Tag.Get("env") != "" {
				newPrefix = fieldType.Tag.Get("env")
			} else {
				newPrefix = fmt.Sprintf("%s_%s", prefix, strings.ToUpper(fieldType.Name))
			}

			if err := loadConfig(field.Addr().Interface(), newPrefix); err != nil {
				return err
			}
			continue
		}

		// Получаем тег env для поля
		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			envTag = fmt.Sprintf("%s_%s", prefix, strings.ToUpper(fieldType.Name))
		}

		// Получаем значение из переменных окружения
		envValue := os.Getenv(envTag)
		if envValue == "" {
			// Проверяем required из тега validate
			if strings.Contains(fieldType.Tag.Get("validate"), "required") {
				return fmt.Errorf("missing required field: %s", envTag)
			}
			continue
		}

		// Устанавливаем значение в поле
		if err := setFieldValue(&field, envValue); err != nil {
			return fmt.Errorf("failed to set field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// setFieldValue устанавливает значение поля на основе его типа
func setFieldValue(field *reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Для кастомного типа Duration
		if field.Type().String() == "util.Duration" {
			var dur util.Duration
			if err := dur.UnmarshalText([]byte(value)); err != nil {
				return err
			}
			field.Set(reflect.ValueOf(dur))
		} else {
			// Для обычных int
			var intVal int64
			fmt.Sscanf(value, "%d", &intVal)
			field.SetInt(intVal)
		}
	default:
		// Пытаемся использовать TextUnmarshaler интерфейс
		if field.CanInterface() {
			if unmarshaler, ok := field.Interface().(interface {
				UnmarshalText([]byte) error
			}); ok {
				return unmarshaler.UnmarshalText([]byte(value))
			}
		}
		return fmt.Errorf("unsupported field type: %s", field.Type())
	}
	return nil
}
