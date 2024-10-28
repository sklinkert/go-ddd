package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/interfaces"
	"github.com/sklinkert/go-ddd/internal/application/mapper"
	"github.com/sklinkert/go-ddd/internal/application/query"
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
	storedSeller, err := s.sellerRepository.FindById(productCommand.SellerId)
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

	result := command.CreateProductCommandResult{
		Result: mapper.NewProductResultFromValidatedEntity(validatedProduct),
	}

	return &result, nil
}

func (s *ProductService) FindAllProducts() (*query.ProductQueryListResult, error) {
	storedProducts, err := s.productRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var queryListResult query.ProductQueryListResult
	for _, product := range storedProducts {
		queryListResult.Result = append(queryListResult.Result, mapper.NewProductResultFromEntity(product))
	}

	return &queryListResult, nil
}

func (s *ProductService) FindProductById(id uuid.UUID) (*query.ProductQueryResult, error) {
	storedProduct, err := s.productRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	var queryResult query.ProductQueryResult
	queryResult.Result = mapper.NewProductResultFromEntity(storedProduct)

	return &queryResult, nil
}
