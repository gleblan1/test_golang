package dto

import "time"

// AddCurrencyRequest представляет запрос на добавление криптовалюты
type AddCurrencyRequest struct {
	Symbol   string `json:"symbol" binding:"required" example:"bitcoin"`
	Interval int    `json:"interval" binding:"required,min=30" example:"60"`
}

// RemoveCurrencyRequest представляет запрос на удаление криптовалюты
type RemoveCurrencyRequest struct {
	Symbol string `json:"symbol" binding:"required" example:"bitcoin"`
}

// GetPriceRequest представляет запрос на получение цены
type GetPriceRequest struct {
	Coin      string `form:"coin" binding:"required" example:"bitcoin"`
	Timestamp int64  `form:"timestamp" binding:"required" example:"1640995200"`
}

// CurrencyResponse представляет ответ с информацией о криптовалюте
type CurrencyResponse struct {
	ID        uint      `json:"id"`
	Symbol    string    `json:"symbol"`
	Interval  int       `json:"interval"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PriceResponse представляет ответ с информацией о цене
type PriceResponse struct {
	ID        uint      `json:"id"`
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SuccessResponse представляет успешный ответ
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
