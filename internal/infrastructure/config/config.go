package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config содержит всю конфигурацию приложения
type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	API      APIConfig      `mapstructure:"api"`
	Worker   WorkerConfig   `mapstructure:"worker"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

// DatabaseConfig содержит конфигурацию базы данных
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// APIConfig содержит конфигурацию API сервера
type APIConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// WorkerConfig содержит конфигурацию worker сервиса
type WorkerConfig struct {
	Interval        int    `mapstructure:"interval"`
	CoinGeckoAPIURL string `mapstructure:"coingecko_api_url"`
}

// LoggingConfig содержит конфигурацию логирования
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load загружает конфигурацию из файла и переменных окружения
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Устанавливаем значения по умолчанию
	setDefaults()

	// Устанавливаем алиасы для переменных окружения
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.ssl_mode", "DB_SSL_MODE")
	viper.BindEnv("api.port", "API_PORT")
	viper.BindEnv("api.host", "API_HOST")
	viper.BindEnv("worker.interval", "WORKER_INTERVAL")
	viper.BindEnv("worker.coingecko_api_url", "COINGECKO_API_URL")

	// Читаем переменные окружения
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Читаем конфигурационный файл (если существует)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults устанавливает значения по умолчанию
func setDefaults() {
	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.name", "crypto_tracker")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.ssl_mode", "disable")

	// API defaults
	viper.SetDefault("api.port", "8080")
	viper.SetDefault("api.host", "0.0.0.0")

	// Worker defaults
	viper.SetDefault("worker.interval", 60)
	viper.SetDefault("worker.coingecko_api_url", "https://api.coingecko.com/api/v3")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
}

// GetDSN возвращает строку подключения к базе данных
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}
