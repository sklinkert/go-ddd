package postgres

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	Id        uuid.UUID `gorm:"primaryKey"`
	Name      string
	Price     float64
	SellerId  uuid.UUID `gorm:"index"`
	Seller    Seller    `gorm:"foreignKey:SellerId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Seller struct {
	Id        uuid.UUID `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
