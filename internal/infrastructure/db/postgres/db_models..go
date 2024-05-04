package postgres

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	Id        uuid.UUID `gorm:"primaryKey"`
	Name      string
	Price     float64
	SellerID  uuid.UUID `gorm:"index"`
	Seller    Seller    `gorm:"foreignKey:SellerID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Seller struct {
	Id        uuid.UUID `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
