package common

import (
	"github.com/google/uuid"
	"time"
)

type SellerResult struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
