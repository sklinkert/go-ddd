package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Seller struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

func NewSeller(name string) *Seller {
	return &Seller{
		Id:        uuid.Must(uuid.NewV7()),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
}

func (s *Seller) validate() error {
	if s.Name == "" {
		return fmt.Errorf("%w: name must not be empty", ErrValidation)
	}
	if s.CreatedAt.After(s.UpdatedAt) {
		return fmt.Errorf("%w: created_at must be before updated_at", ErrValidation)
	}

	return nil
}

func (s *Seller) UpdateName(name string) error {
	s.Name = name
	s.UpdatedAt = time.Now()

	return s.validate()
}
