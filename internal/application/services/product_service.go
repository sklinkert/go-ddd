package services

import (
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/domain/repositories"
)

type ProductService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(product *entities.Product) error {
	return s.repo.Save(product)
}

func (s *ProductService) GetAllProducts() ([]*entities.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) FindProductByID(id int) (*entities.Product, error) {
	return s.repo.FindByID(id)
}
