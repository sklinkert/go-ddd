package services

import (
	"context"
	"encoding/json"
	"errors"

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
	idempotencyRepo   repositories.IdempotencyRepository
}

func NewProductService(
	productRepository repositories.ProductRepository,
	sellerRepository repositories.SellerRepository,
	idempotencyRepo repositories.IdempotencyRepository,
) interfaces.ProductService {
	return &ProductService{
		productRepository: productRepository,
		sellerRepository:  sellerRepository,
		idempotencyRepo:   idempotencyRepo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, productCommand *command.CreateProductCommand) (*command.CreateProductCommandResult, error) {
	// Check idempotency key
	if productCommand.IdempotencyKey != "" {
		existingRecord, err := s.idempotencyRepo.FindByKey(ctx, productCommand.IdempotencyKey)
		if err != nil {
			return nil, err
		}

		if existingRecord != nil {
			// Return cached response
			var result command.CreateProductCommandResult
			if err := json.Unmarshal([]byte(existingRecord.Response), &result); err != nil {
				return nil, err
			}
			return &result, nil
		}
	}

	// Create idempotency record
	idempotencyRecord := newIdempotencyRecord(ctx, productCommand.IdempotencyKey, productCommand)

	storedSeller, err := s.sellerRepository.FindById(ctx, productCommand.SellerId)
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

	_, err = s.productRepository.Create(ctx, validatedProduct)
	if err != nil {
		return nil, err
	}

	result := command.CreateProductCommandResult{
		Result: mapper.NewProductResultFromValidatedEntity(validatedProduct),
	}

	storeIdempotencyResponse(ctx, s.idempotencyRepo, idempotencyRecord, result)

	return &result, nil
}

func (s *ProductService) FindAllProducts(ctx context.Context) (*query.GetAllProductsQueryResult, error) {
	storedProducts, err := s.productRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var queryListResult query.GetAllProductsQueryResult
	for _, product := range storedProducts {
		queryListResult.Result = append(queryListResult.Result, mapper.NewProductResultFromEntity(product))
	}

	return &queryListResult, nil
}

func (s *ProductService) FindProductById(ctx context.Context, productQuery *query.GetProductByIdQuery) (*query.GetProductByIdQueryResult, error) {
	storedProduct, err := s.productRepository.FindById(ctx, productQuery.Id)
	if err != nil {
		return nil, err
	}

	// Not found: let the caller translate this into a 404.
	if storedProduct == nil {
		return nil, nil
	}

	var queryResult query.GetProductByIdQueryResult
	queryResult.Result = mapper.NewProductResultFromEntity(storedProduct)

	return &queryResult, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, productCommand *command.UpdateProductCommand) (*command.UpdateProductCommandResult, error) {
	// Check idempotency key
	if productCommand.IdempotencyKey != "" {
		existingRecord, err := s.idempotencyRepo.FindByKey(ctx, productCommand.IdempotencyKey)
		if err != nil {
			return nil, err
		}

		if existingRecord != nil {
			// Return cached response
			var result command.UpdateProductCommandResult
			if err := json.Unmarshal([]byte(existingRecord.Response), &result); err != nil {
				return nil, err
			}
			return &result, nil
		}
	}

	// Create idempotency record
	idempotencyRecord := newIdempotencyRecord(ctx, productCommand.IdempotencyKey, productCommand)

	// Find existing product
	existingProduct, err := s.productRepository.FindById(ctx, productCommand.Id)
	if err != nil {
		return nil, err
	}

	if existingProduct == nil {
		return nil, errors.New("product not found")
	}

	// Find seller if different
	if productCommand.SellerId != existingProduct.Seller.Id {
		storedSeller, err := s.sellerRepository.FindById(ctx, productCommand.SellerId)
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
		existingProduct.Seller = validatedSeller.Seller
	}

	// Update product fields
	if err := existingProduct.UpdateName(productCommand.Name); err != nil {
		return nil, err
	}

	if err := existingProduct.UpdatePrice(productCommand.Price); err != nil {
		return nil, err
	}

	validatedProduct, err := entities.NewValidatedProduct(existingProduct)
	if err != nil {
		return nil, err
	}

	_, err = s.productRepository.Update(ctx, validatedProduct)
	if err != nil {
		return nil, err
	}

	result := command.UpdateProductCommandResult{
		Result: mapper.NewProductResultFromValidatedEntity(validatedProduct),
	}

	storeIdempotencyResponse(ctx, s.idempotencyRepo, idempotencyRecord, result)

	return &result, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, productCommand *command.DeleteProductCommand) (*command.DeleteProductCommandResult, error) {
	// Check idempotency key
	if productCommand.IdempotencyKey != "" {
		existingRecord, err := s.idempotencyRepo.FindByKey(ctx, productCommand.IdempotencyKey)
		if err != nil {
			return nil, err
		}

		if existingRecord != nil {
			// Return cached response
			var result command.DeleteProductCommandResult
			if err := json.Unmarshal([]byte(existingRecord.Response), &result); err != nil {
				return nil, err
			}
			return &result, nil
		}
	}

	// Create idempotency record
	idempotencyRecord := newIdempotencyRecord(ctx, productCommand.IdempotencyKey, productCommand)

	// Check if product exists
	existingProduct, err := s.productRepository.FindById(ctx, productCommand.Id)
	if err != nil {
		return nil, err
	}

	if existingProduct == nil {
		return nil, errors.New("product not found")
	}

	// Delete product
	err = s.productRepository.Delete(ctx, productCommand.Id)
	if err != nil {
		return nil, err
	}

	result := command.DeleteProductCommandResult{
		Success: true,
	}

	storeIdempotencyResponse(ctx, s.idempotencyRepo, idempotencyRecord, result)

	return &result, nil
}
