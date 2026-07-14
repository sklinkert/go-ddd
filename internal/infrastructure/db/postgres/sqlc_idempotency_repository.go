package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/domain/repositories"
	db "github.com/sklinkert/go-ddd/internal/infrastructure/db/sqlc"
)

type SqlcIdempotencyRepository struct {
	queries *db.Queries
}

func NewSqlcIdempotencyRepository(queries *db.Queries) repositories.IdempotencyRepository {
	return &SqlcIdempotencyRepository{queries: queries}
}

func (r *SqlcIdempotencyRepository) Reserve(ctx context.Context, record *entities.IdempotencyRecord) (bool, error) {
	rows, err := r.queries.ReserveIdempotencyKey(ctx, db.ReserveIdempotencyKeyParams{
		ID:        record.Id,
		Key:       record.Key,
		Request:   record.Request,
		CreatedAt: timestamptzFromTime(record.CreatedAt),
	})
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (r *SqlcIdempotencyRepository) FindByKey(ctx context.Context, key string) (*entities.IdempotencyRecord, error) {
	dbRecord, err := r.queries.GetIdempotencyRecordByKey(ctx, key)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entities.IdempotencyRecord{
		Id:         dbRecord.ID,
		Key:        dbRecord.Key,
		Request:    dbRecord.Request,
		Response:   dbRecord.Response,
		StatusCode: int(dbRecord.StatusCode),
		CreatedAt:  timeFromTimestamptz(dbRecord.CreatedAt),
	}, nil
}

func (r *SqlcIdempotencyRepository) SetResponse(ctx context.Context, key string, response string, statusCode int) error {
	return r.queries.SetIdempotencyResponse(ctx, db.SetIdempotencyResponseParams{
		Key:        key,
		Response:   response,
		StatusCode: int32(statusCode),
	})
}

func (r *SqlcIdempotencyRepository) Delete(ctx context.Context, key string) error {
	return r.queries.DeleteIdempotencyRecord(ctx, key)
}
