package mapper

import (
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

func NewProductResultFromEntity(product *entities.ValidatedProduct) command.ProductResult {
	return command.ProductResult{
		Id:     product.ID,
		Name:   product.Name,
		Price:  product.Price,
		Seller: NewSellerResultFromEntity(product.Seller),
	}
}
