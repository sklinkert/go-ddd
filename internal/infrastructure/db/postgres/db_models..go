package postgres

import (
	"github.com/google/uuid"
)

type Product struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	Name     string
	Price    float64
	SellerID uuid.UUID `gorm:"index"`
	Seller   Seller    `gorm:"foreignKey:SellerID"`
}

type Seller struct {
	ID   uuid.UUID `gorm:"primaryKey"`
	Name string
}
