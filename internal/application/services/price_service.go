package services

import (
	"context"
	"time"

	"crypto-price-tracker-app/internal/domain/models"
	"crypto-price-tracker-app/internal/domain/repository"

	"go.uber.org/zap"
)

type PriceService struct {
	priceRepo    repository.PriceRepository
	currencyRepo repository.CurrencyRepository
	priceAPI     repository.PriceAPI
	logger       *zap.Logger
}

func NewPriceService(priceRepo repository.PriceRepository, currencyRepo repository.CurrencyRepository, priceAPI repository.PriceAPI, logger *zap.Logger) *PriceService {
	return &PriceService{
		priceRepo:    priceRepo,
		currencyRepo: currencyRepo,
		priceAPI:     priceAPI,
		logger:       logger,
	}
}

func (s *PriceService) UpdatePrices(ctx context.Context) error {
	currencies, err := s.currencyRepo.GetAllActive(ctx)
	if err != nil {
		s.logger.Error("Failed to get active currencies for price update", zap.Error(err))
		return err
	}

	s.logger.Info("Starting price update for currencies", zap.Int("count", len(currencies)))
	for _, currencyInterface := range currencies {
		currency := currencyInterface.(*models.Currency)
		if err := s.updateCurrencyPrice(ctx, currency); err != nil {
			s.logger.Error("Failed to update price for currency", zap.String("symbol", currency.Symbol), zap.Error(err))
		} else {
			s.logger.Debug("Price updated successfully", zap.String("symbol", currency.Symbol))
		}
	}

	s.logger.Info("Price update completed")
	return nil
}

func (s *PriceService) updateCurrencyPrice(ctx context.Context, currency *models.Currency) error {
	price, err := s.priceAPI.GetPrice(ctx, currency.ApiID)
	if err != nil {
		s.logger.Error("Failed to get price from API", zap.String("symbol", currency.Symbol), zap.String("api_id", currency.ApiID), zap.Error(err))
		return err
	}

	priceModel := &models.Price{
		CurrencyID: currency.ID,
		Price:      price,
		Timestamp:  time.Now(),
	}

	if err := s.priceRepo.Create(ctx, priceModel); err != nil {
		s.logger.Error("Failed to save price to database", zap.String("symbol", currency.Symbol), zap.Float64("price", price), zap.Error(err))
		return err
	}

	s.logger.Debug("Price saved successfully", zap.String("symbol", currency.Symbol), zap.Float64("price", price))
	return nil
}

func (s *PriceService) GetLatestPrices(ctx context.Context) ([]models.Price, error) {
	currencies, err := s.currencyRepo.GetAllActive(ctx)
	if err != nil {
		s.logger.Error("Failed to get active currencies for latest prices", zap.Error(err))
		return nil, err
	}

	var prices []models.Price
	for _, currencyInterface := range currencies {
		currency := currencyInterface.(*models.Currency)
		priceInterface, err := s.priceRepo.GetLatestPrice(ctx, currency.ID)
		if err != nil || priceInterface == nil {
			s.logger.Debug("No latest price found for currency", zap.String("symbol", currency.Symbol))
			continue
		}
		price := priceInterface.(*models.Price)
		prices = append(prices, *price)
	}

	s.logger.Debug("Retrieved latest prices", zap.Int("count", len(prices)))
	return prices, nil
}
