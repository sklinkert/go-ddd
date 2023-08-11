package entities

import (
	"github.com/google/uuid"
)

type Seller struct {
	ID   uuid.UUID
	Name string
}

func NewSeller(name string) *Seller {
	return &Seller{
		ID:   uuid.New(),
		Name: name,
	}
}
