package postgres

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/testhelpers"
)

func TestSqlcProductRepository_Create_WritesOutboxEvent(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcProductRepository(testDB.Pool)
	validatedSeller := createTestSeller(t, testDB, "Outbox Seller")

	product := entities.NewProduct("Outbox Product", mustMoney(t, 1299, entities.EUR), *validatedSeller)
	validatedProduct, err := entities.NewValidatedProduct(product)
	require.NoError(t, err)

	_, err = repo.Create(context.Background(), validatedProduct)
	require.NoError(t, err)

	// The ProductCreated event must be committed together with the product.
	events, err := testDB.Queries.GetUnpublishedOutboxEvents(context.Background(), 10)
	require.NoError(t, err)
	require.Len(t, events, 1)

	event := events[0]
	assert.Equal(t, "product.created", event.EventName)
	assert.Equal(t, product.Id, event.AggregateID)
	assert.False(t, event.PublishedAt.Valid)

	var payload struct {
		Name       string `json:"Name"`
		PriceCents int64  `json:"PriceCents"`
		Currency   string `json:"Currency"`
	}
	require.NoError(t, json.Unmarshal(event.Payload, &payload))
	assert.Equal(t, "Outbox Product", payload.Name)
	assert.Equal(t, int64(1299), payload.PriceCents)
	assert.Equal(t, "EUR", payload.Currency)

	// Marking published removes it from the unpublished set (relay behavior).
	require.NoError(t, testDB.Queries.MarkOutboxEventPublished(context.Background(), event.ID))
	remaining, err := testDB.Queries.GetUnpublishedOutboxEvents(context.Background(), 10)
	require.NoError(t, err)
	assert.Empty(t, remaining)
}
