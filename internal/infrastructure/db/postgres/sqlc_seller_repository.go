package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

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

func (repo *SqlcSellerRepository) Create(ctx context.Context, seller *entities.ValidatedSeller) (*entities.Seller, error) {
	createdSeller, err := repo.queries.CreateSeller(ctx, db.CreateSellerParams{
		ID:        seller.Id,
		Name:      seller.Name,
		CreatedAt: timestamptzFromTime(seller.CreatedAt),
		UpdatedAt: timestamptzFromTime(seller.UpdatedAt),
	})
	if err != nil {
		return nil, err
	}

	return repo.FindById(ctx, createdSeller.ID)
}

func (repo *SqlcSellerRepository) FindById(ctx context.Context, id uuid.UUID) (*entities.Seller, error) {
	dbSeller, err := repo.queries.GetSellerById(ctx, id)
	if err != nil {
		// A missing row is not an error: return (nil, nil) so callers can
		// translate it into a 404 instead of a 500.
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return fromSqlcSellerRow(&dbSeller), nil
}

func (repo *SqlcSellerRepository) FindAll(ctx context.Context) ([]*entities.Seller, error) {
	dbSellers, err := repo.queries.GetAllSellers(ctx)
	if err != nil {
		return nil, err
	}

	sellers := make([]*entities.Seller, len(dbSellers))
	for i, dbSeller := range dbSellers {
		sellers[i] = fromSqlcSellerAllRow(&dbSeller)
	}

	return sellers, nil
}

func (repo *SqlcSellerRepository) Update(ctx context.Context, seller *entities.ValidatedSeller) (*entities.Seller, error) {
	rows, err := repo.queries.UpdateSeller(ctx, db.UpdateSellerParams{
		ID:        seller.Id,
		Name:      seller.Name,
		UpdatedAt: timestamptzFromTime(seller.UpdatedAt),
	})
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		// Nothing matched: the seller does not exist (or is soft-deleted).
		return nil, pgx.ErrNoRows
	}

	return repo.FindById(ctx, seller.Id)
}

func (repo *SqlcSellerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return repo.queries.DeleteSeller(ctx, id)
}

func fromSqlcSellerRow(dbSeller *db.GetSellerByIdRow) *entities.Seller {
	seller := &entities.Seller{
		Name:      dbSeller.Name,
		CreatedAt: timeFromTimestamptz(dbSeller.CreatedAt),
		UpdatedAt: timeFromTimestamptz(dbSeller.UpdatedAt),
	}
	seller.Id = dbSeller.ID
	return seller
}

func fromSqlcSellerAllRow(dbSeller *db.GetAllSellersRow) *entities.Seller {
	seller := &entities.Seller{
		Name:      dbSeller.Name,
		CreatedAt: timeFromTimestamptz(dbSeller.CreatedAt),
		UpdatedAt: timeFromTimestamptz(dbSeller.UpdatedAt),
	}
	seller.Id = dbSeller.ID
	return seller
}
