package section

// Monitor конфигурация мониторинга и логирования
type Monitor struct {
	LogLevel    string `env:"APP_MONITOR_LOG_LEVEL" validate:"required,oneof=debug info warn error"`
	Environment string `env:"APP_MONITOR_ENVIRONMENT" validate:"required,oneof=development staging production"`
}
