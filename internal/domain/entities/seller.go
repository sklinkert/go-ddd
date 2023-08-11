package entities

import (
	"errors"
	"github.com/google/uuid"
)

type Seller struct {
	ID   uuid.UUID
	Name string
}

type ValidatedSeller struct {
	Seller
	isValidated bool
}

func (s *Seller) validate() error {
	if s.Name == "" {
		return errors.New("invalid seller details")
	}

	return nil
}

func (vp *ValidatedSeller) IsValid() bool {
	return vp.isValidated
}

func NewSeller(name string) *Seller {
	return &Seller{
		ID:   uuid.New(),
		Name: name,
	}
}

func NewValidatedSeller(seller *Seller) (*ValidatedSeller, error) {
	if err := seller.validate(); err != nil {
		return nil, err
	}

	return &ValidatedSeller{
		Seller:      *seller,
		isValidated: true,
	}, nil
}
