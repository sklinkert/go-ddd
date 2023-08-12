package services

import (
	"github.com/google/uuid"
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
func (s *SellerService) CreateSeller(seller *entities.Seller) error {
	validatedSeller, err := entities.NewValidatedSeller(seller)
	if err != nil {
		return err
	}

	return s.repo.Create(validatedSeller)
}

// GetAllSellers fetches all sellers
func (s *SellerService) GetAllSellers() ([]*entities.ValidatedSeller, error) {
	return s.repo.GetAll()
}

// GetSellerByID fetches a specific seller by ID
func (s *SellerService) GetSellerByID(id uuid.UUID) (*entities.ValidatedSeller, error) {
	return s.repo.FindByID(id)
}

// UpdateSeller updates a seller
func (s *SellerService) UpdateSeller(seller *entities.ValidatedSeller) error {
	return s.repo.Update(seller)
}
