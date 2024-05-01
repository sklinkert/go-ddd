package common

import "github.com/google/uuid"

type ProductResult struct {
	Id     uuid.UUID
	Name   string
	Price  float64
	Seller *SellerResult
}
