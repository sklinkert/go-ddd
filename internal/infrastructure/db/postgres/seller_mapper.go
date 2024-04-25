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

// fromDBSeller maps DB persistence model to domain Seller entity.
func fromDBSeller(dbSeller *Seller) *entities.Seller {
	s := &entities.Seller{
		Name: dbSeller.Name,
	}
	s.ID = dbSeller.ID

	return s
}
