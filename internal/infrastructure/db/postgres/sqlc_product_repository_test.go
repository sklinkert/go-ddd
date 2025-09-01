package postgres

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/testhelpers"
)

func TestSqlcProductRepository_Create(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Queries)
	sellerRepo := NewSqlcSellerRepository(testDB.Queries)

	// Create a seller first
	seller := entities.NewSeller("Test Seller")
	validatedSeller, err := entities.NewValidatedSeller(seller)
	require.NoError(t, err)

	_, err = sellerRepo.Create(validatedSeller)
	require.NoError(t, err)

	// Create a product
	product := entities.NewProduct("Test Product", 99.99, *validatedSeller)
	validatedProduct, err := entities.NewValidatedProduct(product)
	require.NoError(t, err)

	createdProduct, err := repo.Create(validatedProduct)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, createdProduct)
	assert.Equal(t, validatedProduct.Name, createdProduct.Name)
	assert.Equal(t, validatedProduct.Price, createdProduct.Price)
	assert.Equal(t, validatedSeller.Id, createdProduct.Seller.Id)
	assert.NotEqual(t, uuid.Nil, createdProduct.Id)
	assert.False(t, createdProduct.CreatedAt.IsZero())
	assert.False(t, createdProduct.UpdatedAt.IsZero())
}

func TestSqlcProductRepository_FindById(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Queries)
	sellerRepo := NewSqlcSellerRepository(testDB.Queries)

	// Create test data
	seller := entities.NewSeller("Test Seller")
	validatedSeller, err := entities.NewValidatedSeller(seller)
	require.NoError(t, err)

	_, err = sellerRepo.Create(validatedSeller)
	require.NoError(t, err)

	product := entities.NewProduct("Test Product", 99.99, *validatedSeller)
	validatedProduct, err := entities.NewValidatedProduct(product)
	require.NoError(t, err)

	createdProduct, err := repo.Create(validatedProduct)
	require.NoError(t, err)

	// Test finding by ID
	foundProduct, err := repo.FindById(createdProduct.Id)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, foundProduct)
	assert.Equal(t, createdProduct.Id, foundProduct.Id)
	assert.Equal(t, createdProduct.Name, foundProduct.Name)
	assert.Equal(t, createdProduct.Price, foundProduct.Price)
	assert.Equal(t, createdProduct.Seller.Id, foundProduct.Seller.Id)
	assert.Equal(t, createdProduct.Seller.Name, foundProduct.Seller.Name)
}

func TestSqlcProductRepository_FindById_NotFound(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Queries)

	// Test finding non-existent product
	nonExistentId := uuid.New()
	foundProduct, err := repo.FindById(nonExistentId)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, foundProduct)
}

func TestSqlcProductRepository_FindAll(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Queries)
	sellerRepo := NewSqlcSellerRepository(testDB.Queries)

	// Create test seller
	seller := entities.NewSeller("Test Seller")
	validatedSeller, err := entities.NewValidatedSeller(seller)
	require.NoError(t, err)

	_, err = sellerRepo.Create(validatedSeller)
	require.NoError(t, err)

	// Create multiple products
	product1 := entities.NewProduct("Product 1", 10.00, *validatedSeller)
	validatedProduct1, err := entities.NewValidatedProduct(product1)
	require.NoError(t, err)

	product2 := entities.NewProduct("Product 2", 20.00, *validatedSeller)
	validatedProduct2, err := entities.NewValidatedProduct(product2)
	require.NoError(t, err)

	createdProduct1, err := repo.Create(validatedProduct1)
	require.NoError(t, err)

	createdProduct2, err := repo.Create(validatedProduct2)
	require.NoError(t, err)

	// Test finding all products
	products, err := repo.FindAll()

	// Assertions
	require.NoError(t, err)
	require.Len(t, products, 2)

	// Verify both products are present
	var foundProduct1, foundProduct2 *entities.Product
	for _, p := range products {
		if p.Id == createdProduct1.Id {
			foundProduct1 = p
		} else if p.Id == createdProduct2.Id {
			foundProduct2 = p
		}
	}

	require.NotNil(t, foundProduct1)
	require.NotNil(t, foundProduct2)
	assert.Equal(t, "Product 1", foundProduct1.Name)
	assert.Equal(t, "Product 2", foundProduct2.Name)
}

