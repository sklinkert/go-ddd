package postgres

import (
	"context"

	"github.com/google/uuid"

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

func (repo *SqlcProductRepository) Create(product *entities.ValidatedProduct) (*entities.Product, error) {
	ctx := context.Background()

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

	return repo.FindById(createdProduct.ID)
}

func (repo *SqlcProductRepository) FindById(id uuid.UUID) (*entities.Product, error) {
	ctx := context.Background()

	row, err := repo.queries.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}

	return fromSqlcProductRow(&row), nil
}

func (repo *SqlcProductRepository) FindAll() ([]*entities.Product, error) {
	ctx := context.Background()

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

func (repo *SqlcProductRepository) Update(product *entities.ValidatedProduct) (*entities.Product, error) {
	ctx := context.Background()

	err := repo.queries.UpdateProduct(ctx, db.UpdateProductParams{
		ID:        product.Id,
		Name:      product.Name,
		Price:     numericFromFloat64(product.Price),
		SellerID:  product.Seller.Id,
		UpdatedAt: timestamptzFromTime(product.UpdatedAt),
	})
	if err != nil {
		return nil, err
	}

	return repo.FindById(product.Id)
}

func (repo *SqlcProductRepository) Delete(id uuid.UUID) error {
	ctx := context.Background()
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
