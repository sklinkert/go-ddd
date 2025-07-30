package request

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
)

type CreateProductRequest struct {
	IdempotencyKey string  `json:"idempotency_key"`
	Name           string  `json:"Name"`
	Price          float64 `json:"Price"`
	SellerId       string  `json:"SellerId"`
}

func (req *CreateProductRequest) ToCreateProductCommand() (*command.CreateProductCommand, error) {
	sellerId, err := uuid.Parse(req.SellerId)
	if err != nil {
		return nil, err
	}

	return &command.CreateProductCommand{
		IdempotencyKey: req.IdempotencyKey,
		Name:           req.Name,
		Price:          req.Price,
		SellerId:       sellerId,
	}, nil
}
