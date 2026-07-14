package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewProduct(t *testing.T) {
	seller := NewSeller("Example Seller")
	validatedSeller, err := NewValidatedSeller(seller)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	price, err := NewMoney(1000, USD)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	product := NewProduct("Example Product", price, *validatedSeller)

	if product.Name != "Example Product" {
		t.Errorf("Expected product name to be 'Example Product', but got %s", product.Name)
	}

	if product.Price != price {
		t.Errorf("Expected product price to be %s, but got %s", price, product.Price)
	}

	if product.SellerId != validatedSeller.Id {
		t.Errorf("Expected product seller id to be %s, but got %s", validatedSeller.Id, product.SellerId)
	}

	if product.Id == (uuid.UUID{}) {
		t.Error("Expected product Id to be set, but got zero value")
	}
}
