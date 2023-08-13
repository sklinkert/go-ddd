package interfaces

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type SellerService interface {
	CreateSeller(sellerCommand *command.CreateSellerCommand) (*command.CreateSellerCommandResult, error)
	GetAllSellers() ([]*entities.ValidatedSeller, error)
	GetSellerByID(id uuid.UUID) (*entities.ValidatedSeller, error)
	UpdateSeller(seller *entities.ValidatedSeller) error
}
