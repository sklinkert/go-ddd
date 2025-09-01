package interfaces

import (
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/query"
)

type SellerService interface {
	CreateSeller(sellerCommand *command.CreateSellerCommand) (*command.CreateSellerCommandResult, error)
	FindAllSellers() (*query.GetAllSellersQueryResult, error)
	FindSellerById(query *query.GetSellerByIdQuery) (*query.GetSellerByIdQueryResult, error)
	UpdateSeller(updateCommand *command.UpdateSellerCommand) (*command.UpdateSellerCommandResult, error)
	DeleteSeller(sellerCommand *command.DeleteSellerCommand) (*command.DeleteSellerCommandResult, error)
}
