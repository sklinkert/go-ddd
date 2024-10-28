package command

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/common"
)

type CreateProductCommand struct {
	// TODO: Implement idempotency key

	ID       uuid.UUID
	Name     string
	Price    float64
	SellerID uuid.UUID
}

type CreateProductCommandResult struct {
	Result *common.ProductResult
}
