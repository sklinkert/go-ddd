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

func (m *MockProductRepository) Create(product *entities.ValidatedProduct) (*entities.Product, error) {
	m.products = append(m.products, product)
	return &product.Product, nil
}

func (m *MockProductRepository) FindAll() ([]*entities.Product, error) {
	var products []*entities.Product
	for _, p := range m.products {
		products = append(products, &p.Product)
	}
	return products, nil
}

func (m *MockProductRepository) Update(product *entities.ValidatedProduct) (*entities.Product, error) {
	for index, p := range m.products {
		if p.Id == product.Id {
			m.products[index] = product
			return &product.Product, nil
		}
	}
	return nil, errors.New("product not found for update")
}

func (m *MockProductRepository) Delete(id uuid.UUID) error {
	for index, p := range m.products {
		if p.Id == id {
			m.products = append(m.products[:index], m.products[index+1:]...)
			return nil
		}
	}
	return errors.New("product not found for delete")
}

func (m *MockProductRepository) FindById(id uuid.UUID) (*entities.Product, error) {
	for _, p := range m.products {
		if p.Id == id {
			return &p.Product, nil
		}
		fmt.Printf("Id: mem:%s - %s\n", p.Id, id)
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

	products, err := service.FindAllProducts()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if len(products.Result) != 2 {
		t.Errorf("Expected 2 products, but got %d", len(products.Result))
	}
}

func TestProductService_FindProductById(t *testing.T) {
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

	foundProduct, err := service.FindProductById(result.Result.Id)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if foundProduct.Result.Name != "Example" {
		t.Errorf("Expected product name 'Example', but got %s", foundProduct.Result.Name)
	}

	_, err = service.FindProductById(uuid.New()) // some non-existent Id
	if err == nil {
		t.Error("Expected error for non-existent product, but got none")
	}
}

func getCreateProductCommand(product *entities.Product) *command.CreateProductCommand {
	return &command.CreateProductCommand{
		Name:     product.Name,
		Price:    product.Price,
		SellerId: product.Seller.Id,
	}
}

func createPersistedSeller(t *testing.T, sellerRepo *MockSellerRepository) *entities.ValidatedSeller {
	seller := entities.NewSeller("John Doe")
	validatedSeller, err := entities.NewValidatedSeller(seller)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	_, err = sellerRepo.Create(validatedSeller)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	return validatedSeller
}
