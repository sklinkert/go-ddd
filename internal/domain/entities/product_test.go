package entities

import (
	"github.com/google/uuid"
	"testing"
)

func TestNewProduct(t *testing.T) {
	seller := NewSeller("Example Seller")
	product := NewProduct("Example Product", 10.0, seller)

	if product.Name != "Example Product" {
		t.Errorf("Expected product name to be 'Example Product', but got %s", product.Name)
	}

	if product.Price != 10.0 {
		t.Errorf("Expected product price to be 10.0, but got %f", product.Price)
	}

	if product.ID == (uuid.UUID{}) {
		t.Error("Expected product ID to be set, but got zero value")
	}
}
