package postgres

import (
	"context"
	"errors"
	"time"

	"crypto-price-tracker-app/internal/domain/models"
	"crypto-price-tracker-app/internal/domain/repository"

	"gorm.io/gorm"
)

type PriceRepository struct {
	db *gorm.DB
}

func NewPriceRepository(db *gorm.DB) repository.PriceRepository {
	return &PriceRepository{db: db}
}

func (r *PriceRepository) Create(ctx context.Context, price interface{}) error {
	priceModel := price.(*models.Price)
	return r.db.WithContext(ctx).Create(priceModel).Error
}

func (r *PriceRepository) GetByCurrencyAndTime(ctx context.Context, currencyID uint, timestamp time.Time) (interface{}, error) {
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

func (r *PriceRepository) GetNearestPrice(ctx context.Context, currencyID uint, timestamp time.Time) (interface{}, error) {
	var price models.Price

	err := r.db.WithContext(ctx).
		Where("currency_id = ? AND timestamp <= ?", currencyID, timestamp).
		Order("timestamp DESC").
		First(&price).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

func (r *PriceRepository) GetLatestPrice(ctx context.Context, currencyID uint) (interface{}, error) {
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

func (r *PriceRepository) GetPriceHistory(ctx context.Context, currencyID uint, from, to time.Time) ([]interface{}, error) {
	var prices []models.Price
	err := r.db.WithContext(ctx).
		Where("currency_id = ? AND timestamp BETWEEN ? AND ?", currencyID, from, to).
		Order("timestamp ASC").
		Find(&prices).Error
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(prices))
	for i, price := range prices {
		result[i] = &price
	}
	return result, nil
}
