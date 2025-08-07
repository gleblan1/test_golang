package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"crypto-price-tracker-app/internal/application/dto"
	"crypto-price-tracker-app/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCurrencyRepository - мок для репозитория криптовалют
type MockCurrencyRepository struct {
	mock.Mock
}

func (m *MockCurrencyRepository) Create(ctx context.Context, currency *models.Currency) error {
	args := m.Called(ctx, currency)
	return args.Error(0)
}

func (m *MockCurrencyRepository) GetBySymbol(ctx context.Context, symbol string) (*models.Currency, error) {
	args := m.Called(ctx, symbol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Currency), args.Error(1)
}

func (m *MockCurrencyRepository) GetAllActive(ctx context.Context) ([]models.Currency, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Currency), args.Error(1)
}

func (m *MockCurrencyRepository) Update(ctx context.Context, currency *models.Currency) error {
	args := m.Called(ctx, currency)
	return args.Error(0)
}

func (m *MockCurrencyRepository) Delete(ctx context.Context, symbol string) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

func (m *MockCurrencyRepository) Deactivate(ctx context.Context, symbol string) error {
	args := m.Called(ctx, symbol)
	return args.Error(0)
}

// MockPriceRepository - мок для репозитория цен
type MockPriceRepository struct {
	mock.Mock
}

func (m *MockPriceRepository) Create(ctx context.Context, price *models.Price) error {
	args := m.Called(ctx, price)
	return args.Error(0)
}

func (m *MockPriceRepository) GetByCurrencyAndTime(ctx context.Context, currencyID uint, timestamp time.Time) (*models.Price, error) {
	args := m.Called(ctx, currencyID, timestamp)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Price), args.Error(1)
}

func (m *MockPriceRepository) GetNearestPrice(ctx context.Context, currencyID uint, timestamp time.Time) (*models.Price, error) {
	args := m.Called(ctx, currencyID, timestamp)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Price), args.Error(1)
}

func (m *MockPriceRepository) GetLatestPrice(ctx context.Context, currencyID uint) (*models.Price, error) {
	args := m.Called(ctx, currencyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Price), args.Error(1)
}

func (m *MockPriceRepository) GetPriceHistory(ctx context.Context, currencyID uint, from, to time.Time) ([]models.Price, error) {
	args := m.Called(ctx, currencyID, from, to)
	return args.Get(0).([]models.Price), args.Error(1)
}

func TestCurrencyService_AddCurrency(t *testing.T) {
	tests := []struct {
		name          string
		req           *dto.AddCurrencyRequest
		setupMocks    func(*MockCurrencyRepository)
		expectedError string
	}{
		{
			name: "successful add currency",
			req: &dto.AddCurrencyRequest{
				Symbol:   "bitcoin",
				Interval: 60,
			},
			setupMocks: func(mockRepo *MockCurrencyRepository) {
				mockRepo.On("GetBySymbol", mock.Anything, "bitcoin").Return(nil, errors.New("not found"))
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Currency")).Return(nil)
			},
		},
		{
			name: "currency already exists",
			req: &dto.AddCurrencyRequest{
				Symbol:   "bitcoin",
				Interval: 60,
			},
			setupMocks: func(mockRepo *MockCurrencyRepository) {
				existingCurrency := &models.Currency{
					ID:       1,
					Symbol:   "bitcoin",
					Interval: 60,
					IsActive: true,
				}
				mockRepo.On("GetBySymbol", mock.Anything, "bitcoin").Return(existingCurrency, nil)
			},
			expectedError: "currency already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCurrencyRepo := new(MockCurrencyRepository)
			mockPriceRepo := new(MockPriceRepository)

			if tt.setupMocks != nil {
				tt.setupMocks(mockCurrencyRepo)
			}

			service := NewCurrencyService(mockCurrencyRepo, mockPriceRepo)
			result, err := service.AddCurrency(context.Background(), tt.req)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.req.Symbol, result.Symbol)
				assert.Equal(t, tt.req.Interval, result.Interval)
			}

			mockCurrencyRepo.AssertExpectations(t)
		})
	}
}

func TestCurrencyService_RemoveCurrency(t *testing.T) {
	tests := []struct {
		name          string
		req           *dto.RemoveCurrencyRequest
		setupMocks    func(*MockCurrencyRepository)
		expectedError string
	}{
		{
			name: "successful remove currency",
			req: &dto.RemoveCurrencyRequest{
				Symbol: "bitcoin",
			},
			setupMocks: func(mockRepo *MockCurrencyRepository) {
				existingCurrency := &models.Currency{
					ID:       1,
					Symbol:   "bitcoin",
					Interval: 60,
					IsActive: true,
				}
				mockRepo.On("GetBySymbol", mock.Anything, "bitcoin").Return(existingCurrency, nil)
				mockRepo.On("Deactivate", mock.Anything, "bitcoin").Return(nil)
			},
		},
		{
			name: "currency not found",
			req: &dto.RemoveCurrencyRequest{
				Symbol: "nonexistent",
			},
			setupMocks: func(mockRepo *MockCurrencyRepository) {
				mockRepo.On("GetBySymbol", mock.Anything, "nonexistent").Return(nil, errors.New("not found"))
			},
			expectedError: "currency not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCurrencyRepo := new(MockCurrencyRepository)
			mockPriceRepo := new(MockPriceRepository)

			if tt.setupMocks != nil {
				tt.setupMocks(mockCurrencyRepo)
			}

			service := NewCurrencyService(mockCurrencyRepo, mockPriceRepo)
			err := service.RemoveCurrency(context.Background(), tt.req)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockCurrencyRepo.AssertExpectations(t)
		})
	}
}
