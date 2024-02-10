package services

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/mapper"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/domain/repositories"
)

type SellerService struct {
	repo repositories.SellerRepository
}

// NewSellerService - Constructor for the service
func NewSellerService(repo repositories.SellerRepository) *SellerService {
	return &SellerService{repo: repo}
}

// CreateSeller saves a new seller
func (s *SellerService) CreateSeller(sellerCommand *command.CreateSellerCommand) (*command.CreateSellerCommandResult, error) {
	var newSeller = entities.NewSeller(sellerCommand.Name)

	validatedSeller, err := entities.NewValidatedSeller(newSeller)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(validatedSeller)
	if err != nil {
		return nil, err
	}

	var result command.CreateSellerCommandResult
	result.Result = mapper.NewSellerResultFromEntity(*validatedSeller)

	return &result, nil
}

// GetAllSellers fetches all sellers
func (s *SellerService) GetAllSellers() ([]*entities.ValidatedSeller, error) {
	return s.repo.FindAll()
}

// GetSellerByID fetches a specific seller by ID
func (s *SellerService) GetSellerByID(id uuid.UUID) (*entities.ValidatedSeller, error) {
	return s.repo.FindById(id)
}

// UpdateSeller updates a seller
func (s *SellerService) UpdateSeller(seller *entities.ValidatedSeller) error {
	return s.repo.Update(seller)
}
