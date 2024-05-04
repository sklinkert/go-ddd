package entities

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Seller struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSeller(name string) *Seller {
	return &Seller{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (s *Seller) validate() error {
	if s.Name == "" {
		return errors.New("name must not be empty")
	}
	if s.CreatedAt.After(s.UpdatedAt) {
		return errors.New("created_at must be before updated_at")
	}

	return nil
}

func (s *Seller) UpdateName(name string) error {
	s.Name = name
	s.UpdatedAt = time.Now()

	return s.validate()
}
