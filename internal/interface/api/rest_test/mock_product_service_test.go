package rest_test

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/mapper"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(productCommand *command.CreateProductCommand) (*command.CreateProductCommandResult, error) {
	args := m.Called(productCommand)

	var seller = &entities.Seller{
		ID:   productCommand.SellerID,
		Name: "Test Seller",
	}

	var validatedSeller, err = entities.NewValidatedSeller(seller)
	if err != nil {
		return nil, err
	}

	var newProduct = entities.NewProduct(
		productCommand.Name,
		productCommand.Price,
		*validatedSeller,
	)

	validatedProduct, err := entities.NewValidatedProduct(newProduct)
	if err != nil {
		return nil, err
	}

	var result command.CreateProductCommandResult
	result.Result = mapper.NewProductResultFromEntity(validatedProduct)

	return &result, args.Error(1)
}

func (m *MockProductService) FindAllProducts() ([]*entities.Product, error) {
	args := m.Called()
	return args.Get(0).([]*entities.Product), args.Error(1)
}

func (m *MockProductService) FindProductById(id uuid.UUID) (*entities.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Product), args.Error(1)
}
