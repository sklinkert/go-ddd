package interfaces

import (
	"context"

	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/query"
)

type SellerService interface {
	CreateSeller(ctx context.Context, sellerCommand *command.CreateSellerCommand) (*command.CreateSellerCommandResult, error)
	FindAllSellers(ctx context.Context) (*query.GetAllSellersQueryResult, error)
	FindSellerById(ctx context.Context, query *query.GetSellerByIdQuery) (*query.GetSellerByIdQueryResult, error)
	UpdateSeller(ctx context.Context, updateCommand *command.UpdateSellerCommand) (*command.UpdateSellerCommandResult, error)
	DeleteSeller(ctx context.Context, sellerCommand *command.DeleteSellerCommand) (*command.DeleteSellerCommandResult, error)
}
