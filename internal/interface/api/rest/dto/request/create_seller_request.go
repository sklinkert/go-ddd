package request

import "github.com/sklinkert/go-ddd/internal/application/command"

type CreateSellerRequest struct {
	IdempotencyKey string `json:"idempotency_key"`
	Name           string `json:"Name"`
}

func (req *CreateSellerRequest) ToCreateSellerCommand() (*command.CreateSellerCommand, error) {
	return &command.CreateSellerCommand{
		IdempotencyKey: req.IdempotencyKey,
		Name:           req.Name,
	}, nil
}
