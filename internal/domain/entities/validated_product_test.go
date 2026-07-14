package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestProductValidation(t *testing.T) {
	price, err := NewMoney(1000, USD)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}
	sellerId := uuid.Must(uuid.NewV7())

	// Test valid product
	validProduct := &Product{Name: "Valid Product", Price: price, SellerId: sellerId}
	if err := validProduct.validate(); err != nil {
		t.Errorf("Expected product to be valid, but got error: %s", err)
	}

	// Test product with empty name
	invalidProduct1 := &Product{Name: "", Price: price, SellerId: sellerId}
	if err := invalidProduct1.validate(); err == nil {
		t.Error("Expected product with empty name to be invalid, but got no error")
	}

	// Test product with zero price
	invalidProduct2 := &Product{Name: "Product", Price: Money{}, SellerId: sellerId}
	if err := invalidProduct2.validate(); err == nil {
		t.Error("Expected product with zero price to be invalid, but got no error")
	}

	// Test product without seller
	invalidProduct3 := &Product{Name: "Product", Price: price}
	if err := invalidProduct3.validate(); err == nil {
		t.Error("Expected product without seller id to be invalid, but got no error")
	}
}

func TestNewValidatedProduct(t *testing.T) {
	// Test valid product
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
	validatedProduct, err := NewValidatedProduct(product)
	if err != nil {
		t.Errorf("Expected product to be valid, but got error: %s", err)
	}
	if !validatedProduct.IsValid() {
		t.Error("Expected ValidatedProduct to be valid")
	}

	// Test invalid product (empty name, zero price)
	invalidProduct := NewProduct("", Money{}, *validatedSeller)
	validatedProduct, err = NewValidatedProduct(invalidProduct)
	if err == nil {
		t.Error("Expected error when validating invalid product, but got none")
	}
	if validatedProduct != nil {
		t.Error("Expected ValidatedProduct to be nil for invalid input")
	}
}
