package rest

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/interfaces"
	"github.com/sklinkert/go-ddd/internal/application/mapper"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type MockSellerService struct {
	sellers map[uuid.UUID]*entities.ValidatedSeller
}

func NewMockSellerService() interfaces.SellerService {
	return &MockSellerService{
		sellers: make(map[uuid.UUID]*entities.ValidatedSeller),
	}
}

func (m *MockSellerService) CreateSeller(seller *command.CreateSellerCommand) (*command.CreateSellerCommandResult, error) {
	var result command.CreateSellerCommandResult

	newSeller := entities.NewSeller(seller.Name)

	validatedSeller, err := entities.NewValidatedSeller(newSeller)
	if err != nil {
		return nil, err
	}

	m.sellers[validatedSeller.ID] = validatedSeller

	result.Result = mapper.NewSellerResultFromEntity(*validatedSeller)

	return &result, nil
}

func (m *MockSellerService) GetAllSellers() ([]*entities.ValidatedSeller, error) {
	var allSellers []*entities.ValidatedSeller
	for _, v := range m.sellers {
		allSellers = append(allSellers, v)
	}
	return allSellers, nil
}

func (m *MockSellerService) GetSellerByID(id uuid.UUID) (*entities.ValidatedSeller, error) {
	if seller, exists := m.sellers[id]; exists {
		return seller, nil
	}
	return nil, errors.New("seller not found")
}

func (m *MockSellerService) UpdateSeller(seller *entities.ValidatedSeller) error {
	if _, exists := m.sellers[seller.ID]; exists {
		m.sellers[seller.ID] = seller
		return nil
	}
	return errors.New("seller not found")
}