func TestSqlcProductRepository_FindAll_Empty(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Queries)

	// Test finding all when no products exist
	products, err := repo.FindAll()

	// Assertions
	require.NoError(t, err)
	assert.Empty(t, products)
}

func TestSqlcProductRepository_Update(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Queries)
	sellerRepo := NewSqlcSellerRepository(testDB.Queries)

	// Create test data
	seller := entities.NewSeller("Test Seller")
	validatedSeller, err := entities.NewValidatedSeller(seller)
	require.NoError(t, err)

	_, err = sellerRepo.Create(validatedSeller)
	require.NoError(t, err)

	product := entities.NewProduct("Original Product", 50.00, *validatedSeller)
	validatedProduct, err := entities.NewValidatedProduct(product)
	require.NoError(t, err)

	createdProduct, err := repo.Create(validatedProduct)
	require.NoError(t, err)

	// Update the product
	updatedProduct := &entities.Product{
		Id:        createdProduct.Id,
		Name:      "Updated Product",
		Price:     75.00,
		Seller:    createdProduct.Seller,
		CreatedAt: createdProduct.CreatedAt,
		UpdatedAt: time.Now(),
	}

	validatedUpdatedProduct, err := entities.NewValidatedProduct(updatedProduct)
	require.NoError(t, err)

	result, err := repo.Update(validatedUpdatedProduct)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Product", result.Name)
	assert.Equal(t, 75.00, result.Price)
	assert.Equal(t, createdProduct.Id, result.Id)
	assert.True(t, result.UpdatedAt.After(createdProduct.UpdatedAt))
}

func TestSqlcProductRepository_Update_NotFound(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Queries)

	// Create a seller for the product
	seller := entities.NewSeller("Test Seller")
	validatedSeller, err := entities.NewValidatedSeller(seller)
	require.NoError(t, err)

	// Create a product with non-existent ID
	nonExistentProduct := &entities.Product{
		Id:        uuid.New(),
		Name:      "Non-existent Product",
		Price:     100.00,
		Seller:    validatedSeller.Seller,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	validatedNonExistentProduct, err := entities.NewValidatedProduct(nonExistentProduct)
	require.NoError(t, err)

	result, err := repo.Update(validatedNonExistentProduct)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSqlcProductRepository_Delete(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Queries)
	sellerRepo := NewSqlcSellerRepository(testDB.Queries)

	// Create test data
	seller := entities.NewSeller("Test Seller")
	validatedSeller, err := entities.NewValidatedSeller(seller)
	require.NoError(t, err)

	_, err = sellerRepo.Create(validatedSeller)
	require.NoError(t, err)

	product := entities.NewProduct("Test Product", 99.99, *validatedSeller)
	validatedProduct, err := entities.NewValidatedProduct(product)
	require.NoError(t, err)

	createdProduct, err := repo.Create(validatedProduct)
	require.NoError(t, err)

	// Delete the product
	err = repo.Delete(createdProduct.Id)
	require.NoError(t, err)

	// Verify product is deleted
	deletedProduct, err := repo.FindById(createdProduct.Id)
	assert.Error(t, err)
	assert.Nil(t, deletedProduct)
}

func TestSqlcProductRepository_Delete_NotFound(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Queries)

	// Try to delete non-existent product
	nonExistentId := uuid.New()
	err := repo.Delete(nonExistentId)

	// Note: PostgreSQL DELETE doesn't fail if the row doesn't exist
	// So this should not return an error
	assert.NoError(t, err)
}

func TestSqlcProductRepository_Create_WithInvalidSeller(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Queries)

	// Create a product with non-existent seller
	invalidSeller := &entities.Seller{
		Id:        uuid.New(), // This ID doesn't exist in the database
		Name:      "Non-existent Seller",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	product := &entities.Product{
		Id:        uuid.New(),
		Name:      "Test Product",
		Price:     99.99,
		Seller:    *invalidSeller,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	validatedProduct, err := entities.NewValidatedProduct(product)
	require.NoError(t, err)

	createdProduct, err := repo.Create(validatedProduct)

	// Should fail due to foreign key constraint
	assert.Error(t, err)
	assert.Nil(t, createdProduct)
}
