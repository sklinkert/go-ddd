package mapper

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSellerResultFromEntity(t *testing.T) {
	now := time.Now()
	seller := &entities.Seller{Name: "Acme", CreatedAt: now, UpdatedAt: now}
	seller.Id = uuid.New()

	result := NewSellerResultFromEntity(seller)

	assert.NotNil(t, result)
	assert.Equal(t, seller.Id, result.Id)
	assert.Equal(t, "Acme", result.Name)
	assert.Equal(t, now, result.CreatedAt)
	assert.Equal(t, now, result.UpdatedAt)
}

func TestNewSellerResultFromEntity_Nil(t *testing.T) {
	assert.Nil(t, NewSellerResultFromEntity(nil))
}

func TestNewSellerResultFromValidatedEntity(t *testing.T) {
	validated, err := entities.NewValidatedSeller(entities.NewSeller("Acme"))
	assert.NoError(t, err)

	result := NewSellerResultFromValidatedEntity(validated)

	assert.NotNil(t, result)
	assert.Equal(t, validated.Id, result.Id)
	assert.Equal(t, "Acme", result.Name)
}

func TestNewProductResultFromEntity(t *testing.T) {
	now := time.Now()
	price, err := entities.NewMoney(999, entities.USD)
	require.NoError(t, err)
	sellerId := uuid.New()
	product := &entities.Product{Name: "Widget", Price: price, SellerId: sellerId, CreatedAt: now, UpdatedAt: now}
	product.Id = uuid.New()

	result := NewProductResultFromEntity(product)

	assert.NotNil(t, result)
	assert.Equal(t, product.Id, result.Id)
	assert.Equal(t, "Widget", result.Name)
	assert.Equal(t, price, result.Price)
	assert.Equal(t, sellerId, result.SellerId)
	assert.Equal(t, now, result.CreatedAt)
	assert.Equal(t, now, result.UpdatedAt)
}

func TestNewProductResultFromEntity_Nil(t *testing.T) {
	assert.Nil(t, NewProductResultFromEntity(nil))
}

func TestNewProductResultFromValidatedEntity(t *testing.T) {
	validatedSeller, err := entities.NewValidatedSeller(entities.NewSeller("Acme"))
	require.NoError(t, err)
	price, err := entities.NewMoney(999, entities.USD)
	require.NoError(t, err)
	validatedProduct, err := entities.NewValidatedProduct(entities.NewProduct("Widget", price, *validatedSeller))
	require.NoError(t, err)

	result := NewProductResultFromValidatedEntity(validatedProduct)

	assert.NotNil(t, result)
	assert.Equal(t, validatedProduct.Id, result.Id)
	assert.Equal(t, "Widget", result.Name)
	assert.Equal(t, price, result.Price)
	assert.Equal(t, validatedSeller.Id, result.SellerId)
}
