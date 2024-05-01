package mapper

import (
	"github.com/sklinkert/go-ddd/internal/application/common"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

func NewSellerResultFromEntity(seller entities.Seller) common.SellerResult {
	return common.SellerResult{
		ID:   seller.ID,
		Name: seller.Name,
	}
}
