package services

import (
	"context"
	"errors"
	"time"

	"crypto-price-tracker-app/internal/application/dto"
	"crypto-price-tracker-app/internal/domain/models"
	"crypto-price-tracker-app/internal/domain/repository"

	"go.uber.org/zap"
)

type CurrencyService struct {
	currencyRepo repository.CurrencyRepository
	priceRepo    repository.PriceRepository
	logger       *zap.Logger
}

func NewCurrencyService(currencyRepo repository.CurrencyRepository, priceRepo repository.PriceRepository, logger *zap.Logger) *CurrencyService {
	return &CurrencyService{
		currencyRepo: currencyRepo,
		priceRepo:    priceRepo,
		logger:       logger,
	}
}

func (s *CurrencyService) AddCurrency(ctx context.Context, req *dto.AddCurrencyRequest) (*dto.CurrencyResponse, error) {
	existing, err := s.currencyRepo.GetBySymbol(ctx, req.Symbol)
	if err == nil && existing != nil {
		s.logger.Warn("Currency already exists", zap.String("symbol", req.Symbol))
		return nil, errors.New("currency already exists")
	}

	currency := &models.Currency{
		Symbol:   req.Symbol,
		ApiID:    req.ApiID,
		Interval: req.Interval,
		IsActive: true,
	}

	if err := s.currencyRepo.Create(ctx, currency); err != nil {
		s.logger.Error("Failed to create currency", zap.String("symbol", req.Symbol), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Currency added successfully", zap.String("symbol", req.Symbol), zap.Uint("id", currency.ID))
	return &dto.CurrencyResponse{
		ID:        currency.ID,
		Symbol:    currency.Symbol,
		ApiID:     currency.ApiID,
		Interval:  currency.Interval,
		IsActive:  currency.IsActive,
		CreatedAt: currency.CreatedAt,
		UpdatedAt: currency.UpdatedAt,
	}, nil
}

func (s *CurrencyService) RemoveCurrency(ctx context.Context, req *dto.RemoveCurrencyRequest) error {
	_, err := s.currencyRepo.GetBySymbol(ctx, req.Symbol)
	if err != nil {
		s.logger.Warn("Currency not found for removal", zap.String("symbol", req.Symbol))
		return errors.New("currency not found")
	}

	if err := s.currencyRepo.Deactivate(ctx, req.Symbol); err != nil {
		s.logger.Error("Failed to deactivate currency", zap.String("symbol", req.Symbol), zap.Error(err))
		return err
	}

	s.logger.Info("Currency removed successfully", zap.String("symbol", req.Symbol))
	return nil
}

func (s *CurrencyService) GetPrice(ctx context.Context, req *dto.GetPriceRequest) (*dto.PriceResponse, error) {
	currencyInterface, err := s.currencyRepo.GetBySymbol(ctx, req.Coin)
	if err != nil {
		s.logger.Error("Failed to get currency", zap.String("symbol", req.Coin), zap.Error(err))
		return nil, errors.New("currency not found")
	}

	if currencyInterface == nil {
		s.logger.Warn("Currency not found", zap.String("symbol", req.Coin))
		return nil, errors.New("currency not found")
	}

	currency := currencyInterface.(*models.Currency)

	timestamp := time.Unix(req.Timestamp, 0)

	priceInterface, err := s.priceRepo.GetByCurrencyAndTime(ctx, currency.ID, timestamp)
	if err != nil {
		s.logger.Error("Failed to get price by time", zap.String("symbol", req.Coin), zap.Time("timestamp", timestamp), zap.Error(err))
		return nil, errors.New("price not found")
	}

	if priceInterface == nil {
		s.logger.Debug("Exact price not found, searching for nearest", zap.String("symbol", req.Coin), zap.Time("timestamp", timestamp))
		priceInterface, err = s.priceRepo.GetNearestPrice(ctx, currency.ID, timestamp)
		if err != nil {
			s.logger.Error("Failed to get nearest price", zap.String("symbol", req.Coin), zap.Time("timestamp", timestamp), zap.Error(err))
			return nil, errors.New("price not found")
		}
		if priceInterface == nil {
			s.logger.Warn("No price found for currency", zap.String("symbol", req.Coin), zap.Time("timestamp", timestamp))
			return nil, errors.New("price not found")
		}
	}

	price := priceInterface.(*models.Price)

	s.logger.Debug("Price retrieved successfully", zap.String("symbol", req.Coin), zap.Float64("price", price.Price), zap.Time("timestamp", price.Timestamp))
	return &dto.PriceResponse{
		ID:        price.ID,
		Symbol:    currency.Symbol,
		Price:     price.Price,
		Timestamp: price.Timestamp,
		CreatedAt: price.CreatedAt,
	}, nil
}

func (s *CurrencyService) GetAllActiveCurrencies(ctx context.Context) ([]dto.CurrencyResponse, error) {
	currenciesInterface, err := s.currencyRepo.GetAllActive(ctx)
	if err != nil {
		s.logger.Error("Failed to get active currencies", zap.Error(err))
		return nil, err
	}

	var responses []dto.CurrencyResponse
	for _, currencyInterface := range currenciesInterface {
		currency := currencyInterface.(*models.Currency)
		responses = append(responses, dto.CurrencyResponse{
			ID:        currency.ID,
			Symbol:    currency.Symbol,
			ApiID:     currency.ApiID,
			Interval:  currency.Interval,
			IsActive:  currency.IsActive,
			CreatedAt: currency.CreatedAt,
			UpdatedAt: currency.UpdatedAt,
		})
	}

	s.logger.Debug("Retrieved active currencies", zap.Int("count", len(responses)))
	return responses, nil
}
