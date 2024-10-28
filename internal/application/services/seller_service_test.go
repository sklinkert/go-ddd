package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"testing"
)

// MockSellerRepository is a mock implementation of the SellerRepository interface
type MockSellerRepository struct {
	sellers []*entities.ValidatedSeller
}

func (m *MockSellerRepository) Create(seller *entities.ValidatedSeller) (*entities.Seller, error) {
	m.sellers = append(m.sellers, seller)
	return &seller.Seller, nil
}

func (m *MockSellerRepository) FindAll() ([]*entities.Seller, error) {
	var sellers []*entities.Seller
	for _, s := range m.sellers {
		sellers = append(sellers, &s.Seller)
	}
	return sellers, nil
}

func (m *MockSellerRepository) FindById(id uuid.UUID) (*entities.Seller, error) {
	for _, s := range m.sellers {
		if s.Id == id {
			return &s.Seller, nil
		} else {
			fmt.Printf("Id: %s - %s\n", s.Id, id)
		}
	}
	return nil, errors.New("seller not found")
}

func (m *MockSellerRepository) Delete(id uuid.UUID) error {
	for index, s := range m.sellers {
		if s.Id == id {
			m.sellers = append(m.sellers[:index], m.sellers[index+1:]...)
			return nil
		}
	}
	return errors.New("seller not found for deletion")
}

func (m *MockSellerRepository) Update(seller *entities.ValidatedSeller) (*entities.Seller, error) {
	for index, s := range m.sellers {
		if s.Id == seller.Id {
			m.sellers[index] = seller
			return &seller.Seller, nil
		}
	}
	return nil, errors.New("seller not found for update")
}

func TestSellerService_CreateSeller(t *testing.T) {
	repo := &MockSellerRepository{}
	service := NewSellerService(repo)

	_, err := service.CreateSeller(getCreateSellerCommand("John Doe"))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if len(repo.sellers) != 1 {
		t.Errorf("Expected 1 seller in productRepository, but got %d", len(repo.sellers))
	}
}

func TestSellerService_GetAllSellers(t *testing.T) {
	repo := &MockSellerRepository{}
	service := NewSellerService(repo)

	// Add two sellers
	_, _ = service.CreateSeller(getCreateSellerCommand("John Doe"))
	_, _ = service.CreateSeller(getCreateSellerCommand("Jane Doe"))

	sellers, err := service.FindAllSellers()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if len(sellers.Result) != 2 {
		t.Errorf("Expected 2 sellers, but got %d", len(sellers.Result))
	}
}

func TestSellerService_GetSellerById(t *testing.T) {
	repo := &MockSellerRepository{}
	service := NewSellerService(repo)

	createdSellerResult, _ := service.CreateSeller(getCreateSellerCommand("John Doe"))
	sellerID := createdSellerResult.Result.Id

	foundSeller, err := service.FindSellerById(sellerID)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if foundSeller.Result.Name != "John Doe" {
		t.Errorf("Expected seller name 'John Doe', but got %s", foundSeller.Result.Name)
	}

	_, err = service.FindSellerById(uuid.New()) // some non-existent Id
	if err == nil {
		t.Error("Expected error for non-existent seller, but got none")
	}
}

func TestSellerService_UpdateSeller(t *testing.T) {
	repo := &MockSellerRepository{}
	service := NewSellerService(repo)

	createdSellerResult, _ := service.CreateSeller(getCreateSellerCommand("John Doe"))
	sellerId := createdSellerResult.Result.Id

	var updatableSeller = entities.Seller{
		Id:   sellerId,
		Name: "Doe Johnny",
	}

	_, err := service.UpdateSeller(&command.UpdateSellerCommand{
		Id:   sellerId,
		Name: updatableSeller.Name,
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	updatedSeller, _ := service.FindSellerById(sellerId)
	if updatedSeller.Result.Name != "Doe Johnny" {
		t.Errorf("Expected seller name 'Johnny Doe', but got %s", updatedSeller.Result.Name)
	}
}

func getCreateSellerCommand(name string) *command.CreateSellerCommand {
	return &command.CreateSellerCommand{
		Name: name,
	}
}
