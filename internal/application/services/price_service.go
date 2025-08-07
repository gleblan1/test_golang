package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"crypto-price-tracker-app/internal/domain/models"
	"crypto-price-tracker-app/internal/domain/repository"
)

// PriceService предоставляет бизнес-логику для работы с ценами
type PriceService struct {
	currencyRepo repository.CurrencyRepository
	priceRepo    repository.PriceRepository
	priceAPI     PriceAPI
}

// PriceAPI определяет интерфейс для получения цен из внешнего API
type PriceAPI interface {
	GetPrice(ctx context.Context, symbol string) (float64, error)
}

// NewPriceService создает новый экземпляр PriceService
func NewPriceService(currencyRepo repository.CurrencyRepository, priceRepo repository.PriceRepository, priceAPI PriceAPI) *PriceService {
	return &PriceService{
		currencyRepo: currencyRepo,
		priceRepo:    priceRepo,
		priceAPI:     priceAPI,
	}
}

// UpdatePrices обновляет цены для всех активных криптовалют
func (s *PriceService) UpdatePrices(ctx context.Context) error {
	// Получаем все активные криптовалюты
	currencies, err := s.currencyRepo.GetAllActive(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active currencies: %w", err)
	}

	// Обновляем цены для каждой криптовалюты
	for _, currency := range currencies {
		if err := s.updateCurrencyPrice(ctx, &currency); err != nil {
			log.Printf("Failed to update price for %s: %v", currency.Symbol, err)
			continue
		}
	}

	return nil
}

// updateCurrencyPrice обновляет цену для конкретной криптовалюты
func (s *PriceService) updateCurrencyPrice(ctx context.Context, currency *models.Currency) error {
	// Получаем текущую цену из внешнего API
	price, err := s.priceAPI.GetPrice(ctx, currency.Symbol)
	if err != nil {
		return fmt.Errorf("failed to get price from API: %w", err)
	}

	// Создаем новую запись о цене
	priceRecord := &models.Price{
		CurrencyID: currency.ID,
		Price:      price,
		Timestamp:  time.Now(),
	}

	// Сохраняем цену в базу данных
	if err := s.priceRepo.Create(ctx, priceRecord); err != nil {
		return fmt.Errorf("failed to save price to database: %w", err)
	}

	log.Printf("Updated price for %s: $%.2f", currency.Symbol, price)
	return nil
}

// GetLatestPrices возвращает последние цены для всех активных криптовалют
func (s *PriceService) GetLatestPrices(ctx context.Context) (map[string]float64, error) {
	currencies, err := s.currencyRepo.GetAllActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active currencies: %w", err)
	}

	prices := make(map[string]float64)
	for _, currency := range currencies {
		latestPrice, err := s.priceRepo.GetLatestPrice(ctx, currency.ID)
		if err != nil {
			log.Printf("Failed to get latest price for %s: %v", currency.Symbol, err)
			continue
		}
		prices[currency.Symbol] = latestPrice.Price
	}

	return prices, nil
}
