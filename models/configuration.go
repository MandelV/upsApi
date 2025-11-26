package models

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Configuration struct {
	VERSION     string `env:"VERSION" envDefault:"0.0.1"`
	ENV         string `env:"ENV" envDefault:"dev"`
	LOG_LEVEL   int    `env:"LOG_LEVEL" envDefault:"0"` // 0: debug, 1: info, 2: warn, 3: error, 4: fatal
	LISTEN_PORT int    `env:"LISTEN_PORT" envDefault:"9695"`
	UPS_USER    string `env:"UPS_USER" envDefault:"eaton"`
	UPS_HOST    string `env:"UPS_HOST" envDefault:"localhost"`
}

func NewConfiguration() *Configuration {

	godotenv.Load(".env")

	conf := &Configuration{}

	if err := env.Parse(conf); err != nil {
		panic(err)
	}
	return conf
}
