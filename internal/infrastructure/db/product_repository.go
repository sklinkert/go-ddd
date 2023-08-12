package db

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/domain/repositories"
	"gorm.io/gorm"
)

type GormProductRepository struct {
	db *gorm.DB
}

func NewGormProductRepository(db *gorm.DB) repositories.ProductRepository {
	return &GormProductRepository{db: db}
}

func (repo *GormProductRepository) Create(product *entities.ValidatedProduct) error {
	// Map domain entity to DB model
	dbProduct := ToDBProduct(product)

	if err := repo.db.Create(dbProduct).Error; err != nil {
		return err
	}

	// Read row from DB to never return different data than persisted
	storedProduct, err := repo.FindByID(dbProduct.ID)
	if err != nil {
		return err
	}

	// Map back to domain entity
	*product = *storedProduct

	return nil
}

func (repo *GormProductRepository) FindByID(id uuid.UUID) (*entities.ValidatedProduct, error) {
	var dbProduct Product
	if err := repo.db.First(&dbProduct, id).Error; err != nil {
		return nil, err
	}

	// Map back to domain entity
	return FromDBProduct(&dbProduct)
}

func (repo *GormProductRepository) GetAll() ([]*entities.ValidatedProduct, error) {
	var dbProducts []Product
	var err error

	if err := repo.db.Find(&dbProducts).Error; err != nil {
		return nil, err
	}

	products := make([]*entities.ValidatedProduct, len(dbProducts))
	for i, dbProduct := range dbProducts {
		products[i], err = FromDBProduct(&dbProduct)
		if err != nil {
			return nil, err
		}
	}
	return products, nil
}

func (repo *GormProductRepository) Update(product *entities.ValidatedProduct) error {
	dbProduct := ToDBProduct(product)
	return repo.db.Model(&Product{}).Where("id = ?", dbProduct.ID).Updates(dbProduct).Error
}

func (repo *GormProductRepository) Delete(id uuid.UUID) error {
	return repo.db.Delete(&Product{}, id).Error
}
