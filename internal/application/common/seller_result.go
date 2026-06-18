package common

import (
	"time"

	"github.com/google/uuid"
)

type SellerResult struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
