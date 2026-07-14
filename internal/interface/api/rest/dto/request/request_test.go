package request

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProductRequest_ToCreateProductCommand(t *testing.T) {
	sellerId := uuid.New()
	req := &CreateProductRequest{
		IdempotencyKey: "key-1",
		Name:           "Widget",
		PriceCents:     999,
		Currency:       "USD",
		SellerId:       sellerId.String(),
	}

	cmd, err := req.ToCreateProductCommand()

	assert.NoError(t, err)
	assert.Equal(t, "key-1", cmd.IdempotencyKey)
	assert.Equal(t, "Widget", cmd.Name)
	assert.Equal(t, int64(999), cmd.PriceCents)
	assert.Equal(t, entities.USD, cmd.Currency)
	assert.Equal(t, sellerId, cmd.SellerId)
}

func TestCreateProductRequest_ToCreateProductCommand_InvalidSellerId(t *testing.T) {
	req := &CreateProductRequest{Name: "Widget", PriceCents: 999, Currency: "USD", SellerId: "not-a-uuid"}

	cmd, err := req.ToCreateProductCommand()

	assert.Error(t, err)
	assert.Nil(t, cmd)
}

func TestCreateProductRequest_JsonTags(t *testing.T) {
	sellerId := uuid.New()
	body := `{"idempotency_key":"key-1","name":"Widget","price_cents":1234,"currency":"EUR","seller_id":"` + sellerId.String() + `"}`

	var req CreateProductRequest
	require.NoError(t, json.Unmarshal([]byte(body), &req))

	assert.Equal(t, "key-1", req.IdempotencyKey)
	assert.Equal(t, "Widget", req.Name)
	assert.Equal(t, int64(1234), req.PriceCents)
	assert.Equal(t, "EUR", req.Currency)
	assert.Equal(t, sellerId.String(), req.SellerId)
}

func TestUpdateProductRequest_ToUpdateProductCommand(t *testing.T) {
	sellerId := uuid.New()
	productId := uuid.New()
	req := &UpdateProductRequest{
		IdempotencyKey: "key-2",
		Name:           "Widget v2",
		PriceCents:     1999,
		Currency:       "EUR",
		SellerId:       sellerId.String(),
	}

	cmd, err := req.ToUpdateProductCommand(productId)

	assert.NoError(t, err)
	assert.Equal(t, productId, cmd.Id)
	assert.Equal(t, "key-2", cmd.IdempotencyKey)
	assert.Equal(t, "Widget v2", cmd.Name)
	assert.Equal(t, int64(1999), cmd.PriceCents)
	assert.Equal(t, entities.EUR, cmd.Currency)
	assert.Equal(t, sellerId, cmd.SellerId)
}

func TestUpdateProductRequest_ToUpdateProductCommand_InvalidSellerId(t *testing.T) {
	req := &UpdateProductRequest{Name: "Widget", SellerId: "nope"}

	cmd, err := req.ToUpdateProductCommand(uuid.New())

	assert.Error(t, err)
	assert.Nil(t, cmd)
}

func TestCreateSellerRequest_ToCreateSellerCommand(t *testing.T) {
	req := &CreateSellerRequest{IdempotencyKey: "key-3", Name: "Acme"}

	cmd, err := req.ToCreateSellerCommand()

	assert.NoError(t, err)
	assert.Equal(t, "key-3", cmd.IdempotencyKey)
	assert.Equal(t, "Acme", cmd.Name)
}

func TestUpdateSellerRequest_ToUpdateSellerCommand(t *testing.T) {
	id := uuid.New()
	req := &UpdateSellerRequest{IdempotencyKey: "key-4", Id: id, Name: "Acme v2"}

	cmd, err := req.ToUpdateSellerCommand()

	assert.NoError(t, err)
	assert.Equal(t, id, cmd.Id)
	assert.Equal(t, "key-4", cmd.IdempotencyKey)
	assert.Equal(t, "Acme v2", cmd.Name)
}

func TestUpdateSellerRequest_JsonTags(t *testing.T) {
	id := uuid.New()
	body := `{"idempotency_key":"key-4","id":"` + id.String() + `","name":"Acme v2"}`

	var req UpdateSellerRequest
	require.NoError(t, json.Unmarshal([]byte(body), &req))

	assert.Equal(t, id, req.Id)
	assert.Equal(t, "Acme v2", req.Name)
}
