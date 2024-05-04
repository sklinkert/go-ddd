package entities

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Product struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Price     float64
	Seller    Seller
}

func (p *Product) validate() error {
	if p.Name == "" {
		return errors.New("name must not be empty")
	}
	if p.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if p.CreatedAt.After(p.UpdatedAt) {
		return errors.New("created_at must be before updated_at")
	}

	return nil
}

func NewProduct(name string, price float64, seller ValidatedSeller) *Product {
	return &Product{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Price:     price,
		Seller:    seller.Seller,
	}
}

func (p *Product) UpdateName(name string) error {
	p.Name = name
	p.UpdatedAt = time.Now()

	return p.validate()
}

func (p *Product) UpdatePrice(price float64) error {
	p.Price = price
	p.UpdatedAt = time.Now()

	return p.validate()
}
