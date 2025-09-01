package entities

import (
	"time"

	"github.com/google/uuid"
)

type IdempotencyRecord struct {
	ID         uuid.UUID
	Key        string
	Request    string
	Response   string
	StatusCode int
	CreatedAt  time.Time
}

func NewIdempotencyRecord(key string, request string) *IdempotencyRecord {
	return &IdempotencyRecord{
		ID:        uuid.New(),
		Key:       key,
		Request:   request,
		CreatedAt: time.Now(),
	}
}

func (i *IdempotencyRecord) SetResponse(response string, statusCode int) {
	i.Response = response
	i.StatusCode = statusCode
}
