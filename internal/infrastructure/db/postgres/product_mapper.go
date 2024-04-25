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

func FromDBProduct(dbProduct *Product) *entities.Product {
	var seller = &entities.Seller{
		ID:   dbProduct.Seller.ID,
		Name: dbProduct.Seller.Name,
	}

	var p = &entities.Product{
		Name:   dbProduct.Name,
		Price:  dbProduct.Price,
		Seller: *seller,
	}
	p.ID = dbProduct.ID

	return p
}
