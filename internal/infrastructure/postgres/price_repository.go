package postgres

import (
	"context"
	"errors"
	"time"

	"crypto-price-tracker-app/internal/domain/models"
	"crypto-price-tracker-app/internal/domain/repository"

	"gorm.io/gorm"
)

// PriceRepository реализует интерфейс repository.PriceRepository
type PriceRepository struct {
	db *gorm.DB
}

// NewPriceRepository создает новый экземпляр PriceRepository
func NewPriceRepository(db *gorm.DB) repository.PriceRepository {
	return &PriceRepository{db: db}
}

// Create создает новую запись о цене
func (r *PriceRepository) Create(ctx context.Context, price *models.Price) error {
	return r.db.WithContext(ctx).Create(price).Error
}

// GetByCurrencyAndTime возвращает цену для криптовалюты в указанное время
func (r *PriceRepository) GetByCurrencyAndTime(ctx context.Context, currencyID uint, timestamp time.Time) (*models.Price, error) {
	var price models.Price
	err := r.db.WithContext(ctx).
		Where("currency_id = ? AND timestamp = ?", currencyID, timestamp).
		First(&price).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &price, nil
}

// GetNearestPrice возвращает ближайшую доступную цену
func (r *PriceRepository) GetNearestPrice(ctx context.Context, currencyID uint, timestamp time.Time) (*models.Price, error) {
	var price models.Price

	// Ищем ближайшую цену до указанного времени
	err := r.db.WithContext(ctx).
		Where("currency_id = ? AND timestamp <= ?", currencyID, timestamp).
		Order("timestamp DESC").
		First(&price).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Если нет цен до указанного времени, ищем ближайшую после
			err = r.db.WithContext(ctx).
				Where("currency_id = ? AND timestamp >= ?", currencyID, timestamp).
				Order("timestamp ASC").
				First(&price).Error

			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, nil
				}
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &price, nil
}

// GetLatestPrice возвращает последнюю цену для криптовалюты
func (r *PriceRepository) GetLatestPrice(ctx context.Context, currencyID uint) (*models.Price, error) {
	var price models.Price
	err := r.db.WithContext(ctx).
		Where("currency_id = ?", currencyID).
		Order("timestamp DESC").
		First(&price).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &price, nil
}

// GetPriceHistory возвращает историю цен для криптовалюты
func (r *PriceRepository) GetPriceHistory(ctx context.Context, currencyID uint, from, to time.Time) ([]models.Price, error) {
	var prices []models.Price
	err := r.db.WithContext(ctx).
		Where("currency_id = ? AND timestamp BETWEEN ? AND ?", currencyID, from, to).
		Order("timestamp ASC").
		Find(&prices).Error
	return prices, err
}
