package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entities.ValidatedProduct) (*entities.Product, error)
	FindById(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	FindAll(ctx context.Context) ([]*entities.Product, error)
	Update(ctx context.Context, product *entities.ValidatedProduct) (*entities.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
