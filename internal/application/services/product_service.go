package services

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/domain/repositories"
)

type ProductService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(product *entities.Product) (*entities.ValidatedProduct, error) {
	validatedProduct, err := entities.NewValidatedProduct(product)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(validatedProduct)
}

func (s *ProductService) GetAllProducts() ([]*entities.ValidatedProduct, error) {
	return s.repo.GetAll()
}

func (s *ProductService) FindProductByID(id uuid.UUID) (*entities.ValidatedProduct, error) {
	return s.repo.FindByID(id)
}
