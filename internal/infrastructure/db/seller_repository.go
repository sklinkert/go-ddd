package db

import (
	"github.com/google/uuid"
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

func (repo *GormSellerRepository) Create(seller *entities.ValidatedSeller) error {
	dbSeller := ToDBSeller(seller)

	if err := repo.db.Create(dbSeller).Error; err != nil {
		return err
	}

	storedSeller, err := repo.FindByID(dbSeller.ID)
	if err != nil {
		return err
	}

	*seller = *storedSeller

	return nil
}

func (repo *GormSellerRepository) FindByID(id uuid.UUID) (*entities.ValidatedSeller, error) {
	var dbSeller Seller
	if err := repo.db.First(&dbSeller, id).Error; err != nil {
		return nil, err
	}
	return FromDBSeller(&dbSeller)
}

func (repo *GormSellerRepository) GetAll() ([]*entities.ValidatedSeller, error) {
	var err error
	var dbSellers []Seller
	if err := repo.db.Find(&dbSellers).Error; err != nil {
		return nil, err
	}

	sellers := make([]*entities.ValidatedSeller, len(dbSellers))
	for i, dbSeller := range dbSellers {
		sellers[i], err = FromDBSeller(&dbSeller)
		if err != nil {
			return nil, err
		}
	}
	return sellers, nil
}

func (repo *GormSellerRepository) Update(seller *entities.ValidatedSeller) error {
	dbSeller := ToDBSeller(seller)
	return repo.db.Model(&Seller{}).Where("id = ?", dbSeller.ID).Updates(dbSeller).Error
}

func (repo *GormSellerRepository) Delete(id uuid.UUID) error {
	return repo.db.Delete(&Seller{}, id).Error
}
