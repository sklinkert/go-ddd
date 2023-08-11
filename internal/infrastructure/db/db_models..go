package db

import (
	"github.com/google/uuid"
)

type Product struct {
	ID       uuid.UUID
	Name     string
	Price    float64
	SellerID uuid.UUID
}

type Seller struct {
	ID   uuid.UUID
	Name string
}
