package rest

import (
	"context"
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

func (m *MockSellerService) CreateSeller(ctx context.Context, seller *command.CreateSellerCommand) (*command.CreateSellerCommandResult, error) {
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

func (m *MockSellerService) FindAllSellers(ctx context.Context) (*query.GetAllSellersQueryResult, error) {
	var allSellers query.GetAllSellersQueryResult
	for _, v := range m.sellers {
		allSellers.Result = append(allSellers.Result, mapper.NewSellerResultFromEntity(&v.Seller))
	}
	return &allSellers, nil
}

func (m *MockSellerService) FindSellerById(ctx context.Context, sellerQuery *query.GetSellerByIdQuery) (*query.GetSellerByIdQueryResult, error) {
	if seller, exists := m.sellers[sellerQuery.Id]; exists {
		return &query.GetSellerByIdQueryResult{
			Result: mapper.NewSellerResultFromEntity(&seller.Seller),
		}, nil
	}
	return nil, errors.New("seller not found")
}

func (m *MockSellerService) UpdateSeller(ctx context.Context, updateCommand *command.UpdateSellerCommand) (*command.UpdateSellerCommandResult, error) {
	if _, exists := m.sellers[updateCommand.Id]; exists {
		m.sellers[updateCommand.Id].Name = updateCommand.Name
		return &command.UpdateSellerCommandResult{
			Result: mapper.NewSellerResultFromEntity(&m.sellers[updateCommand.Id].Seller),
		}, nil
	}
	return nil, errors.New("seller not found")
}

func (m *MockSellerService) DeleteSeller(ctx context.Context, deleteCommand *command.DeleteSellerCommand) (*command.DeleteSellerCommandResult, error) {
	if _, exists := m.sellers[deleteCommand.Id]; exists {
		delete(m.sellers, deleteCommand.Id)
		return &command.DeleteSellerCommandResult{Success: true}, nil
	}
	return nil, errors.New("seller not found")
}
