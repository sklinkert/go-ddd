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

func (repo *GormProductRepository) Create(product *entities.ValidatedProduct) (*entities.ValidatedProduct, error) {
	// Map domain entity to DB model
	dbProduct := ToDBProduct(product)

	if err := repo.db.Save(dbProduct).Error; err != nil {
		return nil, err
	}

	// read newly created row from db to never operate on data that was not persisted
	if err := repo.db.First(dbProduct, dbProduct.ID).Error; err != nil {
		return nil, err
	}

	// Map back to domain entity
	storedProduct, err := FromDBProduct(dbProduct)
	if err != nil {
		return nil, err
	}

	return storedProduct, nil
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
