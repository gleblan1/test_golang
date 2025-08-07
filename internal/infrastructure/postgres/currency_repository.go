package postgres

import (
	"context"
	"errors"

	"crypto-price-tracker-app/internal/domain/models"
	"crypto-price-tracker-app/internal/domain/repository"

	"gorm.io/gorm"
)

type CurrencyRepository struct {
	db *gorm.DB
}

func NewCurrencyRepository(db *gorm.DB) repository.CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r *CurrencyRepository) Create(ctx context.Context, currency interface{}) error {
	currencyModel := currency.(*models.Currency)
	return r.db.WithContext(ctx).Create(currencyModel).Error
}

func (r *CurrencyRepository) GetBySymbol(ctx context.Context, symbol string) (interface{}, error) {
	var currency models.Currency
	err := r.db.WithContext(ctx).Where("symbol = ?", symbol).First(&currency).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &currency, nil
}

func (r *CurrencyRepository) GetAllActive(ctx context.Context) ([]interface{}, error) {
	var currencies []models.Currency
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&currencies).Error
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(currencies))
	for i, currency := range currencies {
		result[i] = &currency
	}
	return result, nil
}

func (r *CurrencyRepository) Update(ctx context.Context, currency interface{}) error {
	currencyModel := currency.(*models.Currency)
	return r.db.WithContext(ctx).Save(currencyModel).Error
}

func (r *CurrencyRepository) Delete(ctx context.Context, symbol string) error {
	return r.db.WithContext(ctx).Where("symbol = ?", symbol).Delete(&models.Currency{}).Error
}

func (r *CurrencyRepository) Deactivate(ctx context.Context, symbol string) error {
	return r.db.WithContext(ctx).Model(&models.Currency{}).Where("symbol = ?", symbol).Update("is_active", false).Error
}
