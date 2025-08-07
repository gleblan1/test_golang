package models

import (
	"time"
)

// Currency представляет криптовалюту для отслеживания
type Currency struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Symbol    string    `json:"symbol" gorm:"uniqueIndex;not null"`
	Interval  int       `json:"interval" gorm:"not null"` // интервал обновления в секундах
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Price представляет цену криптовалюты в определенный момент времени
type Price struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CurrencyID uint      `json:"currency_id" gorm:"not null"`
	Currency   Currency  `json:"currency" gorm:"foreignKey:CurrencyID"`
	Price      float64   `json:"price" gorm:"not null"`
	Timestamp  time.Time `json:"timestamp" gorm:"not null;index"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName возвращает имя таблицы для модели Currency
func (Currency) TableName() string {
	return "currencies"
}

// TableName возвращает имя таблицы для модели Price
func (Price) TableName() string {
	return "prices"
}
