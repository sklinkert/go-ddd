package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/domain/repositories"
	db "github.com/sklinkert/go-ddd/internal/infrastructure/db/sqlc"
)

type SqlcProductRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewSqlcProductRepository(pool *pgxpool.Pool) repositories.ProductRepository {
	return &SqlcProductRepository{pool: pool, queries: db.New(pool)}
}

// Create persists the product and its recorded domain events in one
// transaction (transactional outbox): either both are committed or neither.
// The read-after-write happens inside the same transaction, so a transient
// failure cannot surface after the commit already succeeded.
func (repo *SqlcProductRepository) Create(ctx context.Context, product *entities.ValidatedProduct) (*entities.Product, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := repo.queries.WithTx(tx)

	if _, err := qtx.CreateProduct(ctx, db.CreateProductParams{
		ID:              product.Id,
		Name:            product.Name,
		PriceMinorUnits: product.Price.MinorUnits(),
		Currency:        string(product.Price.Currency()),
		SellerID:        product.SellerId,
		CreatedAt:       timestamptzFromTime(product.CreatedAt),
		UpdatedAt:       timestamptzFromTime(product.UpdatedAt),
	}); err != nil {
		return nil, err
	}

	if err := insertOutboxEvents(ctx, qtx, product.PullEvents()); err != nil {
		return nil, err
	}

	row, err := qtx.GetProductById(ctx, product.Id)
	if err != nil {
		return nil, err
	}

	created, err := productFromRow(row.ID, row.Name, row.PriceMinorUnits, row.Currency, row.SellerID, row.CreatedAt, row.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return created, nil
}

func (repo *SqlcProductRepository) FindById(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	row, err := repo.queries.GetProductById(ctx, id)
	if err != nil {
		// A missing row is not an error: return (nil, nil) so callers can
		// translate it into a 404 instead of a 500.
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return productFromRow(row.ID, row.Name, row.PriceMinorUnits, row.Currency, row.SellerID, row.CreatedAt, row.UpdatedAt)
}

func (repo *SqlcProductRepository) FindAll(ctx context.Context) ([]*entities.Product, error) {
	rows, err := repo.queries.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]*entities.Product, len(rows))
	for i, row := range rows {
		product, err := productFromRow(row.ID, row.Name, row.PriceMinorUnits, row.Currency, row.SellerID, row.CreatedAt, row.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products[i] = product
	}

	return products, nil
}

func (repo *SqlcProductRepository) Update(ctx context.Context, product *entities.ValidatedProduct) (*entities.Product, error) {
	rows, err := repo.queries.UpdateProduct(ctx, db.UpdateProductParams{
		ID:              product.Id,
		Name:            product.Name,
		PriceMinorUnits: product.Price.MinorUnits(),
		Currency:        string(product.Price.Currency()),
		SellerID:        product.SellerId,
		UpdatedAt:       timestamptzFromTime(product.UpdatedAt),
	})
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		// Nothing matched: the product does not exist (or is soft-deleted).
		return nil, entities.ErrProductNotFound
	}

	return repo.FindById(ctx, product.Id)
}

func (repo *SqlcProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return repo.queries.DeleteProduct(ctx, id)
}

func productFromRow(id uuid.UUID, name string, priceMinorUnits int64, currency string, sellerId uuid.UUID, createdAt, updatedAt pgtype.Timestamptz) (*entities.Product, error) {
	price, err := entities.NewMoney(priceMinorUnits, entities.Currency(currency))
	if err != nil {
		return nil, err
	}

	return &entities.Product{
		Id:        id,
		Name:      name,
		Price:     price,
		SellerId:  sellerId,
		CreatedAt: timeFromTimestamptz(createdAt),
		UpdatedAt: timeFromTimestamptz(updatedAt),
	}, nil
}
