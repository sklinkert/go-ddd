package repositories

import "github.com/sklinkert/go-ddd/internal/domain/entities"

type ProductRepository interface {
	Save(product *entities.Product) error
	FindByID(id int) (*entities.Product, error)
	GetAll() ([]*entities.Product, error)
	Update(product *entities.Product) error
	Delete(id int) error
}
