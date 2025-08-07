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
	// Инициализируем логгер
	logger, err := initLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Crypto Price Tracker Worker")

	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Подключаемся к базе данных
	db, err := postgres.NewConnection(&postgres.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.Name,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer postgres.CloseConnection(db)

	// Инициализируем репозитории
	currencyRepo := postgres.NewCurrencyRepository(db)
	priceRepo := postgres.NewPriceRepository(db)

	// Инициализируем внешний API клиент
	coingeckoClient := coingecko.NewClient(cfg.Worker.CoinGeckoAPIURL)

	// Инициализируем сервис обновления цен
	priceService := services.NewPriceService(currencyRepo, priceRepo, coingeckoClient)

	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем worker в горутине
	go runWorker(ctx, priceService, cfg.Worker.Interval, logger)

	// Ждем сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker...")
	cancel()

	// Даем время для завершения текущих операций
	time.Sleep(5 * time.Second)
	logger.Info("Worker exited")
}

// initLogger инициализирует логгер
func initLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	return config.Build()
}

// runWorker запускает основной цикл worker'а
func runWorker(ctx context.Context, priceService *services.PriceService, interval int, logger *zap.Logger) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	logger.Info("Worker started", zap.Int("interval_seconds", interval))

	// Выполняем первое обновление сразу
	if err := priceService.UpdatePrices(ctx); err != nil {
		logger.Error("Failed to update prices", zap.Error(err))
	}

	// Основной цикл обновления цен
	for {
		select {
		case <-ctx.Done():
			logger.Info("Worker context cancelled")
			return
		case <-ticker.C:
			if err := priceService.UpdatePrices(ctx); err != nil {
				logger.Error("Failed to update prices", zap.Error(err))
			} else {
				logger.Info("Prices updated successfully")
			}
		}
	}
}
