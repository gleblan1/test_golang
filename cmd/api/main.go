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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title Crypto Price Tracker API
// @version 1.0
// @description Микросервис для отслеживания цен криптовалют
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Инициализируем логгер
	logger, err := initLogger()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Crypto Price Tracker API")

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

	// Инициализируем сервисы
	currencyService := services.NewCurrencyService(currencyRepo, priceRepo)
	_ = services.NewPriceService(currencyRepo, priceRepo, coingeckoClient)

	// Инициализируем HTTP handlers
	currencyHandler := handlers.NewCurrencyHandler(currencyService)

	// Настраиваем Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Добавляем middleware
	router.Use(middleware.LoggerMiddleware(logger))
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.CORSMiddleware())

	// Настраиваем маршруты
	setupRoutes(router, currencyHandler)

	// Создаем HTTP сервер
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.API.Host, cfg.API.Port),
		Handler: router,
	}

	// Запускаем сервер в горутине
	go func() {
		logger.Info("Starting HTTP server", zap.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Ждем сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

// initLogger инициализирует логгер
func initLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	return config.Build()
}

// setupRoutes настраивает маршруты API
func setupRoutes(router *gin.Engine, currencyHandler *handlers.CurrencyHandler) {
	// API v1
	v1 := router.Group("/api/v1")
	{
		// Currency endpoints
		currency := v1.Group("/currency")
		{
			currency.POST("/add", currencyHandler.AddCurrency)
			currency.POST("/remove", currencyHandler.RemoveCurrency)
			currency.GET("/price", currencyHandler.GetPrice)
			currency.GET("/list", currencyHandler.GetAllCurrencies)
		}
	}

	// Health check
	router.GET("/health", currencyHandler.HealthCheck)

	// Swagger documentation
	router.Static("/docs", "./docs")
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
