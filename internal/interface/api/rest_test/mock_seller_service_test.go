package rest

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/interfaces"
	"github.com/sklinkert/go-ddd/internal/application/mapper"
	"github.com/sklinkert/go-ddd/internal/application/query"
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

	m.sellers[validatedSeller.Id] = validatedSeller

	result.Result = mapper.NewSellerResultFromEntity(&validatedSeller.Seller)

	return &result, nil
}

func (m *MockSellerService) FindAllSellers() (*query.SellerQueryListResult, error) {
	var allSellers query.SellerQueryListResult
	for _, v := range m.sellers {
		allSellers.Result = append(allSellers.Result, mapper.NewSellerResultFromEntity(&v.Seller))
	}
	return &allSellers, nil
}

func (m *MockSellerService) FindSellerById(id uuid.UUID) (*query.SellerQueryResult, error) {
	if seller, exists := m.sellers[id]; exists {
		return &query.SellerQueryResult{
			Result: mapper.NewSellerResultFromEntity(&seller.Seller),
		}, nil
	}
	return nil, errors.New("seller not found")
}

func (m *MockSellerService) UpdateSeller(updateCommand *command.UpdateSellerCommand) (*command.UpdateSellerCommandResult, error) {
	if _, exists := m.sellers[updateCommand.Id]; exists {
		m.sellers[updateCommand.Id].Name = updateCommand.Name
		return &command.UpdateSellerCommandResult{
			Result: mapper.NewSellerResultFromEntity(&m.sellers[updateCommand.Id].Seller),
		}, nil
	}
	return nil, errors.New("seller not found")
}

func (m *MockSellerService) DeleteSeller(id uuid.UUID) error {
	if _, exists := m.sellers[id]; exists {
		delete(m.sellers, id)
		return nil
	}
	return errors.New("seller not found")
}
