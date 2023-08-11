package repositories

import "github.com/sklinkert/go-ddd/internal/domain/entities"

type SellerRepository interface {
	Save(seller *entities.Seller) error
	FindByID(id int) (*entities.Seller, error)
	GetAll() ([]*entities.Seller, error)
	Update(seller *entities.Seller) error
	Delete(id int) error
}
