package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crypto-price-tracker-app/internal/application/services"
	"crypto-price-tracker-app/internal/infrastructure/coingecko"
	"crypto-price-tracker-app/internal/infrastructure/config"
	"crypto-price-tracker-app/internal/infrastructure/postgres"

	"go.uber.org/zap"
)

func main() {
	logger, err := initLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Crypto Price Tracker Worker")

	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		logger.Warn("Failed to load config file, using defaults", zap.Error(err))
		cfg = &config.Config{
			Database: config.DatabaseConfig{
				Host:     "postgres",
				Port:     5432,
				User:     "postgres",
				Password: "password",
				DBName:   "crypto_tracker",
				SSLMode:  "disable",
			},
			API: config.APIConfig{
				Port: "8080",
			},
			Worker: config.WorkerConfig{
				Interval: 60,
			},
			Logging: config.LoggingConfig{
				Level: "info",
			},
		}
	}

	db, err := postgres.NewConnection(&postgres.Config{
		Host:     cfg.Database.Host,
		Port:     fmt.Sprintf("%d", cfg.Database.Port),
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer postgres.CloseConnection(db)

	currencyRepo := postgres.NewCurrencyRepository(db)
	priceRepo := postgres.NewPriceRepository(db)

	coingeckoClient := coingecko.NewClient("https://api.coingecko.com/api/v3")

	priceService := services.NewPriceService(priceRepo, currencyRepo, coingeckoClient, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go runWorker(ctx, priceService, cfg.Worker.Interval, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker...")
	cancel()

	time.Sleep(5 * time.Second)
	logger.Info("Worker exited")
}

func initLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	return config.Build()
}

func runWorker(ctx context.Context, priceService *services.PriceService, interval int, logger *zap.Logger) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	logger.Info("Worker started", zap.Int("interval_seconds", interval))

	if err := priceService.UpdatePrices(ctx); err != nil {
		logger.Error("Failed to update prices", zap.Error(err))
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("Worker context cancelled")
			return
		case <-ticker.C:
			if err := priceService.UpdatePrices(ctx); err != nil {
				logger.Error("Failed to update prices", zap.Error(err))
			}
		}
	}
}
