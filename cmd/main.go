package main

import (
	"auth-service/internal/app"
	"auth-service/internal/config"
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
)

const configPath = "config.json"

func run() error {
	ctx, stopCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopCtx()

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("Load(%q): %w", configPath, err)
	}

	if err := app.Run(ctx, cfg); err != nil {
		return fmt.Errorf("Run(%+v): %w", cfg, err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
