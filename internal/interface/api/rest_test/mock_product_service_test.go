package rest_test

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(product *entities.Product) (*entities.ValidatedProduct, error) {
	args := m.Called(product)
	return args.Get(0).(*entities.ValidatedProduct), args.Error(1)
}

func (m *MockProductService) GetAllProducts() ([]*entities.ValidatedProduct, error) {
	args := m.Called()
	return args.Get(0).([]*entities.ValidatedProduct), args.Error(1)
}

func (m *MockProductService) FindProductByID(id uuid.UUID) (*entities.ValidatedProduct, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.ValidatedProduct), args.Error(1)
}
