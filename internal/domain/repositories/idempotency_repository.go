package repositories

import (
	"context"

	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type IdempotencyRepository interface {
	// Reserve atomically claims the record's key. It returns false when the
	// key is already claimed by another request.
	Reserve(ctx context.Context, record *entities.IdempotencyRecord) (bool, error)
	FindByKey(ctx context.Context, key string) (*entities.IdempotencyRecord, error)
	SetResponse(ctx context.Context, key string, response string, statusCode int) error
	// Delete releases a reserved key, e.g. when the operation failed and the
	// client should be able to retry.
	Delete(ctx context.Context, key string) error
}
