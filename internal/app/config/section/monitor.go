package section

// Monitor конфигурация мониторинга и логирования
type Monitor struct {
	LogLevel    string `split_words:"true" default:"debug"`
	Environment string `default:"development"`
}
