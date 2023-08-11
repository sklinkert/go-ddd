package db

import (
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/domain/repositories"
	"gorm.io/gorm"
)

type GormSellerRepository struct {
	db *gorm.DB
}

func NewGormSellerRepository(db *gorm.DB) repositories.SellerRepository {
	return &GormSellerRepository{db: db}
}

func (repo *GormSellerRepository) Save(seller *entities.Seller) error {
	dbSeller := ToDBSeller(seller)
	return repo.db.Save(dbSeller).Error
}

func (repo *GormSellerRepository) FindByID(id int) (*entities.Seller, error) {
	var dbSeller Seller
	if err := repo.db.First(&dbSeller, id).Error; err != nil {
		return nil, err
	}
	return FromDBSeller(&dbSeller), nil
}

func (repo *GormSellerRepository) GetAll() ([]*entities.Seller, error) {
	var dbSellers []Seller
	if err := repo.db.Find(&dbSellers).Error; err != nil {
		return nil, err
	}

	sellers := make([]*entities.Seller, len(dbSellers))
	for i, dbSeller := range dbSellers {
		sellers[i] = FromDBSeller(&dbSeller)
	}
	return sellers, nil
}

func (repo *GormSellerRepository) Update(seller *entities.Seller) error {
	dbSeller := ToDBSeller(seller)
	return repo.db.Model(&Seller{}).Where("id = ?", dbSeller.ID).Updates(dbSeller).Error
}

func (repo *GormSellerRepository) Delete(id int) error {
	return repo.db.Delete(&Seller{}, id).Error
}
