package command

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/common"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type CreateProductCommand struct {
	IdempotencyKey string
	Id             uuid.UUID
	Name           string
	PriceCents     int64
	Currency       entities.Currency
	SellerId       uuid.UUID
}

type CreateProductCommandResult struct {
	Result *common.ProductResult
}
