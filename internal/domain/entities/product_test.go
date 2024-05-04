package entities

import (
	"github.com/google/uuid"
	"testing"
)

func TestNewProduct(t *testing.T) {
	seller := NewSeller("Example Seller")
	validatedSeller, err := NewValidatedSeller(seller)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}

	product := NewProduct("Example Product", 10.0, *validatedSeller)

	if product.Name != "Example Product" {
		t.Errorf("Expected product name to be 'Example Product', but got %s", product.Name)
	}

	if product.Price != 10.0 {
		t.Errorf("Expected product price to be 10.0, but got %f", product.Price)
	}

	if product.Id == (uuid.UUID{}) {
		t.Error("Expected product Id to be set, but got zero value")
	}
}
