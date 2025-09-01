package postgres

import (
	"context"
	"database/sql"
	"errors"

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

func (r *SqlcIdempotencyRepository) FindByKey(ctx context.Context, key string) (*entities.IdempotencyRecord, error) {
	dbRecord, err := r.queries.GetIdempotencyRecordByKey(ctx, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entities.IdempotencyRecord{
		ID:         dbRecord.ID,
		Key:        dbRecord.Key,
		Request:    dbRecord.Request,
		Response:   dbRecord.Response,
		StatusCode: int(dbRecord.StatusCode),
		CreatedAt:  timeFromTimestamptz(dbRecord.CreatedAt),
	}, nil
}

func (r *SqlcIdempotencyRepository) Create(ctx context.Context, record *entities.IdempotencyRecord) (*entities.IdempotencyRecord, error) {
	createdRecord, err := r.queries.CreateIdempotencyRecord(ctx, db.CreateIdempotencyRecordParams{
		ID:         record.ID,
		Key:        record.Key,
		Request:    record.Request,
		Response:   record.Response,
		StatusCode: int32(record.StatusCode),
		CreatedAt:  timestamptzFromTime(record.CreatedAt),
	})
	if err != nil {
		return nil, err
	}

	return &entities.IdempotencyRecord{
		ID:         createdRecord.ID,
		Key:        createdRecord.Key,
		Request:    createdRecord.Request,
		Response:   createdRecord.Response,
		StatusCode: int(createdRecord.StatusCode),
		CreatedAt:  timeFromTimestamptz(createdRecord.CreatedAt),
	}, nil
}

func (r *SqlcIdempotencyRepository) Update(ctx context.Context, record *entities.IdempotencyRecord) (*entities.IdempotencyRecord, error) {
	updatedRecord, err := r.queries.UpdateIdempotencyRecord(ctx, db.UpdateIdempotencyRecordParams{
		ID:         record.ID,
		Request:    record.Request,
		Response:   record.Response,
		StatusCode: int32(record.StatusCode),
	})
	if err != nil {
		return nil, err
	}

	return &entities.IdempotencyRecord{
		ID:         updatedRecord.ID,
		Key:        updatedRecord.Key,
		Request:    updatedRecord.Request,
		Response:   updatedRecord.Response,
		StatusCode: int(updatedRecord.StatusCode),
		CreatedAt:  timeFromTimestamptz(updatedRecord.CreatedAt),
	}, nil
}
