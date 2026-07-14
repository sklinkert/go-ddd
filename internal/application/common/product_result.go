package common

import (
	"time"

	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

type ProductResult struct {
	Id        uuid.UUID
	Name      string
	Price     entities.Money
	SellerId  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
