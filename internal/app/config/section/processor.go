package section

// ProcessorWebServer конфигурация веб-сервера
type ProcessorWebServer struct {
	ListenPort int `env:"APP_PROCESSOR_WEB_SERVER_LISTEN_PORT" validate:"required"`
}

// Processor секция конфигурации процессоров
type Processor struct {
	WebServer ProcessorWebServer
}
