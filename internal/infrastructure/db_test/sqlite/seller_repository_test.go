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

	err := repo.Create(validatedSeller)
	assert.Nil(t, err)

	// More assertions related to saving can go here.
}

func TestSellerRepositoryFindById(t *testing.T) {
	gormDB, cleanup := setupDatabase()
	defer cleanup()

	repo := postgres.NewGormSellerRepository(gormDB)

	seller := entities.NewSeller("John")
	validatedSeller, _ := entities.NewValidatedSeller(seller)
	_ = repo.Create(validatedSeller)

	fetchedSeller, err := repo.FindById(validatedSeller.Seller.ID)
	assert.Nil(t, err)
	assert.Equal(t, "John", fetchedSeller.Seller.Name)

	// More assertions related to fetching by ID can go here.
}

func TestSellerRepositoryGetAll(t *testing.T) {
	gormDB, cleanup := setupDatabase()
	defer cleanup()

	repo := postgres.NewGormSellerRepository(gormDB)

	seller1 := entities.NewSeller("John")
	validatedSeller1, _ := entities.NewValidatedSeller(seller1)
	_ = repo.Create(validatedSeller1)

	seller2 := entities.NewSeller("Jane")
	validatedSeller2, _ := entities.NewValidatedSeller(seller2)
	_ = repo.Create(validatedSeller2)

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
	_ = repo.Create(validatedSeller)

	// Update name and validate
	validatedSeller.Seller.Name = "Johnny"
	err := repo.Update(validatedSeller)
	assert.Nil(t, err)

	updatedSeller, _ := repo.FindById(validatedSeller.Seller.ID)
	assert.Equal(t, "Johnny", updatedSeller.Seller.Name)
}

func TestSellerRepositoryDelete(t *testing.T) {
	gormDB, cleanup := setupDatabase()
	defer cleanup()

	repo := postgres.NewGormSellerRepository(gormDB)

	seller := entities.NewSeller("John")
	validatedSeller, _ := entities.NewValidatedSeller(seller)
	_ = repo.Create(validatedSeller)

	err := repo.Delete(validatedSeller.Seller.ID)
	assert.Nil(t, err)

	// Try to find the deleted seller
	deletedSeller, err := repo.FindById(validatedSeller.Seller.ID)
	assert.NotNil(t, err) // Expect an error since the seller should be deleted
	assert.Nil(t, deletedSeller)
}
