package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crypto-price-tracker-app/internal/application/services"
	handlers "crypto-price-tracker-app/internal/delivery/http"
	"crypto-price-tracker-app/internal/delivery/middleware"
	"crypto-price-tracker-app/internal/infrastructure/coingecko"
	"crypto-price-tracker-app/internal/infrastructure/config"
	"crypto-price-tracker-app/internal/infrastructure/postgres"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, err := initLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Crypto Price Tracker API")

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

	currencyService := services.NewCurrencyService(currencyRepo, priceRepo, logger)
	priceService := services.NewPriceService(priceRepo, currencyRepo, coingeckoClient, logger)

	handlers := handlers.NewHandlers(currencyService, priceService)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORSMiddleware())

	setupRoutes(router, handlers)

	server := &http.Server{
		Addr:    ":" + cfg.API.Port,
		Handler: router,
	}

	go func() {
		logger.Info("Starting HTTP server", zap.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func initLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	return config.Build()
}

func setupRoutes(router *gin.Engine, handlers *handlers.Handlers) {
	v1 := router.Group("/api/v1")
	{
		currency := v1.Group("/currency")
		{
			currency.POST("/add", handlers.AddCurrency)
			currency.POST("/remove", handlers.RemoveCurrency)
			currency.GET("/price", handlers.GetPrice)
			currency.GET("/list", handlers.GetAllCurrencies)
		}
	}

	router.GET("/health", handlers.HealthCheck)

	router.Static("/docs", "./docs")
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})
}
