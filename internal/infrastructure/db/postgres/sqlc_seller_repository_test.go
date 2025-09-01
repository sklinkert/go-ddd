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

func TestSqlcSellerRepository_Create(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcSellerRepository(testDB.Queries)

	// Create a seller
	seller := entities.NewSeller("Test Seller")
	validatedSeller, err := entities.NewValidatedSeller(seller)
	require.NoError(t, err)

	createdSeller, err := repo.Create(validatedSeller)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, createdSeller)
	assert.Equal(t, validatedSeller.Name, createdSeller.Name)
	assert.NotEqual(t, uuid.Nil, createdSeller.Id)
	assert.False(t, createdSeller.CreatedAt.IsZero())
	assert.False(t, createdSeller.UpdatedAt.IsZero())
}

func TestSqlcSellerRepository_FindById(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcSellerRepository(testDB.Queries)

	// Create test data
	seller := entities.NewSeller("Test Seller")
	validatedSeller, err := entities.NewValidatedSeller(seller)
	require.NoError(t, err)

	createdSeller, err := repo.Create(validatedSeller)
	require.NoError(t, err)

	// Test finding by ID
	foundSeller, err := repo.FindById(createdSeller.Id)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, foundSeller)
	assert.Equal(t, createdSeller.Id, foundSeller.Id)
	assert.Equal(t, createdSeller.Name, foundSeller.Name)
	assert.Equal(t, createdSeller.CreatedAt.Unix(), foundSeller.CreatedAt.Unix()) // Compare Unix timestamps to avoid precision issues
	assert.Equal(t, createdSeller.UpdatedAt.Unix(), foundSeller.UpdatedAt.Unix())
}

func TestSqlcSellerRepository_FindById_NotFound(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcSellerRepository(testDB.Queries)

	// Test finding non-existent seller
	nonExistentId := uuid.New()
	foundSeller, err := repo.FindById(nonExistentId)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, foundSeller)
}

func TestSqlcSellerRepository_FindAll(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcSellerRepository(testDB.Queries)

	// Create multiple sellers
	seller1 := entities.NewSeller("Seller One")
	validatedSeller1, err := entities.NewValidatedSeller(seller1)
	require.NoError(t, err)

	seller2 := entities.NewSeller("Seller Two")
	validatedSeller2, err := entities.NewValidatedSeller(seller2)
	require.NoError(t, err)

	createdSeller1, err := repo.Create(validatedSeller1)
	require.NoError(t, err)

	createdSeller2, err := repo.Create(validatedSeller2)
	require.NoError(t, err)

	// Test finding all sellers
	sellers, err := repo.FindAll()

	// Assertions
	require.NoError(t, err)
	require.Len(t, sellers, 2)

	// Verify both sellers are present
	var foundSeller1, foundSeller2 *entities.Seller
	for _, s := range sellers {
		if s.Id == createdSeller1.Id {
			foundSeller1 = s
		} else if s.Id == createdSeller2.Id {
			foundSeller2 = s
		}
	}

	require.NotNil(t, foundSeller1)
	require.NotNil(t, foundSeller2)
	assert.Equal(t, "Seller One", foundSeller1.Name)
	assert.Equal(t, "Seller Two", foundSeller2.Name)
}

func TestSqlcSellerRepository_FindAll_Empty(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcSellerRepository(testDB.Queries)

	// Test finding all when no sellers exist
	sellers, err := repo.FindAll()

	// Assertions
	require.NoError(t, err)
	assert.Empty(t, sellers)
}

