package events

import "github.com/google/uuid"

const ProductCreatedEventName = "product.created"

type ProductCreated struct {
	BaseEvent
	Name       string
	PriceCents int64
	Currency   string
	SellerId   uuid.UUID
}

func NewProductCreated(productId uuid.UUID, name string, priceCents int64, currency string, sellerId uuid.UUID) ProductCreated {
	return ProductCreated{
		BaseEvent:  NewBaseEvent(productId),
		Name:       name,
		PriceCents: priceCents,
		Currency:   currency,
		SellerId:   sellerId,
	}
}

func (e ProductCreated) EventName() string { return ProductCreatedEventName }
