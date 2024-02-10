package repositories

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type SellerRepository interface {
	Create(seller *entities.ValidatedSeller) error
	FindById(id uuid.UUID) (*entities.ValidatedSeller, error)
	FindAll() ([]*entities.ValidatedSeller, error)
	Update(seller *entities.ValidatedSeller) error
	Delete(id uuid.UUID) error
}
