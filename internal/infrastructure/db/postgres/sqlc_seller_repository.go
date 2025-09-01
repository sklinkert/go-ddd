package postgres

import (
	"context"

	"github.com/google/uuid"

	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/domain/repositories"
	db "github.com/sklinkert/go-ddd/internal/infrastructure/db/sqlc"
)

type SqlcSellerRepository struct {
	queries *db.Queries
}

func NewSqlcSellerRepository(queries *db.Queries) repositories.SellerRepository {
	return &SqlcSellerRepository{queries: queries}
}

func (repo *SqlcSellerRepository) Create(seller *entities.ValidatedSeller) (*entities.Seller, error) {
	ctx := context.Background()

	createdSeller, err := repo.queries.CreateSeller(ctx, db.CreateSellerParams{
		ID:        seller.Id,
		Name:      seller.Name,
		CreatedAt: timestamptzFromTime(seller.CreatedAt),
		UpdatedAt: timestamptzFromTime(seller.UpdatedAt),
	})
	if err != nil {
		return nil, err
	}

	return repo.FindById(createdSeller.ID)
}

func (repo *SqlcSellerRepository) FindById(id uuid.UUID) (*entities.Seller, error) {
	ctx := context.Background()

	dbSeller, err := repo.queries.GetSellerById(ctx, id)
	if err != nil {
		return nil, err
	}

	return fromSqlcSeller(&dbSeller), nil
}

func (repo *SqlcSellerRepository) FindAll() ([]*entities.Seller, error) {
	ctx := context.Background()

	dbSellers, err := repo.queries.GetAllSellers(ctx)
	if err != nil {
		return nil, err
	}

	sellers := make([]*entities.Seller, len(dbSellers))
	for i, dbSeller := range dbSellers {
		sellers[i] = fromSqlcSeller(&dbSeller)
	}

	return sellers, nil
}

func (repo *SqlcSellerRepository) Update(seller *entities.ValidatedSeller) (*entities.Seller, error) {
	ctx := context.Background()

	err := repo.queries.UpdateSeller(ctx, db.UpdateSellerParams{
		ID:        seller.Id,
		Name:      seller.Name,
		UpdatedAt: timestamptzFromTime(seller.UpdatedAt),
	})
	if err != nil {
		return nil, err
	}

	return repo.FindById(seller.Id)
}

func (repo *SqlcSellerRepository) Delete(id uuid.UUID) error {
	ctx := context.Background()
	return repo.queries.DeleteSeller(ctx, id)
}

func fromSqlcSeller(dbSeller *db.Seller) *entities.Seller {
	seller := &entities.Seller{
		Name:      dbSeller.Name,
		CreatedAt: timeFromTimestamptz(dbSeller.CreatedAt),
		UpdatedAt: timeFromTimestamptz(dbSeller.UpdatedAt),
	}
	seller.Id = dbSeller.ID
	return seller
}
