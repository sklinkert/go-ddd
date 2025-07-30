package request

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
)

type UpdateSellerRequest struct {
	IdempotencyKey string    `json:"idempotency_key"`
	Id             uuid.UUID `json:"Id"`
	Name           string    `json:"Name"`
}

func (req *UpdateSellerRequest) ToUpdateSellerCommand() (*command.UpdateSellerCommand, error) {
	return &command.UpdateSellerCommand{
		IdempotencyKey: req.IdempotencyKey,
		Id:             req.Id,
		Name:           req.Name,
	}, nil
}
