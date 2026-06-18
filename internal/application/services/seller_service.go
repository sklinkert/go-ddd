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
	// Check idempotency key
	if sellerCommand.IdempotencyKey != "" {
		existingRecord, err := s.idempotencyRepo.FindByKey(ctx, sellerCommand.IdempotencyKey)
		if err != nil {
			return nil, err
		}

		if existingRecord != nil {
			// Return cached response
			var result command.CreateSellerCommandResult
			if err := json.Unmarshal([]byte(existingRecord.Response), &result); err != nil {
				return nil, err
			}
			return &result, nil
		}
	}

	// Create idempotency record
	idempotencyRecord := newIdempotencyRecord(ctx, sellerCommand.IdempotencyKey, sellerCommand)

	var newSeller = entities.NewSeller(sellerCommand.Name)

	validatedSeller, err := entities.NewValidatedSeller(newSeller)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.Create(ctx, validatedSeller)
	if err != nil {
		return nil, err
	}

	result := command.CreateSellerCommandResult{
		Result: mapper.NewSellerResultFromValidatedEntity(validatedSeller),
	}

	storeIdempotencyResponse(ctx, s.idempotencyRepo, idempotencyRecord, result)

	return &result, nil
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
	// Check idempotency key
	if updateCommand.IdempotencyKey != "" {
		existingRecord, err := s.idempotencyRepo.FindByKey(ctx, updateCommand.IdempotencyKey)
		if err != nil {
			return nil, err
		}

		if existingRecord != nil {
			// Return cached response
			var result command.UpdateSellerCommandResult
			if err := json.Unmarshal([]byte(existingRecord.Response), &result); err != nil {
				return nil, err
			}
			return &result, nil
		}
	}

	// Create idempotency record
	idempotencyRecord := newIdempotencyRecord(ctx, updateCommand.IdempotencyKey, updateCommand)

	seller, err := s.repo.FindById(ctx, updateCommand.Id)
	if err != nil {
		return nil, err
	}

	if seller == nil {
		return nil, errors.New("seller not found")
	}

	if err := seller.UpdateName(updateCommand.Name); err != nil {
		return nil, err
	}

	validatedUpdatedSeller, err := entities.NewValidatedSeller(seller)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.Update(ctx, validatedUpdatedSeller)
	if err != nil {
		return nil, err
	}

	result := command.UpdateSellerCommandResult{
		Result: mapper.NewSellerResultFromEntity(seller),
	}

	storeIdempotencyResponse(ctx, s.idempotencyRepo, idempotencyRecord, result)

	return &result, nil
}

func (s *SellerService) DeleteSeller(ctx context.Context, sellerCommand *command.DeleteSellerCommand) (*command.DeleteSellerCommandResult, error) {
	// Check idempotency key
	if sellerCommand.IdempotencyKey != "" {
		existingRecord, err := s.idempotencyRepo.FindByKey(ctx, sellerCommand.IdempotencyKey)
		if err != nil {
			return nil, err
		}

		if existingRecord != nil {
			// Return cached response
			var result command.DeleteSellerCommandResult
			if err := json.Unmarshal([]byte(existingRecord.Response), &result); err != nil {
				return nil, err
			}
			return &result, nil
		}
	}

	// Create idempotency record
	idempotencyRecord := newIdempotencyRecord(ctx, sellerCommand.IdempotencyKey, sellerCommand)

	// Check if seller exists
	existingSeller, err := s.repo.FindById(ctx, sellerCommand.Id)
	if err != nil {
		return nil, err
	}

	if existingSeller == nil {
		return nil, errors.New("seller not found")
	}

	// Delete seller
	err = s.repo.Delete(ctx, sellerCommand.Id)
	if err != nil {
		return nil, err
	}

	result := command.DeleteSellerCommandResult{
		Success: true,
	}

	storeIdempotencyResponse(ctx, s.idempotencyRepo, idempotencyRecord, result)

	return &result, nil
}
