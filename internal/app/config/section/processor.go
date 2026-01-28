package section

type Processor struct {
	WebServer ProcessorWebServer
}

type ProcessorWebServer struct {
	ListenPort int `validate:"required" split_words:"true"`
}
