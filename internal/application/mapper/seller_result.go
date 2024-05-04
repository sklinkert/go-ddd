package mapper

import (
	"github.com/sklinkert/go-ddd/internal/application/common"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

func NewSellerResultFromValidatedEntity(seller *entities.ValidatedSeller) *common.SellerResult {
	return NewSellerResultFromEntity(&seller.Seller)
}

func NewSellerResultFromEntity(seller *entities.Seller) *common.SellerResult {
	if seller == nil {
		return nil
	}

	return &common.SellerResult{
		Id:        seller.ID,
		Name:      seller.Name,
		CreatedAt: seller.CreatedAt,
		UpdatedAt: seller.UpdatedAt,
	}
}
