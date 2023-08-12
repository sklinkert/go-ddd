package postgres

import (
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

func ToDBProduct(product *entities.ValidatedProduct) *Product {
	var p = &Product{
		Name:     product.Name,
		Price:    product.Price,
		SellerID: product.Seller.ID, // Ensure Seller is non-nil when mapping
	}
	p.ID = product.ID

	return p
}

func FromDBProduct(dbProduct *Product) (*entities.ValidatedProduct, error) {
	var p = &entities.Product{
		Name:  dbProduct.Name,
		Price: dbProduct.Price,
		Seller: &entities.Seller{
			ID: dbProduct.SellerID,
		},
	}
	p.ID = dbProduct.ID

	return entities.NewValidatedProduct(p)
}
