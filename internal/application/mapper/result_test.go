package mapper

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/stretchr/testify/assert"
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
	seller := entities.Seller{Name: "Acme", CreatedAt: now, UpdatedAt: now}
	seller.Id = uuid.New()
	product := &entities.Product{Name: "Widget", Price: 9.99, Seller: seller, CreatedAt: now, UpdatedAt: now}
	product.Id = uuid.New()

	result := NewProductResultFromEntity(product)

	assert.NotNil(t, result)
	assert.Equal(t, product.Id, result.Id)
	assert.Equal(t, "Widget", result.Name)
	assert.Equal(t, 9.99, result.Price)
	// The nested seller must be mapped too.
	assert.NotNil(t, result.Seller)
	assert.Equal(t, seller.Id, result.Seller.Id)
	assert.Equal(t, "Acme", result.Seller.Name)
}

func TestNewProductResultFromEntity_Nil(t *testing.T) {
	assert.Nil(t, NewProductResultFromEntity(nil))
}

func TestNewProductResultFromValidatedEntity(t *testing.T) {
	validatedSeller, err := entities.NewValidatedSeller(entities.NewSeller("Acme"))
	assert.NoError(t, err)
	validatedProduct, err := entities.NewValidatedProduct(entities.NewProduct("Widget", 9.99, *validatedSeller))
	assert.NoError(t, err)

	result := NewProductResultFromValidatedEntity(validatedProduct)

	assert.NotNil(t, result)
	assert.Equal(t, validatedProduct.Id, result.Id)
	assert.Equal(t, "Widget", result.Name)
	assert.Equal(t, 9.99, result.Price)
	assert.Equal(t, validatedSeller.Id, result.Seller.Id)
}
