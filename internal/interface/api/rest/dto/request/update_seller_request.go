package request

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
)

type UpdateSellerRequest struct {
	ID   uuid.UUID `json:"Id"`
	Name string    `json:"Name"`
}

func (req *UpdateSellerRequest) ToUpdateSellerCommand() (*command.UpdateSellerCommand, error) {
	return &command.UpdateSellerCommand{
		ID:   req.ID,
		Name: req.Name,
	}, nil
}
