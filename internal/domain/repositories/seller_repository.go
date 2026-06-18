package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type SellerRepository interface {
	Create(ctx context.Context, seller *entities.ValidatedSeller) (*entities.Seller, error)
	FindById(ctx context.Context, id uuid.UUID) (*entities.Seller, error)
	FindAll(ctx context.Context) ([]*entities.Seller, error)
	Update(ctx context.Context, seller *entities.ValidatedSeller) (*entities.Seller, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
