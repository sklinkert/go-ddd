package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/query"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Product service: error paths ---

func TestProductService_CreateProduct_SellerNotFound(t *testing.T) {
	service := NewProductService(&MockProductRepository{}, &MockSellerRepository{}, NewMockIdempotencyRepository())

	_, err := service.CreateProduct(context.Background(), &command.CreateProductCommand{
		Name:       "Widget",
		PriceCents: 999,
		Currency:   entities.USD,
		SellerId:   uuid.New(), // no such seller persisted
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "seller not found")
}

func TestProductService_CreateProduct_InvalidCurrency(t *testing.T) {
	sellerRepo := &MockSellerRepository{}
	service := NewProductService(&MockProductRepository{}, sellerRepo, NewMockIdempotencyRepository())

	seller := createPersistedSeller(t, sellerRepo)

	_, err := service.CreateProduct(context.Background(), &command.CreateProductCommand{
		Name:       "Widget",
		PriceCents: 999,
		Currency:   "XXX",
		SellerId:   seller.Id,
	})

	assert.ErrorIs(t, err, entities.ErrValidation)
	assert.ErrorContains(t, err, `unsupported currency "XXX"`)
}

func TestProductService_UpdateProduct_NotFound(t *testing.T) {
	service := NewProductService(&MockProductRepository{}, &MockSellerRepository{}, NewMockIdempotencyRepository())

	_, err := service.UpdateProduct(context.Background(), &command.UpdateProductCommand{
		Id:         uuid.New(),
		Name:       "Widget",
		PriceCents: 999,
		Currency:   entities.USD,
		SellerId:   uuid.New(),
	})

	assert.EqualError(t, err, "product not found")
}

func TestProductService_UpdateProduct_ValidationError(t *testing.T) {
	productRepo := &MockProductRepository{}
	sellerRepo := &MockSellerRepository{}
	service := NewProductService(productRepo, sellerRepo, NewMockIdempotencyRepository())

	seller := createPersistedSeller(t, sellerRepo)
	created, err := service.CreateProduct(context.Background(), getCreateProductCommand("Widget", 999, seller.Id))
	assert.NoError(t, err)

	// Empty name must fail domain validation.
	_, err = service.UpdateProduct(context.Background(), &command.UpdateProductCommand{
		Id:         created.Result.Id,
		Name:       "",
		PriceCents: 999,
		Currency:   entities.USD,
		SellerId:   seller.Id,
	})

	assert.ErrorIs(t, err, entities.ErrValidation)
	assert.ErrorContains(t, err, "name must not be empty")
}

func TestProductService_UpdateProduct_Success(t *testing.T) {
	productRepo := &MockProductRepository{}
	sellerRepo := &MockSellerRepository{}
	service := NewProductService(productRepo, sellerRepo, NewMockIdempotencyRepository())

	seller := createPersistedSeller(t, sellerRepo)
	created, err := service.CreateProduct(context.Background(), getCreateProductCommand("Widget", 999, seller.Id))
	assert.NoError(t, err)

	updated, err := service.UpdateProduct(context.Background(), &command.UpdateProductCommand{
		Id:         created.Result.Id,
		Name:       "Widget v2",
		PriceCents: 1999,
		Currency:   entities.USD,
		SellerId:   seller.Id,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Widget v2", updated.Result.Name)
	expectedPrice, err := entities.NewMoney(1999, entities.USD)
	require.NoError(t, err)
	assert.Equal(t, expectedPrice, updated.Result.Price)
}

func TestProductService_DeleteProduct_NotFound(t *testing.T) {
	service := NewProductService(&MockProductRepository{}, &MockSellerRepository{}, NewMockIdempotencyRepository())

	_, err := service.DeleteProduct(context.Background(), &command.DeleteProductCommand{Id: uuid.New()})

	assert.EqualError(t, err, "product not found")
}

func TestProductService_DeleteProduct_Success(t *testing.T) {
	productRepo := &MockProductRepository{}
	sellerRepo := &MockSellerRepository{}
	service := NewProductService(productRepo, sellerRepo, NewMockIdempotencyRepository())

	seller := createPersistedSeller(t, sellerRepo)
	created, err := service.CreateProduct(context.Background(), getCreateProductCommand("Widget", 999, seller.Id))
	assert.NoError(t, err)

	result, err := service.DeleteProduct(context.Background(), &command.DeleteProductCommand{Id: created.Result.Id})

	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Empty(t, productRepo.products)
}

// --- Product service: idempotency replay ---

func TestProductService_CreateProduct_IdempotentReplay(t *testing.T) {
	productRepo := &MockProductRepository{}
	sellerRepo := &MockSellerRepository{}
	service := NewProductService(productRepo, sellerRepo, NewMockIdempotencyRepository())

	seller := createPersistedSeller(t, sellerRepo)
	cmd := getCreateProductCommand("Widget", 999, seller.Id)
	cmd.IdempotencyKey = "create-key"

	first, err := service.CreateProduct(context.Background(), cmd)
	assert.NoError(t, err)

	second, err := service.CreateProduct(context.Background(), cmd)
	assert.NoError(t, err)

	// The second call must replay the cached result, not create a new product.
	assert.Len(t, productRepo.products, 1)
	assert.Equal(t, first.Result.Id, second.Result.Id)
	assert.Equal(t, first.Result.Price, second.Result.Price)
}

func TestProductService_UpdateProduct_SellerChangedNotFound(t *testing.T) {
	productRepo := &MockProductRepository{}
	sellerRepo := &MockSellerRepository{}
	service := NewProductService(productRepo, sellerRepo, NewMockIdempotencyRepository())

	seller := createPersistedSeller(t, sellerRepo)
	created, err := service.CreateProduct(context.Background(), getCreateProductCommand("Widget", 999, seller.Id))
	assert.NoError(t, err)

	// A different (non-existent) seller must fail the seller lookup branch.
	_, err = service.UpdateProduct(context.Background(), &command.UpdateProductCommand{
		Id:         created.Result.Id,
		Name:       "Widget",
		PriceCents: 999,
		Currency:   entities.USD,
		SellerId:   uuid.New(),
	})

	assert.EqualError(t, err, "seller not found")
}

func TestProductService_UpdateProduct_SellerChangedSuccess(t *testing.T) {
	productRepo := &MockProductRepository{}
	sellerRepo := &MockSellerRepository{}
	service := NewProductService(productRepo, sellerRepo, NewMockIdempotencyRepository())

	sellerA := createPersistedSeller(t, sellerRepo)
	created, err := service.CreateProduct(context.Background(), getCreateProductCommand("Widget", 999, sellerA.Id))
	assert.NoError(t, err)

	// Persist a second seller and move the product to it.
	sellerB, err := entities.NewValidatedSeller(entities.NewSeller("Globex"))
	assert.NoError(t, err)
	_, err = sellerRepo.Create(context.Background(), sellerB)
	assert.NoError(t, err)

	updated, err := service.UpdateProduct(context.Background(), &command.UpdateProductCommand{
		Id:         created.Result.Id,
		Name:       "Widget",
		PriceCents: 999,
		Currency:   entities.USD,
		SellerId:   sellerB.Id,
	})

	assert.NoError(t, err)
	assert.Equal(t, sellerB.Id, updated.Result.SellerId)
}

func TestProductService_UpdateProduct_IdempotentReplay(t *testing.T) {
	productRepo := &MockProductRepository{}
	sellerRepo := &MockSellerRepository{}
	service := NewProductService(productRepo, sellerRepo, NewMockIdempotencyRepository())

	seller := createPersistedSeller(t, sellerRepo)
	created, err := service.CreateProduct(context.Background(), getCreateProductCommand("Widget", 999, seller.Id))
	assert.NoError(t, err)

	cmd := &command.UpdateProductCommand{
		IdempotencyKey: "upd-key",
		Id:             created.Result.Id,
		Name:           "Widget v2",
		PriceCents:     1999,
		Currency:       entities.USD,
		SellerId:       seller.Id,
	}

	first, err := service.UpdateProduct(context.Background(), cmd)
	assert.NoError(t, err)
	second, err := service.UpdateProduct(context.Background(), cmd)
	assert.NoError(t, err)

	assert.Equal(t, first.Result.Name, second.Result.Name)
	assert.Equal(t, first.Result.Price, second.Result.Price)
}

func TestProductService_DeleteProduct_IdempotentReplay(t *testing.T) {
	productRepo := &MockProductRepository{}
	sellerRepo := &MockSellerRepository{}
	service := NewProductService(productRepo, sellerRepo, NewMockIdempotencyRepository())

	seller := createPersistedSeller(t, sellerRepo)
	created, err := service.CreateProduct(context.Background(), getCreateProductCommand("Widget", 999, seller.Id))
	assert.NoError(t, err)

	cmd := &command.DeleteProductCommand{IdempotencyKey: "del-key", Id: created.Result.Id}

	first, err := service.DeleteProduct(context.Background(), cmd)
	assert.NoError(t, err)
	// The product is gone, but the replay returns the cached success result
	// rather than re-running (and failing with "product not found").
	second, err := service.DeleteProduct(context.Background(), cmd)
	assert.NoError(t, err)

	assert.True(t, first.Success)
	assert.True(t, second.Success)
}

// --- Seller service: error paths ---

func TestSellerService_UpdateSeller_NotFound(t *testing.T) {
	service := NewSellerService(&MockSellerRepository{}, NewMockIdempotencyRepository())

	_, err := service.UpdateSeller(context.Background(), &command.UpdateSellerCommand{Id: uuid.New(), Name: "Acme"})

	assert.EqualError(t, err, "seller not found")
}

func TestSellerService_UpdateSeller_ValidationError(t *testing.T) {
	repo := &MockSellerRepository{}
	service := NewSellerService(repo, NewMockIdempotencyRepository())

	created, err := service.CreateSeller(context.Background(), getCreateSellerCommand("Acme"))
	assert.NoError(t, err)

	_, err = service.UpdateSeller(context.Background(), &command.UpdateSellerCommand{Id: created.Result.Id, Name: ""})

	assert.ErrorIs(t, err, entities.ErrValidation)
	assert.ErrorContains(t, err, "name must not be empty")
}

func TestSellerService_DeleteSeller_NotFound(t *testing.T) {
	service := NewSellerService(&MockSellerRepository{}, NewMockIdempotencyRepository())

	_, err := service.DeleteSeller(context.Background(), &command.DeleteSellerCommand{Id: uuid.New()})

	assert.EqualError(t, err, "seller not found")
}

func TestSellerService_DeleteSeller_Success(t *testing.T) {
	repo := &MockSellerRepository{}
	service := NewSellerService(repo, NewMockIdempotencyRepository())

	created, err := service.CreateSeller(context.Background(), getCreateSellerCommand("Acme"))
	assert.NoError(t, err)

	result, err := service.DeleteSeller(context.Background(), &command.DeleteSellerCommand{Id: created.Result.Id})

	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Empty(t, repo.sellers)
}

func TestSellerService_CreateSeller_IdempotentReplay(t *testing.T) {
	repo := &MockSellerRepository{}
	service := NewSellerService(repo, NewMockIdempotencyRepository())

	cmd := getCreateSellerCommand("Acme")
	cmd.IdempotencyKey = "seller-key"

	first, err := service.CreateSeller(context.Background(), cmd)
	assert.NoError(t, err)

	second, err := service.CreateSeller(context.Background(), cmd)
	assert.NoError(t, err)

	assert.Len(t, repo.sellers, 1)
	assert.Equal(t, first.Result.Id, second.Result.Id)
}

func TestSellerService_UpdateSeller_IdempotentReplay(t *testing.T) {
	repo := &MockSellerRepository{}
	service := NewSellerService(repo, NewMockIdempotencyRepository())

	created, err := service.CreateSeller(context.Background(), getCreateSellerCommand("Acme"))
	assert.NoError(t, err)

	cmd := &command.UpdateSellerCommand{IdempotencyKey: "upd-seller", Id: created.Result.Id, Name: "Acme v2"}

	first, err := service.UpdateSeller(context.Background(), cmd)
	assert.NoError(t, err)
	second, err := service.UpdateSeller(context.Background(), cmd)
	assert.NoError(t, err)

	assert.Equal(t, first.Result.Name, second.Result.Name)
}

func TestSellerService_DeleteSeller_IdempotentReplay(t *testing.T) {
	repo := &MockSellerRepository{}
	service := NewSellerService(repo, NewMockIdempotencyRepository())

	created, err := service.CreateSeller(context.Background(), getCreateSellerCommand("Acme"))
	assert.NoError(t, err)

	cmd := &command.DeleteSellerCommand{IdempotencyKey: "del-seller", Id: created.Result.Id}

	first, err := service.DeleteSeller(context.Background(), cmd)
	assert.NoError(t, err)
	second, err := service.DeleteSeller(context.Background(), cmd)
	assert.NoError(t, err)

	assert.True(t, first.Success)
	assert.True(t, second.Success)
}

// Sanity: a not-found seller lookup yields a nil result with no error.
func TestSellerService_FindSellerById_NotFound(t *testing.T) {
	service := NewSellerService(&MockSellerRepository{}, NewMockIdempotencyRepository())

	result, err := service.FindSellerById(context.Background(), &query.GetSellerByIdQuery{Id: uuid.New()})

	assert.NoError(t, err)
	assert.Nil(t, result)
}
