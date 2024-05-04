package sqlite_test

import (
	"github.com/sklinkert/go-ddd/internal/infrastructure/db/postgres"
	"testing"

	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestSellerRepositorySave(t *testing.T) {
	gormDB, cleanup := setupDatabase()
	defer cleanup()

	repo := postgres.NewGormSellerRepository(gormDB)

	seller := entities.NewSeller("John")
	validatedSeller, _ := entities.NewValidatedSeller(seller)

	_, err := repo.Create(validatedSeller)
	assert.Nil(t, err)

	// More assertions related to saving can go here.
}

func TestSellerRepositoryFindById(t *testing.T) {
	gormDB, cleanup := setupDatabase()
	defer cleanup()

	repo := postgres.NewGormSellerRepository(gormDB)

	seller := entities.NewSeller("John")
	validatedSeller, _ := entities.NewValidatedSeller(seller)
	_, err := repo.Create(validatedSeller)

	fetchedSeller, err := repo.FindById(validatedSeller.Seller.Id)
	assert.Nil(t, err)
	assert.Equal(t, "John", fetchedSeller.Name)

	// More assertions related to fetching by Id can go here.
}

func TestSellerRepositoryGetAll(t *testing.T) {
	gormDB, cleanup := setupDatabase()
	defer cleanup()

	repo := postgres.NewGormSellerRepository(gormDB)

	seller1 := entities.NewSeller("John")
	validatedSeller1, _ := entities.NewValidatedSeller(seller1)
	_, err := repo.Create(validatedSeller1)
	assert.NoError(t, err)

	seller2 := entities.NewSeller("Jane")
	validatedSeller2, _ := entities.NewValidatedSeller(seller2)
	_, err = repo.Create(validatedSeller2)
	assert.NoError(t, err)

	allSellers, err := repo.FindAll()
	assert.Nil(t, err)
	assert.Len(t, allSellers, 2)
	// You could further assert the contents of the sellers if needed.
}

func TestSellerRepositoryUpdate(t *testing.T) {
	gormDB, cleanup := setupDatabase()
	defer cleanup()

	repo := postgres.NewGormSellerRepository(gormDB)

	seller := entities.NewSeller("John")
	validatedSeller, _ := entities.NewValidatedSeller(seller)
	_, err := repo.Create(validatedSeller)
	assert.NoError(t, err)

	// UpdateName name and validate
	validatedSeller.Seller.Name = "Johnny"
	_, err = repo.Update(validatedSeller)
	assert.Nil(t, err)

	updatedSeller, _ := repo.FindById(validatedSeller.Seller.Id)
	assert.Equal(t, "Johnny", updatedSeller.Name)
}

func TestSellerRepositoryDelete(t *testing.T) {
	gormDB, cleanup := setupDatabase()
	defer cleanup()

	repo := postgres.NewGormSellerRepository(gormDB)

	seller := entities.NewSeller("John")
	validatedSeller, _ := entities.NewValidatedSeller(seller)
	_, err := repo.Create(validatedSeller)
	assert.NoError(t, err)

	err = repo.Delete(validatedSeller.Seller.Id)
	assert.Nil(t, err)

	// Try to find the deleted seller
	deletedSeller, err := repo.FindById(validatedSeller.Seller.Id)
	assert.NotNil(t, err) // Expect an error since the seller should be deleted
	assert.Nil(t, deletedSeller)
}
