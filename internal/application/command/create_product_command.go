package command

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/common"
)

type CreateProductCommand struct {
	// TODO: Implement idempotency key

	Id       uuid.UUID
	Name     string
	Price    float64
	SellerId uuid.UUID
}

type CreateProductCommandResult struct {
	Result *common.ProductResult
}
