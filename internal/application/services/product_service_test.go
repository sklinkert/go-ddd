package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"testing"
)

// MockProductRepository is a mock implementation of the ProductRepository interface
type MockProductRepository struct {
	products []*entities.ValidatedProduct
}

func (m *MockProductRepository) Save(product *entities.ValidatedProduct) error {
	m.products = append(m.products, product)
	return nil
}

func (m *MockProductRepository) GetAll() ([]*entities.ValidatedProduct, error) {
	return m.products, nil
}

func (m *MockProductRepository) Update(product *entities.ValidatedProduct) error {
	for index, p := range m.products {
		if p.ID == product.ID {
			m.products[index] = product
			return nil
		}
	}
	return errors.New("product not found for update")
}

func (m *MockProductRepository) Delete(id uuid.UUID) error {
	for index, p := range m.products {
		if p.ID == id {
			m.products = append(m.products[:index], m.products[index+1:]...)
			return nil
		}
	}
	return errors.New("product not found for delete")
}

func (m *MockProductRepository) FindByID(id uuid.UUID) (*entities.ValidatedProduct, error) {
	for _, p := range m.products {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, errors.New("product not found")
}

func TestProductService_CreateProduct(t *testing.T) {
	repo := &MockProductRepository{}
	service := NewProductService(repo)

	product := entities.NewProduct("Example", 100.0)
	err := service.CreateProduct(product)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if len(repo.products) != 1 {
		t.Errorf("Expected 1 product in repo, but got %d", len(repo.products))
	}
}

func TestProductService_GetAllProducts(t *testing.T) {
	repo := &MockProductRepository{}
	service := NewProductService(repo)

	// Add two products
	_ = service.CreateProduct(entities.NewProduct("Example1", 100.0))
	_ = service.CreateProduct(entities.NewProduct("Example2", 200.0))

	products, err := service.GetAllProducts()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if len(products) != 2 {
		t.Errorf("Expected 2 products, but got %d", len(products))
	}
}

func TestProductService_FindProductByID(t *testing.T) {
	repo := &MockProductRepository{}
	service := NewProductService(repo)

	product := entities.NewProduct("Example", 100.0)
	_ = service.CreateProduct(product)

	foundProduct, err := service.FindProductByID(product.ID)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if foundProduct.Name != "Example" {
		t.Errorf("Expected product name 'Example', but got %s", foundProduct.Name)
	}

	_, err = service.FindProductByID(uuid.New()) // some non-existent ID
	if err == nil {
		t.Error("Expected error for non-existent product, but got none")
	}
}
