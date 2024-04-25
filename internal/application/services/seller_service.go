package services

import (
	"errors"
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

	_, err = s.repo.Create(validatedSeller)
	if err != nil {
		return nil, err
	}

	var result command.CreateSellerCommandResult
	result.Result = mapper.NewSellerResultFromEntity(validatedSeller.Seller)

	return &result, nil
}

// FindAllSellers fetches all sellers
func (s *SellerService) FindAllSellers() ([]*entities.Seller, error) {
	return s.repo.FindAll()
}

// FindSellerById fetches a specific seller by ID
func (s *SellerService) FindSellerById(id uuid.UUID) (*entities.Seller, error) {
	return s.repo.FindById(id)
}

// UpdateSeller updates a seller
func (s *SellerService) UpdateSeller(updateCommand *command.UpdateSellerCommand) (*command.UpdateSellerCommandResult, error) {
	seller, err := s.repo.FindById(updateCommand.ID)
	if err != nil {
		return nil, err
	}

	if seller == nil {
		return nil, errors.New("seller not found")
	}

	if err := seller.Update(updateCommand.Name); err != nil {
		return nil, err
	}

	validatedUpdatedSeller, err := entities.NewValidatedSeller(seller)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.Update(validatedUpdatedSeller)
	if err != nil {
		return nil, err
	}

	var result command.UpdateSellerCommandResult
	result.Result = mapper.NewSellerResultFromEntity(*seller)

	return &result, nil
}
