package rest_test

import (
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/mapper"
	"github.com/sklinkert/go-ddd/internal/application/query"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(productCommand *command.CreateProductCommand) (*command.CreateProductCommandResult, error) {
	args := m.Called(productCommand)

	var now = time.Now()

	var seller = &entities.Seller{
		Id:        productCommand.SellerId,
		Name:      "Test Seller",
		CreatedAt: now,
		UpdatedAt: now,
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
	result.Result = mapper.NewProductResultFromValidatedEntity(validatedProduct)

	return &result, args.Error(1)
}

func (m *MockProductService) FindAllProducts() (*query.GetAllProductsQueryResult, error) {
	args := m.Called()

	productQueryListResult := &query.GetAllProductsQueryResult{}

	for _, product := range args.Get(0).([]*entities.Product) {
		productQueryListResult.Result = append(productQueryListResult.Result, mapper.NewProductResultFromEntity(product))
	}

	return productQueryListResult, args.Error(1)
}

func (m *MockProductService) FindProductById(productQuery *query.GetProductByIdQuery) (*query.GetProductByIdQueryResult, error) {
	args := m.Called(productQuery)

	productQueryResult := &query.GetProductByIdQueryResult{
		Result: mapper.NewProductResultFromEntity(args.Get(0).(*entities.Product)),
	}

	return productQueryResult, args.Error(1)
}

func (m *MockProductService) UpdateProduct(productCommand *command.UpdateProductCommand) (*command.UpdateProductCommandResult, error) {
	args := m.Called(productCommand)
	return args.Get(0).(*command.UpdateProductCommandResult), args.Error(1)
}

func (m *MockProductService) DeleteProduct(productCommand *command.DeleteProductCommand) (*command.DeleteProductCommandResult, error) {
	args := m.Called(productCommand)
	return args.Get(0).(*command.DeleteProductCommandResult), args.Error(1)
}
