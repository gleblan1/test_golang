package repository

import (
	"context"
	"time"

	"crypto-price-tracker-app/internal/domain/models"
)

// CurrencyRepository определяет интерфейс для работы с криптовалютами
type CurrencyRepository interface {
	// Create создает новую криптовалюту
	Create(ctx context.Context, currency *models.Currency) error

	// GetBySymbol возвращает криптовалюту по символу
	GetBySymbol(ctx context.Context, symbol string) (*models.Currency, error)

	// GetAllActive возвращает все активные криптовалюты
	GetAllActive(ctx context.Context) ([]models.Currency, error)

	// Update обновляет криптовалюту
	Update(ctx context.Context, currency *models.Currency) error

	// Delete удаляет криптовалюту по символу
	Delete(ctx context.Context, symbol string) error

	// Deactivate деактивирует криптовалюту
	Deactivate(ctx context.Context, symbol string) error
}

// PriceRepository определяет интерфейс для работы с ценами
type PriceRepository interface {
	// Create создает новую запись о цене
	Create(ctx context.Context, price *models.Price) error

	// GetByCurrencyAndTime возвращает цену для криптовалюты в указанное время
	GetByCurrencyAndTime(ctx context.Context, currencyID uint, timestamp time.Time) (*models.Price, error)

	// GetNearestPrice возвращает ближайшую доступную цену
	GetNearestPrice(ctx context.Context, currencyID uint, timestamp time.Time) (*models.Price, error)

	// GetLatestPrice возвращает последнюю цену для криптовалюты
	GetLatestPrice(ctx context.Context, currencyID uint) (*models.Price, error)

	// GetPriceHistory возвращает историю цен для криптовалюты
	GetPriceHistory(ctx context.Context, currencyID uint, from, to time.Time) ([]models.Price, error)
}
