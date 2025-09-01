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
func (s *SellerService) CreateSeller(sellerCommand *command.CreateSellerCommand) (*command.CreateSellerCommandResult, error) {
	ctx := context.Background()

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
	var idempotencyRecord *entities.IdempotencyRecord
	if sellerCommand.IdempotencyKey != "" {
		requestJSON, _ := json.Marshal(sellerCommand)
		idempotencyRecord = entities.NewIdempotencyRecord(sellerCommand.IdempotencyKey, string(requestJSON))
	}

	var newSeller = entities.NewSeller(sellerCommand.Name)

	validatedSeller, err := entities.NewValidatedSeller(newSeller)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.Create(validatedSeller)
	if err != nil {
		return nil, err
	}

	result := command.CreateSellerCommandResult{
		Result: mapper.NewSellerResultFromValidatedEntity(validatedSeller),
	}

	// Store response in idempotency record
	if idempotencyRecord != nil {
		responseJSON, _ := json.Marshal(result)
		idempotencyRecord.SetResponse(string(responseJSON), 200)
		_, err = s.idempotencyRepo.Create(ctx, idempotencyRecord)
		if err != nil {
			// Log error but don't fail the operation
			// In production, you might want to handle this differently
		}
	}

	return &result, nil
}

// FindAllSellers fetches all sellers
func (s *SellerService) FindAllSellers(sellerQuery *query.GetAllSellersQuery) (*query.GetAllSellersQueryResult, error) {
	storedSellers, err := s.repo.FindAll()
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
func (s *SellerService) FindSellerById(sellerQuery *query.GetSellerByIdQuery) (*query.GetSellerByIdQueryResult, error) {
	storedSeller, err := s.repo.FindById(sellerQuery.Id)
	if err != nil {
		return nil, err
	}

	var queryResult query.GetSellerByIdQueryResult
	queryResult.Result = mapper.NewSellerResultFromEntity(storedSeller)

	return &queryResult, nil
}

// UpdateSeller updates a seller
func (s *SellerService) UpdateSeller(updateCommand *command.UpdateSellerCommand) (*command.UpdateSellerCommandResult, error) {
	ctx := context.Background()

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
	var idempotencyRecord *entities.IdempotencyRecord
	if updateCommand.IdempotencyKey != "" {
		requestJSON, _ := json.Marshal(updateCommand)
		idempotencyRecord = entities.NewIdempotencyRecord(updateCommand.IdempotencyKey, string(requestJSON))
	}

	seller, err := s.repo.FindById(updateCommand.Id)
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

	_, err = s.repo.Update(validatedUpdatedSeller)
	if err != nil {
		return nil, err
	}

	result := command.UpdateSellerCommandResult{
		Result: mapper.NewSellerResultFromEntity(seller),
	}

	// Store response in idempotency record
	if idempotencyRecord != nil {
		responseJSON, _ := json.Marshal(result)
		idempotencyRecord.SetResponse(string(responseJSON), 200)
		_, err = s.idempotencyRepo.Create(ctx, idempotencyRecord)
		if err != nil {
			// Log error but don't fail the operation
			// In production, you might want to handle this differently
		}
	}

	return &result, nil
}

func (s *SellerService) DeleteSeller(sellerCommand *command.DeleteSellerCommand) (*command.DeleteSellerCommandResult, error) {
	ctx := context.Background()

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
	var idempotencyRecord *entities.IdempotencyRecord
	if sellerCommand.IdempotencyKey != "" {
		requestJSON, _ := json.Marshal(sellerCommand)
		idempotencyRecord = entities.NewIdempotencyRecord(sellerCommand.IdempotencyKey, string(requestJSON))
	}

	// Check if seller exists
	existingSeller, err := s.repo.FindById(sellerCommand.Id)
	if err != nil {
		return nil, err
	}

	if existingSeller == nil {
		return nil, errors.New("seller not found")
	}

	// Delete seller
	err = s.repo.Delete(sellerCommand.Id)
	if err != nil {
		return nil, err
	}

	result := command.DeleteSellerCommandResult{
		Success: true,
	}

	// Store response in idempotency record
	if idempotencyRecord != nil {
		responseJSON, _ := json.Marshal(result)
		idempotencyRecord.SetResponse(string(responseJSON), 200)
		_, err = s.idempotencyRepo.Create(ctx, idempotencyRecord)
		if err != nil {
			// Log error but don't fail the operation
			// In production, you might want to handle this differently
		}
	}

	return &result, nil
}
