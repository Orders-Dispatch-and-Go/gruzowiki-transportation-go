package app

import (
	"auth-service/internal/api/controller/carrierController"
	"auth-service/internal/api/middlewares"
	"auth-service/internal/config"
	"auth-service/internal/db/pg"
	"auth-service/internal/delivery"
	"auth-service/internal/repo/carrierRepo"
	"auth-service/internal/service/carrierService"
	"auth-service/internal/utils/auth"
	"auth-service/internal/utils/request"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sync"
)

func Run(ctx context.Context, cfg config.Config) error {
	conn, err := getDBConn(ctx, cfg.PostgresDSN)
	if err != nil {
		return fmt.Errorf("getDBConn(%q): %w", cfg.PostgresDSN, err)
	}

	defer conn.Close()

	auth := auth.NewJWT(cfg.JWT)
	requestReader := request.NewReader()

	carrierRepoV := carrierRepo.New(conn)
	carrierServiceV := carrierService.New(carrierRepoV)
	carrierControllerV := carrierController.New(requestReader, carrierServiceV)

	middlewaresV := middlewares.New(auth)
	deliveryV := delivery.New(cfg.Delivery, carrierControllerV, middlewaresV)

	if err := runDelivery(ctx, cfg, deliveryV); err != nil {
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
