package interfaces

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type SellerService interface {
	CreateSeller(seller *entities.Seller) error
	GetAllSellers() ([]*entities.ValidatedSeller, error)
	GetSellerByID(id uuid.UUID) (*entities.ValidatedSeller, error)
	UpdateSeller(seller *entities.ValidatedSeller) error
}
