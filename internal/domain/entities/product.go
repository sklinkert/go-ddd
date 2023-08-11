package entities

import (
	"errors"
	"github.com/google/uuid"
)

type Product struct {
	ID     uuid.UUID
	Name   string
	Price  float64
	Seller *Seller
}

type ValidatedProduct struct {
	Product
	isValidated bool
}

func (vp *ValidatedProduct) IsValid() bool {
	return vp.isValidated
}

func (p *Product) validate() error {
	if p.Name == "" || p.Price <= 0 {
		return errors.New("invalid product details")
	}

	return nil
}

func NewProduct(name string, price float64) *Product {
	return &Product{
		ID:    uuid.New(),
		Name:  name,
		Price: price,
	}
}

func NewValidatedProduct(product *Product) (*ValidatedProduct, error) {
	if err := product.validate(); err != nil {
		return nil, err
	}

	return &ValidatedProduct{
		Product:     *product,
		isValidated: true,
	}, nil
}
