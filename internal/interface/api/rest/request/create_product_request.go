package request

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
)

type CreateProductRequest struct {
	Name     string  `json:"Name"`
	Price    float64 `json:"Price"`
	SellerID string  `json:"SellerId"`
}

func (req *CreateProductRequest) ToCreateProductCommand() (*command.CreateProductCommand, error) {
	sellerId, err := uuid.Parse(req.SellerID)
	if err != nil {
		return nil, err
	}

	return &command.CreateProductCommand{
		Name:     req.Name,
		Price:    req.Price,
		SellerID: sellerId,
	}, nil
}
