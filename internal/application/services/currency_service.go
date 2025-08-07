package services

import (
	"context"
	"errors"
	"time"

	"crypto-price-tracker-app/internal/application/dto"
	"crypto-price-tracker-app/internal/domain/models"
	"crypto-price-tracker-app/internal/domain/repository"
)

// CurrencyService предоставляет бизнес-логику для работы с криптовалютами
type CurrencyService struct {
	currencyRepo repository.CurrencyRepository
	priceRepo    repository.PriceRepository
}

// NewCurrencyService создает новый экземпляр CurrencyService
func NewCurrencyService(currencyRepo repository.CurrencyRepository, priceRepo repository.PriceRepository) *CurrencyService {
	return &CurrencyService{
		currencyRepo: currencyRepo,
		priceRepo:    priceRepo,
	}
}

// AddCurrency добавляет новую криптовалюту для отслеживания
func (s *CurrencyService) AddCurrency(ctx context.Context, req *dto.AddCurrencyRequest) (*dto.CurrencyResponse, error) {
	// Проверяем, существует ли уже такая криптовалюта
	existing, err := s.currencyRepo.GetBySymbol(ctx, req.Symbol)
	if err == nil && existing != nil {
		return nil, errors.New("currency already exists")
	}

	// Создаем новую криптовалюту
	currency := &models.Currency{
		Symbol:   req.Symbol,
		Interval: req.Interval,
		IsActive: true,
	}

	if err := s.currencyRepo.Create(ctx, currency); err != nil {
		return nil, err
	}

	return &dto.CurrencyResponse{
		ID:        currency.ID,
		Symbol:    currency.Symbol,
		Interval:  currency.Interval,
		IsActive:  currency.IsActive,
		CreatedAt: currency.CreatedAt,
		UpdatedAt: currency.UpdatedAt,
	}, nil
}

// RemoveCurrency удаляет криптовалюту из отслеживания
func (s *CurrencyService) RemoveCurrency(ctx context.Context, req *dto.RemoveCurrencyRequest) error {
	// Проверяем, существует ли криптовалюта
	_, err := s.currencyRepo.GetBySymbol(ctx, req.Symbol)
	if err != nil {
		return errors.New("currency not found")
	}

	// Деактивируем криптовалюту
	return s.currencyRepo.Deactivate(ctx, req.Symbol)
}

// GetPrice возвращает цену криптовалюты в указанное время
func (s *CurrencyService) GetPrice(ctx context.Context, req *dto.GetPriceRequest) (*dto.PriceResponse, error) {
	// Получаем криптовалюту по символу
	currency, err := s.currencyRepo.GetBySymbol(ctx, req.Coin)
	if err != nil {
		return nil, errors.New("currency not found")
	}
	
	// Проверяем, что криптовалюта найдена
	if currency == nil {
		return nil, errors.New("currency not found")
	}

	// Преобразуем timestamp в time.Time
	timestamp := time.Unix(req.Timestamp, 0)

	// Пытаемся найти точную цену
	price, err := s.priceRepo.GetByCurrencyAndTime(ctx, currency.ID, timestamp)
	if err != nil {
		return nil, errors.New("price not found")
	}
	
	// Если точной цены нет, ищем ближайшую
	if price == nil {
		price, err = s.priceRepo.GetNearestPrice(ctx, currency.ID, timestamp)
		if err != nil {
			return nil, errors.New("price not found")
		}
		if price == nil {
			return nil, errors.New("price not found")
		}
	}

	return &dto.PriceResponse{
		ID:        price.ID,
		Symbol:    currency.Symbol,
		Price:     price.Price,
		Timestamp: price.Timestamp,
		CreatedAt: price.CreatedAt,
	}, nil
}

// GetAllActiveCurrencies возвращает все активные криптовалюты
func (s *CurrencyService) GetAllActiveCurrencies(ctx context.Context) ([]dto.CurrencyResponse, error) {
	currencies, err := s.currencyRepo.GetAllActive(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.CurrencyResponse
	for _, currency := range currencies {
		responses = append(responses, dto.CurrencyResponse{
			ID:        currency.ID,
			Symbol:    currency.Symbol,
			Interval:  currency.Interval,
			IsActive:  currency.IsActive,
			CreatedAt: currency.CreatedAt,
			UpdatedAt: currency.UpdatedAt,
		})
	}

	return responses, nil
}
