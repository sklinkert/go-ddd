package entities

import (
	"testing"
)

func TestSellerValidation(t *testing.T) {
	// Test valid seller
	validSeller := &Seller{Name: "Valid Seller"}
	if err := validSeller.validate(); err != nil {
		t.Errorf("Expected seller to be valid, but got error: %s", err)
	}

	// Test seller with empty name
	invalidSeller := &Seller{Name: ""}
	if err := invalidSeller.validate(); err == nil {
		t.Error("Expected seller with empty name to be invalid, but got no error")
	}
}

func TestNewValidatedSeller(t *testing.T) {
	// Test valid seller
	seller := NewSeller("Example Seller")
	validatedSeller, err := NewValidatedSeller(seller)
	if err != nil {
		t.Errorf("Expected seller to be valid, but got error: %s", err)
	}
	if !validatedSeller.IsValid() {
		t.Error("Expected ValidatedSeller to be valid")
	}

	// Test invalid seller
	invalidSeller := NewSeller("")
	validatedSeller, err = NewValidatedSeller(invalidSeller)
	if err == nil {
		t.Error("Expected error when validating invalid seller, but got none")
	}
	if validatedSeller != nil {
		t.Error("Expected ValidatedSeller to be nil for invalid input")
	}
}
