package repositories

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type ProductRepository interface {
	Create(product *entities.ValidatedProduct) error
	FindByID(id uuid.UUID) (*entities.ValidatedProduct, error)
	GetAll() ([]*entities.ValidatedProduct, error)
	Update(product *entities.ValidatedProduct) error
	Delete(id uuid.UUID) error
}