func TestSqlcSellerRepository_Update(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcSellerRepository(testDB.Queries)

	// Create test data
	seller := entities.NewSeller("Original Seller")
	validatedSeller, err := entities.NewValidatedSeller(seller)
	require.NoError(t, err)

	createdSeller, err := repo.Create(validatedSeller)
	require.NoError(t, err)

	// Update the seller
	updatedSeller := &entities.Seller{
		Id:        createdSeller.Id,
		Name:      "Updated Seller",
		CreatedAt: createdSeller.CreatedAt,
		UpdatedAt: time.Now(),
	}

	validatedUpdatedSeller, err := entities.NewValidatedSeller(updatedSeller)
	require.NoError(t, err)

	result, err := repo.Update(validatedUpdatedSeller)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Seller", result.Name)
	assert.Equal(t, createdSeller.Id, result.Id)
	assert.True(t, result.UpdatedAt.After(createdSeller.UpdatedAt))
}

func TestSqlcSellerRepository_Update_NotFound(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcSellerRepository(testDB.Queries)

	// Create a seller with non-existent ID
	nonExistentSeller := &entities.Seller{
		Id:        uuid.New(),
		Name:      "Non-existent Seller",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	validatedNonExistentSeller, err := entities.NewValidatedSeller(nonExistentSeller)
	require.NoError(t, err)

	result, err := repo.Update(validatedNonExistentSeller)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSqlcSellerRepository_Delete(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcSellerRepository(testDB.Queries)

	// Create test data
	seller := entities.NewSeller("Test Seller")
	validatedSeller, err := entities.NewValidatedSeller(seller)
	require.NoError(t, err)

	createdSeller, err := repo.Create(validatedSeller)
	require.NoError(t, err)

	// Delete the seller
	err = repo.Delete(createdSeller.Id)
	require.NoError(t, err)

	// Verify seller is deleted
	deletedSeller, err := repo.FindById(createdSeller.Id)
	assert.Error(t, err)
	assert.Nil(t, deletedSeller)
}

func TestSqlcSellerRepository_Delete_NotFound(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcSellerRepository(testDB.Queries)

	// Try to delete non-existent seller
	nonExistentId := uuid.New()
	err := repo.Delete(nonExistentId)

	// Note: PostgreSQL DELETE doesn't fail if the row doesn't exist
	// So this should not return an error
	assert.NoError(t, err)
}

func TestSqlcSellerRepository_Delete_WithExistingProducts(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	sellerRepo := NewSqlcSellerRepository(testDB.Queries)
	productRepo := NewSqlcProductRepository(testDB.Queries)

	// Create a seller
	seller := entities.NewSeller("Test Seller")
	validatedSeller, err := entities.NewValidatedSeller(seller)
	require.NoError(t, err)

	createdSeller, err := sellerRepo.Create(validatedSeller)
	require.NoError(t, err)

	// Create a product for the seller
	product := entities.NewProduct("Test Product", 99.99, *validatedSeller)
	validatedProduct, err := entities.NewValidatedProduct(product)
	require.NoError(t, err)

	_, err = productRepo.Create(validatedProduct)
	require.NoError(t, err)

	// Try to delete the seller - should fail due to foreign key constraint
	err = sellerRepo.Delete(createdSeller.Id)
	assert.Error(t, err) // Should fail due to foreign key constraint
}

func TestSqlcSellerRepository_Create_EmptyName(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	// Create a seller with empty name
	seller := &entities.Seller{
		Id:        uuid.New(),
		Name:      "", // Empty name
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// This should fail at the validation level
	_, err := entities.NewValidatedSeller(seller)
	assert.Error(t, err) // Should fail validation before reaching repository
}

func TestSqlcSellerRepository_Create_LongName(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcSellerRepository(testDB.Queries)

	// Create a seller with very long name
	longName := make([]byte, 1000)
	for i := range longName {
		longName[i] = 'A'
	}

	seller := entities.NewSeller(string(longName))
	validatedSeller, err := entities.NewValidatedSeller(seller)
	require.NoError(t, err)

	createdSeller, err := repo.Create(validatedSeller)

	// Should succeed as TEXT fields in PostgreSQL can handle large strings
	require.NoError(t, err)
	require.NotNil(t, createdSeller)
	assert.Equal(t, string(longName), createdSeller.Name)
}
