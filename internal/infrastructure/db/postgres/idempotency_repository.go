package postgres

import (
	"context"

	"gorm.io/gorm"

	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/domain/repositories"
)

type idempotencyRepository struct {
	db *gorm.DB
}

func NewIdempotencyRepository(db *gorm.DB) repositories.IdempotencyRepository {
	return &idempotencyRepository{db: db}
}

func (r *idempotencyRepository) FindByKey(ctx context.Context, key string) (*entities.IdempotencyRecord, error) {
	var dbRecord IdempotencyRecord
	result := r.db.WithContext(ctx).Where("key = ?", key).First(&dbRecord)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &entities.IdempotencyRecord{
		ID:         dbRecord.Id,
		Key:        dbRecord.Key,
		Request:    dbRecord.Request,
		Response:   dbRecord.Response,
		StatusCode: dbRecord.StatusCode,
		CreatedAt:  dbRecord.CreatedAt,
	}, nil
}

func (r *idempotencyRepository) Create(ctx context.Context, record *entities.IdempotencyRecord) (*entities.IdempotencyRecord, error) {
	dbRecord := IdempotencyRecord{
		Id:         record.ID,
		Key:        record.Key,
		Request:    record.Request,
		Response:   record.Response,
		StatusCode: record.StatusCode,
		CreatedAt:  record.CreatedAt,
	}

	result := r.db.WithContext(ctx).Create(&dbRecord)
	if result.Error != nil {
		return nil, result.Error
	}

	// Read back the created record
	var createdRecord IdempotencyRecord
	if err := r.db.WithContext(ctx).Where("id = ?", dbRecord.Id).First(&createdRecord).Error; err != nil {
		return nil, err
	}

	return &entities.IdempotencyRecord{
		ID:         createdRecord.Id,
		Key:        createdRecord.Key,
		Request:    createdRecord.Request,
		Response:   createdRecord.Response,
		StatusCode: createdRecord.StatusCode,
		CreatedAt:  createdRecord.CreatedAt,
	}, nil
}

func (r *idempotencyRepository) Update(ctx context.Context, record *entities.IdempotencyRecord) (*entities.IdempotencyRecord, error) {
	dbRecord := IdempotencyRecord{
		Id:         record.ID,
		Key:        record.Key,
		Request:    record.Request,
		Response:   record.Response,
		StatusCode: record.StatusCode,
		CreatedAt:  record.CreatedAt,
	}

	result := r.db.WithContext(ctx).Save(&dbRecord)
	if result.Error != nil {
		return nil, result.Error
	}

	// Read back the updated record
	var updatedRecord IdempotencyRecord
	if err := r.db.WithContext(ctx).Where("id = ?", dbRecord.Id).First(&updatedRecord).Error; err != nil {
		return nil, err
	}

	return &entities.IdempotencyRecord{
		ID:         updatedRecord.Id,
		Key:        updatedRecord.Key,
		Request:    updatedRecord.Request,
		Response:   updatedRecord.Response,
		StatusCode: updatedRecord.StatusCode,
		CreatedAt:  updatedRecord.CreatedAt,
	}, nil
}
