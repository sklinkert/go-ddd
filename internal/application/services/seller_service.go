package services

import (
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
	return s.repo.Save(seller)
}

// GetAllSellers fetches all sellers
func (s *SellerService) GetAllSellers() ([]*entities.Seller, error) {
	return s.repo.GetAll()
}

// GetSellerByID fetches a specific seller by ID
func (s *SellerService) GetSellerByID(id int) (*entities.Seller, error) {
	return s.repo.FindByID(id)
}

// UpdateSeller updates a seller
func (s *SellerService) UpdateSeller(seller *entities.Seller) error {
	return s.repo.Update(seller)
}
