package config

type Config struct {
	AppPort string
	DBHost  string
}

func LoadConfig() *Config {
	return &Config{
		AppPort: "8080",
		DBHost:  "localhost",
	}
}
