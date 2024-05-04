package postgres

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

func (repo *GormSellerRepository) Create(seller *entities.ValidatedSeller) (*entities.Seller, error) {
	dbSeller := toDBSeller(seller)

	if err := repo.db.Create(dbSeller).Error; err != nil {
		return nil, err
	}

	return repo.FindById(dbSeller.Id)
}

func (repo *GormSellerRepository) FindById(id uuid.UUID) (*entities.Seller, error) {
	var dbSeller Seller
	if err := repo.db.First(&dbSeller, id).Error; err != nil {
		return nil, err
	}
	return fromDBSeller(&dbSeller), nil
}

func (repo *GormSellerRepository) FindAll() ([]*entities.Seller, error) {
	var dbSellers []Seller
	if err := repo.db.Find(&dbSellers).Error; err != nil {
		return nil, err
	}

	sellers := make([]*entities.Seller, len(dbSellers))
	for i, dbSeller := range dbSellers {
		sellers[i] = fromDBSeller(&dbSeller)
	}

	return sellers, nil
}

func (repo *GormSellerRepository) Update(seller *entities.ValidatedSeller) (*entities.Seller, error) {
	dbSeller := toDBSeller(seller)

	err := repo.db.Model(&Seller{}).Where("id = ?", dbSeller.Id).Updates(dbSeller).Error
	if err != nil {
		return nil, err
	}

	return repo.FindById(dbSeller.Id)
}

func (repo *GormSellerRepository) Delete(id uuid.UUID) error {
	return repo.db.Delete(&Seller{}, id).Error
}
