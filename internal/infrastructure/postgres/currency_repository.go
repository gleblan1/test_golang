package postgres

import (
	"context"
	"errors"

	"crypto-price-tracker-app/internal/domain/models"
	"crypto-price-tracker-app/internal/domain/repository"

	"gorm.io/gorm"
)

// CurrencyRepository реализует интерфейс repository.CurrencyRepository
type CurrencyRepository struct {
	db *gorm.DB
}

// NewCurrencyRepository создает новый экземпляр CurrencyRepository
func NewCurrencyRepository(db *gorm.DB) repository.CurrencyRepository {
	return &CurrencyRepository{db: db}
}

// Create создает новую криптовалюту
func (r *CurrencyRepository) Create(ctx context.Context, currency *models.Currency) error {
	return r.db.WithContext(ctx).Create(currency).Error
}

// GetBySymbol возвращает криптовалюту по символу
func (r *CurrencyRepository) GetBySymbol(ctx context.Context, symbol string) (*models.Currency, error) {
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

// GetAllActive возвращает все активные криптовалюты
func (r *CurrencyRepository) GetAllActive(ctx context.Context) ([]models.Currency, error) {
	var currencies []models.Currency
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&currencies).Error
	return currencies, err
}

// Update обновляет криптовалюту
func (r *CurrencyRepository) Update(ctx context.Context, currency *models.Currency) error {
	return r.db.WithContext(ctx).Save(currency).Error
}

// Delete удаляет криптовалюту по символу
func (r *CurrencyRepository) Delete(ctx context.Context, symbol string) error {
	return r.db.WithContext(ctx).Where("symbol = ?", symbol).Delete(&models.Currency{}).Error
}

// Deactivate деактивирует криптовалюту
func (r *CurrencyRepository) Deactivate(ctx context.Context, symbol string) error {
	return r.db.WithContext(ctx).Model(&models.Currency{}).Where("symbol = ?", symbol).Update("is_active", false).Error
}
