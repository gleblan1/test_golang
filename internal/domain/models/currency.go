package models

import (
	"time"

	"gorm.io/gorm"
)

type Currency struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Symbol    string         `json:"symbol" gorm:"uniqueIndex;not null"`
	ApiID     string         `json:"api_id" gorm:"not null"`
	Interval  int            `json:"interval" gorm:"not null;default:60"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type Price struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	CurrencyID uint           `json:"currency_id" gorm:"not null"`
	Price      float64        `json:"price" gorm:"not null"`
	Timestamp  time.Time      `json:"timestamp" gorm:"not null;index"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Currency   Currency       `json:"currency,omitempty" gorm:"foreignKey:CurrencyID"`
}
