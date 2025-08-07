package repository

import (
	"context"
	"time"
)

type CurrencyRepository interface {
	Create(ctx context.Context, currency interface{}) error
	GetBySymbol(ctx context.Context, symbol string) (interface{}, error)
	GetAllActive(ctx context.Context) ([]interface{}, error)
	Update(ctx context.Context, currency interface{}) error
	Delete(ctx context.Context, symbol string) error
	Deactivate(ctx context.Context, symbol string) error
}

type PriceRepository interface {
	Create(ctx context.Context, price interface{}) error
	GetByCurrencyAndTime(ctx context.Context, currencyID uint, timestamp time.Time) (interface{}, error)
	GetNearestPrice(ctx context.Context, currencyID uint, timestamp time.Time) (interface{}, error)
	GetLatestPrice(ctx context.Context, currencyID uint) (interface{}, error)
	GetPriceHistory(ctx context.Context, currencyID uint, from, to time.Time) ([]interface{}, error)
}

type PriceAPI interface {
	GetPrice(ctx context.Context, symbol string) (float64, error)
}
