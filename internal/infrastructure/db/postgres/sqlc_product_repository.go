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

type SqlcProductRepository struct {
	queries *db.Queries
}

func NewSqlcProductRepository(queries *db.Queries) repositories.ProductRepository {
	return &SqlcProductRepository{queries: queries}
}

func (repo *SqlcProductRepository) Create(ctx context.Context, product *entities.ValidatedProduct) (*entities.Product, error) {
	createdProduct, err := repo.queries.CreateProduct(ctx, db.CreateProductParams{
		ID:        product.Id,
		Name:      product.Name,
		Price:     numericFromFloat64(product.Price),
		SellerID:  product.Seller.Id,
		CreatedAt: timestamptzFromTime(product.CreatedAt),
		UpdatedAt: timestamptzFromTime(product.UpdatedAt),
	})
	if err != nil {
		return nil, err
	}

	return repo.FindById(ctx, createdProduct.ID)
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

	return fromSqlcProductRow(&row), nil
}

func (repo *SqlcProductRepository) FindAll(ctx context.Context) ([]*entities.Product, error) {
	rows, err := repo.queries.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]*entities.Product, len(rows))
	for i, row := range rows {
		products[i] = fromSqlcProductRowAll(&row)
	}

	return products, nil
}

func (repo *SqlcProductRepository) Update(ctx context.Context, product *entities.ValidatedProduct) (*entities.Product, error) {
	rows, err := repo.queries.UpdateProduct(ctx, db.UpdateProductParams{
		ID:        product.Id,
		Name:      product.Name,
		Price:     numericFromFloat64(product.Price),
		SellerID:  product.Seller.Id,
		UpdatedAt: timestamptzFromTime(product.UpdatedAt),
	})
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		// Nothing matched: the product does not exist (or is soft-deleted).
		return nil, pgx.ErrNoRows
	}

	return repo.FindById(ctx, product.Id)
}

func (repo *SqlcProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return repo.queries.DeleteProduct(ctx, id)
}

func fromSqlcProductRow(row *db.GetProductByIdRow) *entities.Product {
	seller := &entities.Seller{
		Id:        row.SID,
		Name:      row.SName,
		CreatedAt: timeFromTimestamptz(row.SCreatedAt),
		UpdatedAt: timeFromTimestamptz(row.SUpdatedAt),
	}

	product := &entities.Product{
		Name:      row.Name,
		Price:     float64FromNumeric(row.Price),
		Seller:    *seller,
		CreatedAt: timeFromTimestamptz(row.CreatedAt),
		UpdatedAt: timeFromTimestamptz(row.UpdatedAt),
	}
	product.Id = row.ID

	return product
}

func fromSqlcProductRowAll(row *db.GetAllProductsRow) *entities.Product {
	seller := &entities.Seller{
		Id:        row.SID,
		Name:      row.SName,
		CreatedAt: timeFromTimestamptz(row.SCreatedAt),
		UpdatedAt: timeFromTimestamptz(row.SUpdatedAt),
	}

	product := &entities.Product{
		Name:      row.Name,
		Price:     float64FromNumeric(row.Price),
		Seller:    *seller,
		CreatedAt: timeFromTimestamptz(row.CreatedAt),
		UpdatedAt: timeFromTimestamptz(row.UpdatedAt),
	}
	product.Id = row.ID

	return product
}
