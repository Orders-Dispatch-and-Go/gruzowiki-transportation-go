package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"auth-service/internal/api"
	"auth-service/internal/api/middlewares"
	"auth-service/internal/config"
	"auth-service/internal/db/pg"
	"auth-service/internal/delivery"
	"auth-service/internal/repo"
	"auth-service/internal/utils/auth"
	"auth-service/internal/utils/crypto"
	"auth-service/internal/utils/email"
	"auth-service/internal/utils/request"
)

func Run(ctx context.Context, cfg config.Config) error {
	conn, err := getDBConn(ctx, cfg.PostgresDSN)
	if err != nil {
		return fmt.Errorf("getDBConn(%q): %w", cfg.PostgresDSN, err)
	}

	defer conn.Close()

	//passwordHasher := crypto.NewBcryptPasswordHasher()
	auth := auth.NewJWT(cfg.JWT)
	//requestReader := request.NewReader()
	//emailSender := email.NewEmailSender(cfg.Email)

	//repo := repo.New(conn)
	//service := service.New(repo, passwordHasher, auth, emailSender)
	//api := api.New(requestReader, service)
	middlewares := middlewares.New(auth)
	delivery := delivery.New(cfg.Delivery, api, middlewares)

	if err := runDelivery(ctx, cfg, delivery); err != nil {
		return fmt.Errorf("run delivery: %w", err)
	}

	return nil
}

func runDelivery(ctx context.Context, cfg config.Config, delivery *delivery.Delivery) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := delivery.Listen(); err != nil {
			log.Printf("delivery listen failed: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), cfg.ShutdownTimeoutSeconds.Duration)
	defer cancelShutdownCtx()

	if err := delivery.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	return nil
}

func getDBConn(ctx context.Context, dsn string) (pg.Conn, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("New(%q): %w", dsn, err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Ping(): %w", err)
	}

	return pg.NewConn(pool), nil
}
