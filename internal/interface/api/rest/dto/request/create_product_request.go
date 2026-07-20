package request

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type CreateProductRequest struct {
	IdempotencyKey  string `json:"idempotency_key"`
	Name            string `json:"name"`
	PriceMinorUnits int64  `json:"price_minor_units"`
	Currency        string `json:"currency"`
	SellerId        string `json:"seller_id"`
}

func (req *CreateProductRequest) ToCreateProductCommand() (*command.CreateProductCommand, error) {
	sellerId, err := uuid.Parse(req.SellerId)
	if err != nil {
		return nil, err
	}

	return &command.CreateProductCommand{
		IdempotencyKey:  req.IdempotencyKey,
		Name:            req.Name,
		PriceMinorUnits: req.PriceMinorUnits,
		Currency:        entities.Currency(req.Currency),
		SellerId:        sellerId,
	}, nil
}
