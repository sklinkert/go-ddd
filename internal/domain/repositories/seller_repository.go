package repositories

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type SellerRepository interface {
	Save(seller *entities.ValidatedSeller) error
	FindByID(id uuid.UUID) (*entities.ValidatedSeller, error)
	GetAll() ([]*entities.ValidatedSeller, error)
	Update(seller *entities.ValidatedSeller) error
	Delete(id uuid.UUID) error
}
