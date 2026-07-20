package services

import (
	"context"

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
	return withIdempotency(ctx, s.idempotencyRepo, productCommand.IdempotencyKey, productCommand, func() (*command.CreateProductCommandResult, error) {
		validatedSeller, err := s.findValidatedSeller(ctx, productCommand.SellerId)
		if err != nil {
			return nil, err
		}

		price, err := entities.NewMoney(productCommand.PriceMinorUnits, productCommand.Currency)
		if err != nil {
			return nil, err
		}

		newProduct := entities.NewProduct(productCommand.Name, price, *validatedSeller)

		validatedProduct, err := entities.NewValidatedProduct(newProduct)
		if err != nil {
			return nil, err
		}

		if _, err := s.productRepository.Create(ctx, validatedProduct); err != nil {
			return nil, err
		}

		return &command.CreateProductCommandResult{
			Result: mapper.NewProductResultFromValidatedEntity(validatedProduct),
		}, nil
	})
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
	return withIdempotency(ctx, s.idempotencyRepo, productCommand.IdempotencyKey, productCommand, func() (*command.UpdateProductCommandResult, error) {
		existingProduct, err := s.productRepository.FindById(ctx, productCommand.Id)
		if err != nil {
			return nil, err
		}

		if existingProduct == nil {
			return nil, entities.ErrProductNotFound
		}

		if productCommand.SellerId != existingProduct.SellerId {
			validatedSeller, err := s.findValidatedSeller(ctx, productCommand.SellerId)
			if err != nil {
				return nil, err
			}

			if err := existingProduct.AssignSeller(*validatedSeller); err != nil {
				return nil, err
			}
		}

		if err := existingProduct.UpdateName(productCommand.Name); err != nil {
			return nil, err
		}

		price, err := entities.NewMoney(productCommand.PriceMinorUnits, productCommand.Currency)
		if err != nil {
			return nil, err
		}

		if err := existingProduct.UpdatePrice(price); err != nil {
			return nil, err
		}

		validatedProduct, err := entities.NewValidatedProduct(existingProduct)
		if err != nil {
			return nil, err
		}

		if _, err := s.productRepository.Update(ctx, validatedProduct); err != nil {
			return nil, err
		}

		return &command.UpdateProductCommandResult{
			Result: mapper.NewProductResultFromValidatedEntity(validatedProduct),
		}, nil
	})
}

func (s *ProductService) DeleteProduct(ctx context.Context, productCommand *command.DeleteProductCommand) (*command.DeleteProductCommandResult, error) {
	return withIdempotency(ctx, s.idempotencyRepo, productCommand.IdempotencyKey, productCommand, func() (*command.DeleteProductCommandResult, error) {
		existingProduct, err := s.productRepository.FindById(ctx, productCommand.Id)
		if err != nil {
			return nil, err
		}

		if existingProduct == nil {
			return nil, entities.ErrProductNotFound
		}

		if err := s.productRepository.Delete(ctx, productCommand.Id); err != nil {
			return nil, err
		}

		return &command.DeleteProductCommandResult{Success: true}, nil
	})
}

func (s *ProductService) findValidatedSeller(ctx context.Context, sellerId uuid.UUID) (*entities.ValidatedSeller, error) {
	storedSeller, err := s.sellerRepository.FindById(ctx, sellerId)
	if err != nil {
		return nil, err
	}

	if storedSeller == nil {
		return nil, entities.ErrSellerNotFound
	}

	return entities.NewValidatedSeller(storedSeller)
}
