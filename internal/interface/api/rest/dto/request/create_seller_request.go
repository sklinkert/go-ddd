package request

import "github.com/sklinkert/go-ddd/internal/application/command"

type CreateSellerRequest struct {
	Name string `json:"Name"`
}

func (req *CreateSellerRequest) ToCreateSellerCommand() (*command.CreateSellerCommand, error) {
	return &command.CreateSellerCommand{
		Name: req.Name,
	}, nil
}
