package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"testing"
)

// MockProductRepository is a mock implementation of the ProductRepository interface
type MockProductRepository struct {
	products []*entities.ValidatedProduct
}

func (m *MockProductRepository) Create(product *entities.ValidatedProduct) error {
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
		fmt.Printf("ID: mem:%s - %s\n", p.ID, id)
	}
	return nil, errors.New("product not found")
}

func TestProductService_CreateProduct(t *testing.T) {
	productRepo := &MockProductRepository{}
	sellerRepo := &MockSellerRepository{}
	service := NewProductService(productRepo, sellerRepo)

	// Create seller
	seller := createPersistedSeller(t, sellerRepo)

	// Create product
	product := entities.NewProduct("Example", 100.0, *seller)
	productCommand := getCreateProductCommand(product)
	_, err := service.CreateProduct(productCommand)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if len(productRepo.products) != 1 {
		t.Errorf("Expected 1 product in productRepository, but got %d", len(productRepo.products))
	}
}

func TestProductService_GetAllProducts(t *testing.T) {
	productRepo := &MockProductRepository{}
	sellerRepo := &MockSellerRepository{}
	service := NewProductService(productRepo, sellerRepo)

	// Create seller
	seller := createPersistedSeller(t, sellerRepo)

	// Add two products
	_, _ = service.CreateProduct(getCreateProductCommand(entities.NewProduct("Example1", 100.0, *seller)))
	_, _ = service.CreateProduct(getCreateProductCommand(entities.NewProduct("Example2", 200.0, *seller)))

	products, err := service.GetAllProducts()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if len(products) != 2 {
		t.Errorf("Expected 2 products, but got %d", len(products))
	}
}

func TestProductService_FindProductByID(t *testing.T) {
	productRepo := &MockProductRepository{}
	sellerRepo := &MockSellerRepository{}
	service := NewProductService(productRepo, sellerRepo)

	// Create seller
	seller := createPersistedSeller(t, sellerRepo)

	product := entities.NewProduct("Example", 100.0, *seller)
	result, err := service.CreateProduct(getCreateProductCommand(product))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	foundProduct, err := service.FindProductByID(result.Result.Id)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if foundProduct.Name != "Example" {
		t.Errorf("Expected product name 'Example', but got %s", foundProduct.Name)
	}

	_, err = service.FindProductByID(uuid.New()) // some non-existent ID
	if err == nil {
		t.Error("Expected error for non-existent product, but got none")
	}
}

func getCreateProductCommand(product *entities.Product) *command.CreateProductCommand {
	return &command.CreateProductCommand{
		Name:     product.Name,
		Price:    product.Price,
		SellerID: product.Seller.ID,
	}
}

func createPersistedSeller(t *testing.T, sellerRepo *MockSellerRepository) *entities.ValidatedSeller {
	seller := entities.NewSeller("John Doe")
	validatedSeller, err := entities.NewValidatedSeller(seller)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	err = sellerRepo.Create(validatedSeller)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	return validatedSeller
}
