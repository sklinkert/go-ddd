package interfaces

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type ProductService interface {
	CreateProduct(product *entities.Product) (*entities.ValidatedProduct, error)
	GetAllProducts() ([]*entities.ValidatedProduct, error)
	FindProductByID(id uuid.UUID) (*entities.ValidatedProduct, error)
}
