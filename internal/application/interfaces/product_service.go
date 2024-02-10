package interfaces

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type ProductService interface {
	CreateProduct(productCommand *command.CreateProductCommand) (*command.CreateProductCommandResult, error)
	GetAllProducts() ([]*entities.ValidatedProduct, error)
	FindProductById(id uuid.UUID) (*entities.ValidatedProduct, error)
}
