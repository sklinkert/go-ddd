package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/interfaces"
	"github.com/sklinkert/go-ddd/internal/application/mapper"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/domain/repositories"
)

type ProductService struct {
	productRepository repositories.ProductRepository
	sellerRepository  repositories.SellerRepository
}

func NewProductService(
	productRepository repositories.ProductRepository,
	sellerRepository repositories.SellerRepository,
) interfaces.ProductService {
	return &ProductService{productRepository: productRepository, sellerRepository: sellerRepository}
}

func (s *ProductService) CreateProduct(productCommand *command.CreateProductCommand) (*command.CreateProductCommandResult, error) {
	storedSeller, err := s.sellerRepository.FindById(productCommand.SellerID)
	if err != nil {
		return nil, err
	}

	if storedSeller == nil {
		return nil, errors.New("seller not found")
	}

	validatedSeller, err := entities.NewValidatedSeller(storedSeller)
	if err != nil {
		return nil, err
	}

	var newProduct = entities.NewProduct(
		productCommand.Name,
		productCommand.Price,
		*validatedSeller,
	)

	validatedProduct, err := entities.NewValidatedProduct(newProduct)
	if err != nil {
		return nil, err
	}

	_, err = s.productRepository.Create(validatedProduct)
	if err != nil {
		return nil, err
	}

	var result command.CreateProductCommandResult
	result.Result = mapper.NewProductResultFromEntity(validatedProduct)

	return &result, nil
}

func (s *ProductService) FindAllProducts() ([]*entities.Product, error) {
	return s.productRepository.FindAll()
}

func (s *ProductService) FindProductById(id uuid.UUID) (*entities.Product, error) {
	return s.productRepository.FindById(id)
}
