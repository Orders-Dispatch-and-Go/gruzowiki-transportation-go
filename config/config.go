package config

import (
	"flag"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	Address  string `yaml:"update_address" env:"UPDATE_ADDRESS" env-default:"localhost:80"`
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