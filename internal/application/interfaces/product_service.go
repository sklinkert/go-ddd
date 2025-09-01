package interfaces

import (
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/query"
)

type ProductService interface {
	CreateProduct(productCommand *command.CreateProductCommand) (*command.CreateProductCommandResult, error)
	UpdateProduct(productCommand *command.UpdateProductCommand) (*command.UpdateProductCommandResult, error)
	DeleteProduct(productCommand *command.DeleteProductCommand) (*command.DeleteProductCommandResult, error)
	FindAllProducts() (*query.GetAllProductsQueryResult, error)
	FindProductById(query *query.GetProductByIdQuery) (*query.GetProductByIdQueryResult, error)
}
