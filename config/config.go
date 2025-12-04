package config

import (
	"flag"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ClinetUrl string `yaml:"client_url" env:"CLIENT_URL" env-default:"http://localhost:8082"`
	LogLevel  string `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	Address   string `yaml:"transportation_service_address" env:"TRANSPORTATION_SERVICE_ADDRESS" env-default:"localhost:8081"`
	Dsn       string `yaml:"dsn" env:"DSN" env-default:"user='postgres' password='postgres' host='localhost' port=5433 dbname='postgres' sslmode=disable"`
}

func MustLoad() Config {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("cannot read config %q: %s", configPath, err)
		}
	}
	return cfg
}
