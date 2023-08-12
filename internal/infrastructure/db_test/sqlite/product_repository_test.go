package sqlite_test

import (
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/infrastructure/db/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupDatabase() (*gorm.DB, func()) {
	// Use sqlite for testing purposes
	database, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// AutoMigrate our Product model
	err = database.AutoMigrate(&postgres.Product{}, &postgres.Seller{})
	if err != nil {
		panic("Failed to migrate database")
	}

	// Cleanup function to truncate tables
	cleanup := func() {
		database.Exec("DELETE FROM sellers")
		database.Exec("DELETE FROM products")
	}

	return database, cleanup

}

func TestGormProductRepository_Save(t *testing.T) {
	gormDB, cleanup := setupDatabase()
	defer cleanup()

	repo := postgres.NewGormProductRepository(gormDB)

	seller := entities.NewSeller("TestSeller")
	product := entities.NewProduct("TestProduct", 9.99, seller)
	validProduct, _ := entities.NewValidatedProduct(product)

	err := repo.Create(validProduct)
	if err != nil {
		t.Errorf("Unexpected error during save: %s", err)
	}
}

func TestGormProductRepository_FindByID(t *testing.T) {
	gormDB, cleanup := setupDatabase()
	defer cleanup()

	repo := postgres.NewGormProductRepository(gormDB)

	seller := entities.NewSeller("TestSeller")
	product := entities.NewProduct("TestProduct", 9.99, seller)
	validProduct, _ := entities.NewValidatedProduct(product)
	repo.Create(validProduct)

	foundProduct, err := repo.FindByID(validProduct.ID)
	if err != nil || foundProduct.Name != "TestProduct" {
		t.Error("Error fetching or product mismatch")
	}
}

func TestGormProductRepository_Update(t *testing.T) {
	gormDB, cleanup := setupDatabase()
	defer cleanup()

	repo := postgres.NewGormProductRepository(gormDB)

	seller := entities.NewSeller("TestSeller")
	product := entities.NewProduct("TestProduct", 9.99, seller)
	validProduct, _ := entities.NewValidatedProduct(product)
	repo.Create(validProduct)

	validProduct.Name = "UpdatedProduct"
	err := repo.Update(validProduct)
	if err != nil {
		t.Errorf("Unexpected error during update: %s", err)
	}

	updatedProduct, _ := repo.FindByID(validProduct.ID)
	if updatedProduct.Name != "UpdatedProduct" {
		t.Error("Update failed or fetched wrong product")
	}
}

func TestGormProductRepository_GetAll(t *testing.T) {
	gormDB, cleanup := setupDatabase()
	defer cleanup()

	repo := postgres.NewGormProductRepository(gormDB)

	seller := entities.NewSeller("TestSeller")
	product := entities.NewProduct("TestProduct", 9.99, seller)
	validProduct, _ := entities.NewValidatedProduct(product)
	repo.Create(validProduct)

	products, err := repo.GetAll()
	if err != nil || len(products) != 1 {
		t.Error("Error fetching all products or product count mismatch")
	}
}

func TestGormProductRepository_Delete(t *testing.T) {
	gormDB, cleanup := setupDatabase()
	defer cleanup()

	repo := postgres.NewGormProductRepository(gormDB)

	seller := entities.NewSeller("TestSeller")
	product := entities.NewProduct("TestProduct", 9.99, seller)
	validProduct, _ := entities.NewValidatedProduct(product)
	repo.Create(validProduct)

	err := repo.Delete(validProduct.ID)
	if err != nil {
		t.Errorf("Unexpected error during delete: %s", err)
	}

	_, err = repo.FindByID(validProduct.ID)
	if err == nil {
		t.Error("Product should have been deleted, but was found")
	}
}
