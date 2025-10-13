package config

import (
	"encoding/json"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"gruzowiki-transportation/internal/auth"
	"log"
	"os"
)

type Config struct {
	ServerPort  string
	PostgresDSN string
	Jwt         auth.JwtConfig
}

func Load(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("Open(%q): %w", path, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Close(): %v", err)
		}
	}()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("Decode(): %w", err)
	}

	if err := defaults.Set(&cfg); err != nil {
		return Config{}, fmt.Errorf("Set(%+v): %w", cfg, err)
	}

	if err := validator.New().Struct(cfg); err != nil {
		return Config{}, fmt.Errorf("Struct(%+v): %w", cfg, err)
	}

	return cfg, nil
}
