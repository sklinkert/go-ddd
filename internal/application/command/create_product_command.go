package command

import "github.com/google/uuid"

type CreateProductCommand struct {
	// TODO: Implement idempotency key

	ID       uuid.UUID
	Name     string
	Price    float64
	SellerID uuid.UUID
}
