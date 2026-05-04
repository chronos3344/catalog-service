package config

import (
	"io"
	"time"

	"github.com/chronos3344/catalog-service/internal/app/config/section"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Repository section.Repository
	Processor  section.Processor
	Monitor    section.Monitor
}

var Root Config

type LoadArgs struct {
	Output          io.Writer `json:"-"`
	EnableSimpleLog bool
}

func createLogger(level zerolog.Level, output io.Writer) zerolog.Logger {
	return zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()
}

func Load(args LoadArgs) {
	// Этап 1: Предварительный логгер с Debug уровнем

	// Переименование стандартных полей zerolog:
	// "time" -> "timestamp", "message" -> "msg"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.MessageFieldName = "msg"
	zerolog.TimeFieldFormat = time.RFC3339

	if args.EnableSimpleLog {
		args.Output = zerolog.ConsoleWriter{Out: args.Output}
	}

	// log.Logger — глобальная переменная пакета zerolog/log.

	// Все вызовы log.Info(), log.Error() и т.д. используют этот логгер.
	// Перезаписывая его, мы меняем поведение логирования во всём приложении.
	log.Logger = createLogger(zerolog.DebugLevel, args.Output)

	log.Debug().Msg("Logger initialized with Debug level")

	// Загрузка конфигурации
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("No .env file found")
	}

	if err := envconfig.Process("APP", &Root); err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// TODO: Этап 2 — переинициализация с уровнем из конфига
	// 1. Распарсите уровень: zerolog.ParseLevel(Root.Monitor.LogLevel)
	// 2. Пересоздайте логгер: log.Logger = createLogger(level, args.Output)
	// 3. Залогируйте успешную инициализацию на уровне Info

	level, err := zerolog.ParseLevel(Root.Monitor.LogLevel)
	if err != nil {
		log.Warn().Str("log_level", Root.Monitor.LogLevel).Msg("Unknown log level, using debug")
		level = zerolog.DebugLevel
	}

	log.Logger = createLogger(level, args.Output)
	log.Info().Str("log_level", level.String()).Msg("Logger re-initialized with config level")
}
