package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/events"
)

type Product struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Price     Money
	SellerId  uuid.UUID

	domainEvents []events.DomainEvent
}

func (p *Product) validate() error {
	if p.Name == "" {
		return fmt.Errorf("%w: name must not be empty", ErrValidation)
	}
	if p.Price.Cents() == 0 {
		return fmt.Errorf("%w: price must be greater than 0", ErrValidation)
	}
	if p.SellerId == uuid.Nil {
		return fmt.Errorf("%w: seller id must not be empty", ErrValidation)
	}
	if p.CreatedAt.After(p.UpdatedAt) {
		return fmt.Errorf("%w: created_at must be before updated_at", ErrValidation)
	}

	return nil
}

// NewProduct requires a ValidatedSeller so a product can only ever be
// created against a seller that passed validation. The product stores just
// the seller's Id: sellers are a separate aggregate and must not be embedded.
func NewProduct(name string, price Money, seller ValidatedSeller) *Product {
	product := &Product{
		Id:        uuid.Must(uuid.NewV7()),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Price:     price,
		SellerId:  seller.Id,
	}

	product.recordEvent(events.NewProductCreated(product.Id, name, price.Cents(), string(price.Currency()), seller.Id))

	return product
}

func (p *Product) recordEvent(event events.DomainEvent) {
	p.domainEvents = append(p.domainEvents, event)
}

// PullEvents returns the recorded domain events and clears them. The
// repository persists them in the same transaction as the aggregate
// (transactional outbox), so callers pull exactly once per save.
func (p *Product) PullEvents() []events.DomainEvent {
	pulled := p.domainEvents
	p.domainEvents = nil
	return pulled
}

func (p *Product) UpdateName(name string) error {
	p.Name = name
	p.UpdatedAt = time.Now()

	return p.validate()
}

func (p *Product) UpdatePrice(price Money) error {
	p.Price = price
	p.UpdatedAt = time.Now()

	return p.validate()
}

// AssignSeller moves the product to a different (validated) seller.
func (p *Product) AssignSeller(seller ValidatedSeller) error {
	p.SellerId = seller.Id
	p.UpdatedAt = time.Now()

	return p.validate()
}
