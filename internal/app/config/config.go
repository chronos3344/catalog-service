package main

func (c *Config) Load() error {
	_ = godotenv.load()
	c.LoadFromEnv()
}

type Config struct{}

func main() {}
