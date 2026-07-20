package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/testhelpers"
)

func mustMoney(t *testing.T, minorUnits int64, currency entities.Currency) entities.Money {
	t.Helper()
	money, err := entities.NewMoney(minorUnits, currency)
	require.NoError(t, err)
	return money
}

func createTestSeller(t *testing.T, testDB *testhelpers.PostgresTestContainer, name string) *entities.ValidatedSeller {
	t.Helper()
	sellerRepo := NewSqlcSellerRepository(testDB.Queries)

	seller := entities.NewSeller(name)
	validatedSeller, err := entities.NewValidatedSeller(seller)
	require.NoError(t, err)

	_, err = sellerRepo.Create(context.Background(), validatedSeller)
	require.NoError(t, err)

	return validatedSeller
}

func TestSqlcProductRepository_Create(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Pool)
	validatedSeller := createTestSeller(t, testDB, "Test Seller")

	// Create a product
	product := entities.NewProduct("Test Product", mustMoney(t, 9999, entities.USD), *validatedSeller)
	validatedProduct, err := entities.NewValidatedProduct(product)
	require.NoError(t, err)

	createdProduct, err := repo.Create(context.Background(), validatedProduct)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, createdProduct)
	assert.Equal(t, validatedProduct.Name, createdProduct.Name)
	assert.Equal(t, validatedProduct.Price, createdProduct.Price)
	assert.Equal(t, int64(9999), createdProduct.Price.MinorUnits())
	assert.Equal(t, entities.USD, createdProduct.Price.Currency())
	assert.Equal(t, validatedSeller.Id, createdProduct.SellerId)
	assert.NotEqual(t, uuid.Nil, createdProduct.Id)
	assert.False(t, createdProduct.CreatedAt.IsZero())
	assert.False(t, createdProduct.UpdatedAt.IsZero())
}

func TestSqlcProductRepository_FindById(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Pool)
	validatedSeller := createTestSeller(t, testDB, "Test Seller")

	product := entities.NewProduct("Test Product", mustMoney(t, 9999, entities.EUR), *validatedSeller)
	validatedProduct, err := entities.NewValidatedProduct(product)
	require.NoError(t, err)

	createdProduct, err := repo.Create(context.Background(), validatedProduct)
	require.NoError(t, err)

	// Test finding by ID
	foundProduct, err := repo.FindById(context.Background(), createdProduct.Id)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, foundProduct)
	assert.Equal(t, createdProduct.Id, foundProduct.Id)
	assert.Equal(t, createdProduct.Name, foundProduct.Name)
	assert.Equal(t, createdProduct.Price, foundProduct.Price)
	assert.Equal(t, validatedSeller.Id, foundProduct.SellerId)
}

func TestSqlcProductRepository_FindById_NotFound(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Pool)

	// Test finding non-existent product
	nonExistentId := uuid.New()
	foundProduct, err := repo.FindById(context.Background(), nonExistentId)

	// Assertions
	assert.NoError(t, err)
	assert.Nil(t, foundProduct)
}

func TestSqlcProductRepository_FindAll(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Pool)
	validatedSeller := createTestSeller(t, testDB, "Test Seller")

	// Create multiple products
	product1 := entities.NewProduct("Product 1", mustMoney(t, 1000, entities.USD), *validatedSeller)
	validatedProduct1, err := entities.NewValidatedProduct(product1)
	require.NoError(t, err)

	product2 := entities.NewProduct("Product 2", mustMoney(t, 2000, entities.EUR), *validatedSeller)
	validatedProduct2, err := entities.NewValidatedProduct(product2)
	require.NoError(t, err)

	createdProduct1, err := repo.Create(context.Background(), validatedProduct1)
	require.NoError(t, err)

	createdProduct2, err := repo.Create(context.Background(), validatedProduct2)
	require.NoError(t, err)

	// Test finding all products
	products, err := repo.FindAll(context.Background())

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
	assert.Equal(t, mustMoney(t, 1000, entities.USD), foundProduct1.Price)
	assert.Equal(t, "Product 2", foundProduct2.Name)
	assert.Equal(t, mustMoney(t, 2000, entities.EUR), foundProduct2.Price)
}

