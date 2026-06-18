package interfaces

import (
	"context"

	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/query"
)

type ProductService interface {
	CreateProduct(ctx context.Context, productCommand *command.CreateProductCommand) (*command.CreateProductCommandResult, error)
	UpdateProduct(ctx context.Context, productCommand *command.UpdateProductCommand) (*command.UpdateProductCommandResult, error)
	DeleteProduct(ctx context.Context, productCommand *command.DeleteProductCommand) (*command.DeleteProductCommandResult, error)
	FindAllProducts(ctx context.Context) (*query.GetAllProductsQueryResult, error)
	FindProductById(ctx context.Context, query *query.GetProductByIdQuery) (*query.GetProductByIdQueryResult, error)
}
