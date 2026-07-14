package request

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type UpdateProductRequest struct {
	IdempotencyKey string `json:"idempotency_key"`
	Name           string `json:"name"`
	PriceCents     int64  `json:"price_cents"`
	Currency       string `json:"currency"`
	SellerId       string `json:"seller_id"`
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
		PriceCents:     req.PriceCents,
		Currency:       entities.Currency(req.Currency),
		SellerId:       sellerId,
	}, nil
}
