package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"testing"
)

// MockSellerRepository is a mock implementation of the SellerRepository interface
type MockSellerRepository struct {
	sellers []*entities.ValidatedSeller
}

func (m *MockSellerRepository) Create(seller *entities.ValidatedSeller) error {
	m.sellers = append(m.sellers, seller)
	return nil
}

func (m *MockSellerRepository) GetAll() ([]*entities.ValidatedSeller, error) {
	return m.sellers, nil
}

func (m *MockSellerRepository) FindByID(id uuid.UUID) (*entities.ValidatedSeller, error) {
	for _, s := range m.sellers {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, errors.New("seller not found")
}

func (m *MockSellerRepository) Delete(id uuid.UUID) error {
	for index, s := range m.sellers {
		if s.ID == id {
			m.sellers = append(m.sellers[:index], m.sellers[index+1:]...)
			return nil
		}
	}
	return errors.New("seller not found for deletion")
}

func (m *MockSellerRepository) Update(seller *entities.ValidatedSeller) error {
	for index, s := range m.sellers {
		if s.ID == seller.ID {
			m.sellers[index] = seller
			return nil
		}
	}
	return errors.New("seller not found for update")
}

func TestSellerService_CreateSeller(t *testing.T) {
	repo := &MockSellerRepository{}
	service := NewSellerService(repo)

	seller := entities.NewSeller("John Doe")
	err := service.CreateSeller(seller)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if len(repo.sellers) != 1 {
		t.Errorf("Expected 1 seller in repo, but got %d", len(repo.sellers))
	}
}

func TestSellerService_GetAllSellers(t *testing.T) {
	repo := &MockSellerRepository{}
	service := NewSellerService(repo)

	// Add two sellers
	_ = service.CreateSeller(entities.NewSeller("John Doe"))
	_ = service.CreateSeller(entities.NewSeller("Jane Doe"))

	sellers, err := service.GetAllSellers()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if len(sellers) != 2 {
		t.Errorf("Expected 2 sellers, but got %d", len(sellers))
	}
}

func TestSellerService_GetSellerByID(t *testing.T) {
	repo := &MockSellerRepository{}
	service := NewSellerService(repo)

	seller := entities.NewSeller("John Doe")
	_ = service.CreateSeller(seller)

	foundSeller, err := service.GetSellerByID(seller.ID)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if foundSeller.Name != "John Doe" {
		t.Errorf("Expected seller name 'John Doe', but got %s", foundSeller.Name)
	}

	_, err = service.GetSellerByID(uuid.New()) // some non-existent ID
	if err == nil {
		t.Error("Expected error for non-existent seller, but got none")
	}
}

func TestSellerService_UpdateSeller(t *testing.T) {
	repo := &MockSellerRepository{}
	service := NewSellerService(repo)

	seller := entities.NewSeller("John Doe")
	_ = service.CreateSeller(seller)
	seller.Name = "Johnny Doe"
	err := service.UpdateSeller(&entities.ValidatedSeller{Seller: *seller})

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	updatedSeller, _ := service.GetSellerByID(seller.ID)
	if updatedSeller.Name != "Johnny Doe" {
		t.Errorf("Expected seller name 'Johnny Doe', but got %s", updatedSeller.Name)
	}
}
