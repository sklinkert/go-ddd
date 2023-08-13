package entities

import (
	"testing"
)

func TestProductValidation(t *testing.T) {
	// Test valid product
	validProduct := &Product{Name: "Valid Product", Price: 10.0}
	if err := validProduct.validate(); err != nil {
		t.Errorf("Expected product to be valid, but got error: %s", err)
	}

	// Test product with empty name
	invalidProduct1 := &Product{Name: "", Price: 10.0}
	if err := invalidProduct1.validate(); err == nil {
		t.Error("Expected product with empty name to be invalid, but got no error")
	}

	// Test product with non-positive price
	invalidProduct2 := &Product{Name: "Product", Price: -5.0}
	if err := invalidProduct2.validate(); err == nil {
		t.Error("Expected product with negative price to be invalid, but got no error")
	}
}

func TestNewValidatedProduct(t *testing.T) {
	// Test valid product
	seller := NewSeller("Example Seller")
	validatedSeller, err := NewValidatedSeller(seller)
	if err != nil {
		t.Fatalf("Expected no error, but got %s", err.Error())
	}
	product := NewProduct("Example Product", 10.0, *validatedSeller)
	validatedProduct, err := NewValidatedProduct(product)
	if err != nil {
		t.Errorf("Expected product to be valid, but got error: %s", err)
	}
	if !validatedProduct.IsValid() {
		t.Error("Expected ValidatedProduct to be valid")
	}

	// Test invalid product
	invalidProduct := NewProduct("", -10.0, *validatedSeller)
	validatedProduct, err = NewValidatedProduct(invalidProduct)
	if err == nil {
		t.Error("Expected error when validating invalid product, but got none")
	}
	if validatedProduct != nil {
		t.Error("Expected ValidatedProduct to be nil for invalid input")
	}
}