func TestSqlcProductRepository_FindAll_Empty(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Pool)

	// Test finding all when no products exist
	products, err := repo.FindAll(context.Background())

	// Assertions
	require.NoError(t, err)
	assert.Empty(t, products)
}

func TestSqlcProductRepository_Update(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Pool)
	validatedSeller := createTestSeller(t, testDB, "Test Seller")

	product := entities.NewProduct("Original Product", mustMoney(t, 5000, entities.USD), *validatedSeller)
	validatedProduct, err := entities.NewValidatedProduct(product)
	require.NoError(t, err)

	createdProduct, err := repo.Create(context.Background(), validatedProduct)
	require.NoError(t, err)

	// Update the product
	updatedProduct := &entities.Product{
		Id:        createdProduct.Id,
		Name:      "Updated Product",
		Price:     mustMoney(t, 7500, entities.USD),
		SellerId:  createdProduct.SellerId,
		CreatedAt: createdProduct.CreatedAt,
		UpdatedAt: time.Now(),
	}

	validatedUpdatedProduct, err := entities.NewValidatedProduct(updatedProduct)
	require.NoError(t, err)

	result, err := repo.Update(context.Background(), validatedUpdatedProduct)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Product", result.Name)
	assert.Equal(t, mustMoney(t, 7500, entities.USD), result.Price)
	assert.Equal(t, createdProduct.Id, result.Id)
	assert.True(t, result.UpdatedAt.After(createdProduct.UpdatedAt))
}

func TestSqlcProductRepository_Update_NotFound(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Pool)

	// Create a product with non-existent ID
	nonExistentProduct := &entities.Product{
		Id:        uuid.New(),
		Name:      "Non-existent Product",
		Price:     mustMoney(t, 10000, entities.USD),
		SellerId:  uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	validatedNonExistentProduct, err := entities.NewValidatedProduct(nonExistentProduct)
	require.NoError(t, err)

	result, err := repo.Update(context.Background(), validatedNonExistentProduct)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSqlcProductRepository_Delete(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Pool)
	validatedSeller := createTestSeller(t, testDB, "Test Seller")

	product := entities.NewProduct("Test Product", mustMoney(t, 9999, entities.USD), *validatedSeller)
	validatedProduct, err := entities.NewValidatedProduct(product)
	require.NoError(t, err)

	createdProduct, err := repo.Create(context.Background(), validatedProduct)
	require.NoError(t, err)

	// Delete the product
	err = repo.Delete(context.Background(), createdProduct.Id)
	require.NoError(t, err)

	// Verify product is deleted
	deletedProduct, err := repo.FindById(context.Background(), createdProduct.Id)
	assert.NoError(t, err)
	assert.Nil(t, deletedProduct)
}

func TestSqlcProductRepository_Delete_NotFound(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Pool)

	// Try to delete non-existent product
	nonExistentId := uuid.New()
	err := repo.Delete(context.Background(), nonExistentId)

	// Note: PostgreSQL DELETE doesn't fail if the row doesn't exist
	// So this should not return an error
	assert.NoError(t, err)
}

func TestSqlcProductRepository_Create_WithInvalidSeller(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Pool)

	// Reference a seller ID that doesn't exist in the database
	product := &entities.Product{
		Id:        uuid.New(),
		Name:      "Test Product",
		Price:     mustMoney(t, 9999, entities.USD),
		SellerId:  uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	validatedProduct, err := entities.NewValidatedProduct(product)
	require.NoError(t, err)

	createdProduct, err := repo.Create(context.Background(), validatedProduct)

	// Should fail due to foreign key constraint
	assert.Error(t, err)
	assert.Nil(t, createdProduct)
}
