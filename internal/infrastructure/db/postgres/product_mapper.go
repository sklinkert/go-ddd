package postgres

import (
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

func toDBProduct(product *entities.ValidatedProduct) *Product {
	var p = &Product{
		Name:      product.Name,
		Price:     product.Price,
		SellerId:  product.Seller.Id, // Ensure Seller is non-nil when mapping
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
	p.Id = product.Id

	return p
}

func fromDBProduct(dbProduct *Product) *entities.Product {
	var seller = &entities.Seller{
		Id:        dbProduct.Seller.Id,
		Name:      dbProduct.Seller.Name,
		CreatedAt: dbProduct.Seller.CreatedAt,
		UpdatedAt: dbProduct.Seller.UpdatedAt,
	}

	var p = &entities.Product{
		Name:      dbProduct.Name,
		Price:     dbProduct.Price,
		Seller:    *seller,
		CreatedAt: dbProduct.CreatedAt,
		UpdatedAt: dbProduct.UpdatedAt,
	}
	p.Id = dbProduct.Id

	return p
}
