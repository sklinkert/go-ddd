package db

import (
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

func (repo *GormProductRepository) Save(product *entities.Product) error {
	// Map domain entity to DB model
	dbProduct := ToDBProduct(product)

	return repo.db.Save(dbProduct).Error
}

func (repo *GormProductRepository) FindByID(id int) (*entities.Product, error) {
	var dbProduct Product
	if err := repo.db.First(&dbProduct, id).Error; err != nil {
		return nil, err
	}

	// Map back to domain entity
	return FromDBProduct(&dbProduct), nil
}

func (repo *GormProductRepository) GetAll() ([]*entities.Product, error) {
	var dbProducts []Product
	if err := repo.db.Find(&dbProducts).Error; err != nil {
		return nil, err
	}

	products := make([]*entities.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		products[i] = FromDBProduct(&dbProduct)
	}
	return products, nil
}

func (repo *GormProductRepository) Update(product *entities.Product) error {
	dbProduct := ToDBProduct(product)
	return repo.db.Model(&Product{}).Where("id = ?", dbProduct.ID).Updates(dbProduct).Error
}

func (repo *GormProductRepository) Delete(id int) error {
	return repo.db.Delete(&Product{}, id).Error
}
