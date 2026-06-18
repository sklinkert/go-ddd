package entities

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewSeller(t *testing.T) {
	seller := NewSeller("Example Seller")

	if seller.Name != "Example Seller" {
		t.Errorf("Expected seller name to be 'Example Seller', but got %s", seller.Name)
	}

	if seller.Id == (uuid.UUID{}) {
		t.Error("Expected seller Id to be set, but got zero value")
	}
}
