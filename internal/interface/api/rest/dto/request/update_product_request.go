package request

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
)

type UpdateProductRequest struct {
	IdempotencyKey string  `json:"idempotency_key"`
	Name           string  `json:"Name"`
	Price          float64 `json:"Price"`
	SellerId       string  `json:"SellerId"`
}

// ToUpdateProductCommand builds the command. The product Id comes from the
// URL path rather than the body.
func (req *UpdateProductRequest) ToUpdateProductCommand(id uuid.UUID) (*command.UpdateProductCommand, error) {
	sellerId, err := uuid.Parse(req.SellerId)
	if err != nil {
		return nil, err
	}

	return &command.UpdateProductCommand{
		IdempotencyKey: req.IdempotencyKey,
		Id:             id,
		Name:           req.Name,
		Price:          req.Price,
		SellerId:       sellerId,
	}, nil
}
