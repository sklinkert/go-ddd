package services

import (
	"context"

	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/interfaces"
	"github.com/sklinkert/go-ddd/internal/application/mapper"
	"github.com/sklinkert/go-ddd/internal/application/query"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/domain/repositories"
)

type SellerService struct {
	repo            repositories.SellerRepository
	idempotencyRepo repositories.IdempotencyRepository
}

// NewSellerService - Constructor for the service
func NewSellerService(repo repositories.SellerRepository, idempotencyRepo repositories.IdempotencyRepository) interfaces.SellerService {
	return &SellerService{
		repo:            repo,
		idempotencyRepo: idempotencyRepo,
	}
}

// CreateSeller saves a new seller
func (s *SellerService) CreateSeller(ctx context.Context, sellerCommand *command.CreateSellerCommand) (*command.CreateSellerCommandResult, error) {
	return withIdempotency(ctx, s.idempotencyRepo, sellerCommand.IdempotencyKey, sellerCommand, func() (*command.CreateSellerCommandResult, error) {
		newSeller := entities.NewSeller(sellerCommand.Name)

		validatedSeller, err := entities.NewValidatedSeller(newSeller)
		if err != nil {
			return nil, err
		}

		if _, err := s.repo.Create(ctx, validatedSeller); err != nil {
			return nil, err
		}

		return &command.CreateSellerCommandResult{
			Result: mapper.NewSellerResultFromValidatedEntity(validatedSeller),
		}, nil
	})
}

// FindAllSellers fetches all sellers
func (s *SellerService) FindAllSellers(ctx context.Context) (*query.GetAllSellersQueryResult, error) {
	storedSellers, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var queryResult query.GetAllSellersQueryResult
	for _, seller := range storedSellers {
		queryResult.Result = append(queryResult.Result, mapper.NewSellerResultFromEntity(seller))
	}

	return &queryResult, nil
}

// FindSellerById fetches a specific seller by Id
func (s *SellerService) FindSellerById(ctx context.Context, sellerQuery *query.GetSellerByIdQuery) (*query.GetSellerByIdQueryResult, error) {
	storedSeller, err := s.repo.FindById(ctx, sellerQuery.Id)
	if err != nil {
		return nil, err
	}

	// Not found: let the caller translate this into a 404.
	if storedSeller == nil {
		return nil, nil
	}

	var queryResult query.GetSellerByIdQueryResult
	queryResult.Result = mapper.NewSellerResultFromEntity(storedSeller)

	return &queryResult, nil
}

// UpdateSeller updates a seller
func (s *SellerService) UpdateSeller(ctx context.Context, updateCommand *command.UpdateSellerCommand) (*command.UpdateSellerCommandResult, error) {
	return withIdempotency(ctx, s.idempotencyRepo, updateCommand.IdempotencyKey, updateCommand, func() (*command.UpdateSellerCommandResult, error) {
		seller, err := s.repo.FindById(ctx, updateCommand.Id)
		if err != nil {
			return nil, err
		}

		if seller == nil {
			return nil, entities.ErrSellerNotFound
		}

		if err := seller.UpdateName(updateCommand.Name); err != nil {
			return nil, err
		}

		validatedUpdatedSeller, err := entities.NewValidatedSeller(seller)
		if err != nil {
			return nil, err
		}

		if _, err := s.repo.Update(ctx, validatedUpdatedSeller); err != nil {
			return nil, err
		}

		return &command.UpdateSellerCommandResult{
			Result: mapper.NewSellerResultFromValidatedEntity(validatedUpdatedSeller),
		}, nil
	})
}

func (s *SellerService) DeleteSeller(ctx context.Context, sellerCommand *command.DeleteSellerCommand) (*command.DeleteSellerCommandResult, error) {
	return withIdempotency(ctx, s.idempotencyRepo, sellerCommand.IdempotencyKey, sellerCommand, func() (*command.DeleteSellerCommandResult, error) {
		existingSeller, err := s.repo.FindById(ctx, sellerCommand.Id)
		if err != nil {
			return nil, err
		}

		if existingSeller == nil {
			return nil, entities.ErrSellerNotFound
		}

		if err := s.repo.Delete(ctx, sellerCommand.Id); err != nil {
			return nil, err
		}

		return &command.DeleteSellerCommandResult{Success: true}, nil
	})
}
