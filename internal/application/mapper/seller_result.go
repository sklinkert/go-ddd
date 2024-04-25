package mapper

import (
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

func NewSellerResultFromEntity(seller entities.Seller) command.SellerResult {
	return command.SellerResult{
		ID:   seller.ID,
		Name: seller.Name,
	}
}
