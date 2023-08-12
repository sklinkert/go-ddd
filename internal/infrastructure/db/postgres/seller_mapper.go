package postgres

import (
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

// ToDBSeller maps domain Seller entity to DB persistence model.
func ToDBSeller(seller *entities.ValidatedSeller) *Seller {
	s := &Seller{
		Name: seller.Name,
	}
	s.ID = seller.ID

	return s
}

// FromDBSeller maps DB persistence model to domain Seller entity.
func FromDBSeller(dbSeller *Seller) (*entities.ValidatedSeller, error) {
	s := &entities.Seller{
		ID:   dbSeller.ID,
		Name: dbSeller.Name,
	}
	s.ID = dbSeller.ID

	return entities.NewValidatedSeller(s)
}
