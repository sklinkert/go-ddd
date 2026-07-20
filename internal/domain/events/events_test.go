package events

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewProductCreated(t *testing.T) {
	productId := uuid.New()
	sellerId := uuid.New()

	event := NewProductCreated(productId, "Widget", 999, "USD", sellerId)

	assert.Equal(t, "product.created", event.EventName())
	assert.Equal(t, productId, event.AggregateId())
	assert.Equal(t, sellerId, event.SellerId)
	assert.Equal(t, int64(999), event.PriceMinorUnits)
	assert.Equal(t, "USD", event.Currency)
	assert.NotEqual(t, uuid.Nil, event.EventId())
	assert.WithinDuration(t, time.Now(), event.OccurredAt(), time.Second)
}

func TestBaseEvent_UniqueIds(t *testing.T) {
	aggregateId := uuid.New()

	first := NewBaseEvent(aggregateId)
	second := NewBaseEvent(aggregateId)

	assert.NotEqual(t, first.EventId(), second.EventId())
}
