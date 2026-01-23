package section

type ProcessorWebServer struct {
	Port int `yaml:"port" env:"HTTP_PORT" env-default:"9091"`
}
