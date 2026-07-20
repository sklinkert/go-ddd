package services

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/query"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

// MockProductRepository is a mock implementation of the ProductRepository interface
type MockProductRepository struct {
	products []*entities.ValidatedProduct
}

func (m *MockProductRepository) Create(ctx context.Context, product *entities.ValidatedProduct) (*entities.Product, error) {
	m.products = append(m.products, product)
	return &product.Product, nil
}

func (m *MockProductRepository) FindAll(ctx context.Context) ([]*entities.Product, error) {
	var products []*entities.Product
	for _, p := range m.products {
		products = append(products, &p.Product)
	}
	return products, nil
}

func (m *MockProductRepository) Update(ctx context.Context, product *entities.ValidatedProduct) (*entities.Product, error) {
	for index, p := range m.products {
		if p.Id == product.Id {
			m.products[index] = product
			return &product.Product, nil
		}
	}
	return nil, errors.New("product not found for update")
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	for index, p := range m.products {
		if p.Id == id {
			m.products = append(m.products[:index], m.products[index+1:]...)
			return nil
		}
	}
	return errors.New("product not found for delete")
}

func (m *MockProductRepository) FindById(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	for _, p := range m.products {
		if p.Id == id {
			return &p.Product, nil
		}
	}
	return nil, nil
}

// MockIdempotencyRepository is an in-memory implementation of the
// IdempotencyRepository interface with call tracking for assertions.
type MockIdempotencyRepository struct {
	mu             sync.Mutex
	records        map[string]*entities.IdempotencyRecord
	reserveCalls   int
	deleteCalls    int
	deletedKeys    []string
	reserveErr     error
	findErr        error
	setResponseErr error
}

func NewMockIdempotencyRepository() *MockIdempotencyRepository {
	return &MockIdempotencyRepository{
		records: make(map[string]*entities.IdempotencyRecord),
	}
}

func (m *MockIdempotencyRepository) Reserve(ctx context.Context, record *entities.IdempotencyRecord) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.reserveCalls++
	if m.reserveErr != nil {
		return false, m.reserveErr
	}
	if _, exists := m.records[record.Key]; exists {
		return false, nil
	}
	m.records[record.Key] = record
	return true, nil
}

func (m *MockIdempotencyRepository) FindByKey(ctx context.Context, key string) (*entities.IdempotencyRecord, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.findErr != nil {
		return nil, m.findErr
	}
	if record, exists := m.records[key]; exists {
		copied := *record
		return &copied, nil
	}
	return nil, nil
}

func (m *MockIdempotencyRepository) SetResponse(ctx context.Context, key string, response string, statusCode int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.setResponseErr != nil {
		return m.setResponseErr
	}
	if record, exists := m.records[key]; exists {
		record.SetResponse(response, statusCode)
	}
	return nil
}

func (m *MockIdempotencyRepository) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deleteCalls++
	m.deletedKeys = append(m.deletedKeys, key)
	delete(m.records, key)
	return nil
}

func TestProductService_CreateProduct(t *testing.T) {
	productRepo := &MockProductRepository{}
	sellerRepo := &MockSellerRepository{}
	idempotencyRepo := NewMockIdempotencyRepository()
	service := NewProductService(productRepo, sellerRepo, idempotencyRepo)

	// Create seller
	seller := createPersistedSeller(t, sellerRepo)

	// Create product
	productCommand := getCreateProductCommand("Example", 10000, seller.Id)
	_, err := service.CreateProduct(context.Background(), productCommand)
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
	idempotencyRepo := NewMockIdempotencyRepository()
	service := NewProductService(productRepo, sellerRepo, idempotencyRepo)

	// Create seller
	seller := createPersistedSeller(t, sellerRepo)

	// Add two products
	_, _ = service.CreateProduct(context.Background(), getCreateProductCommand("Example1", 10000, seller.Id))
	_, _ = service.CreateProduct(context.Background(), getCreateProductCommand("Example2", 20000, seller.Id))

	products, err := service.FindAllProducts(context.Background())
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
	idempotencyRepo := NewMockIdempotencyRepository()
	service := NewProductService(productRepo, sellerRepo, idempotencyRepo)

	// Create seller
	seller := createPersistedSeller(t, sellerRepo)

	result, err := service.CreateProduct(context.Background(), getCreateProductCommand("Example", 10000, seller.Id))
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	foundProduct, err := service.FindProductById(context.Background(), &query.GetProductByIdQuery{Id: result.Result.Id})
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if foundProduct.Result.Name != "Example" {
		t.Errorf("Expected product name 'Example', but got %s", foundProduct.Result.Name)
	}

	notFound, err := service.FindProductById(context.Background(), &query.GetProductByIdQuery{Id: uuid.New()}) // some non-existent Id
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if notFound != nil {
		t.Error("Expected nil result for non-existent product, but got one")
	}
}

func getCreateProductCommand(name string, priceMinorUnits int64, sellerId uuid.UUID) *command.CreateProductCommand {
	return &command.CreateProductCommand{
		Name:            name,
		PriceMinorUnits: priceMinorUnits,
		Currency:        entities.USD,
		SellerId:        sellerId,
	}
}

func createPersistedSeller(t *testing.T, sellerRepo *MockSellerRepository) *entities.ValidatedSeller {
	seller := entities.NewSeller("John Doe")
	validatedSeller, err := entities.NewValidatedSeller(seller)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	_, err = sellerRepo.Create(context.Background(), validatedSeller)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	return validatedSeller
}
