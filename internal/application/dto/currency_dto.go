package dto

import "time"

type AddCurrencyRequest struct {
	Symbol   string `json:"symbol" binding:"required" example:"BTC"`
	ApiID    string `json:"api_id" binding:"required" example:"bitcoin"`
	Interval int    `json:"interval" binding:"required,min=30" example:"60"`
}

type RemoveCurrencyRequest struct {
	Symbol string `json:"symbol" binding:"required" example:"bitcoin"`
}

type GetPriceRequest struct {
	Coin      string `form:"coin" binding:"required" example:"bitcoin"`
	Timestamp int64  `form:"timestamp" binding:"required" example:"1640995200"`
}

type CurrencyResponse struct {
	ID        uint      `json:"id"`
	Symbol    string    `json:"symbol"`
	ApiID     string    `json:"api_id"`
	Interval  int       `json:"interval"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PriceResponse struct {
	ID        uint      `json:"id"`
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
