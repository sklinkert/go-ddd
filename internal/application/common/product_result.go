package common

import (
	"time"

	"github.com/google/uuid"
)

type ProductResult struct {
	Id        uuid.UUID
	Name      string
	Price     float64
	Seller    *SellerResult
	CreatedAt time.Time
	UpdatedAt time.Time
}
