package internal

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	Address        string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8080"`
	DBURI          string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Key            string `env:"SECRET_KEY" envDefault:"asdqwezxc123"`
	LogFile        string `env:"LOG_FILE" envDefault:"logs.json"`
	Debug          bool   `env:"DEBUG" envDefault:"false"`
}

func NewConfig() Config {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}
	return cfg
}
