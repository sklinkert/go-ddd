package entities

import (
	"errors"
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

func (s *Seller) validate() error {
	if s.Name == "" {
		return errors.New("invalid seller details")
	}

	return nil
}
