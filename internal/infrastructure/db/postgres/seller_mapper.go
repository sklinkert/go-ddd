package postgres

import (
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

// toDBSeller maps domain Seller entity to DB persistence model.
func toDBSeller(seller *entities.ValidatedSeller) *Seller {
	s := &Seller{
		Name: seller.Name,
	}
	s.Id = seller.Id

	return s
}

// fromDBSeller maps DB persistence model to domain Seller entity.
func fromDBSeller(dbSeller *Seller) *entities.Seller {
	s := &entities.Seller{
		Name: dbSeller.Name,
	}
	s.Id = dbSeller.Id

	return s
}
