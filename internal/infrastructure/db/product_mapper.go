package db

import (
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

func ToDBProduct(product *entities.Product) *Product {
	var p = &Product{
		Name:     product.Name,
		Price:    product.Price,
		SellerID: product.Seller.ID, // Ensure Seller is non-nil when mapping
	}
	p.ID = product.ID

	return p
}

func FromDBProduct(dbProduct *Product) *entities.Product {
	var p = &entities.Product{
		Name:  dbProduct.Name,
		Price: dbProduct.Price,
		Seller: &entities.Seller{
			ID: dbProduct.SellerID,
		},
	}
	p.ID = dbProduct.ID

	return p
}
