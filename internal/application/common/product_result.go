package common

import (
	"github.com/google/uuid"
	"time"
)

type ProductResult struct {
	Id        uuid.UUID
	Name      string
	Price     float64
	Seller    *SellerResult
	CreatedAt time.Time
	UpdatedAt time.Time
}
