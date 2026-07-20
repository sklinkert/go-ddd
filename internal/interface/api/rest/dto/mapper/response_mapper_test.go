package mapper

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/common"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustMoney(t *testing.T, minorUnits int64, currency entities.Currency) entities.Money {
	t.Helper()
	money, err := entities.NewMoney(minorUnits, currency)
	require.NoError(t, err)
	return money
}

func TestToSellerResponse(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	result := &common.SellerResult{Id: id, Name: "Acme", CreatedAt: now, UpdatedAt: now}

	resp := ToSellerResponse(result)

	assert.Equal(t, id.String(), resp.Id)
	assert.Equal(t, "Acme", resp.Name)
	assert.Equal(t, now, resp.CreatedAt)
}

func TestToSellerListResponse(t *testing.T) {
	results := []*common.SellerResult{
		{Id: uuid.New(), Name: "Acme"},
		{Id: uuid.New(), Name: "Globex"},
	}

	resp := ToSellerListResponse(results)

	assert.Len(t, resp.Sellers, 2)
	assert.Equal(t, "Acme", resp.Sellers[0].Name)
	assert.Equal(t, "Globex", resp.Sellers[1].Name)
}

func TestToSellerListResponse_Empty(t *testing.T) {
	resp := ToSellerListResponse(nil)
	assert.Empty(t, resp.Sellers)
}

func TestToProductResponse(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	sellerId := uuid.New()
	result := &common.ProductResult{
		Id:        id,
		Name:      "Widget",
		Price:     mustMoney(t, 999, entities.USD),
		SellerId:  sellerId,
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := ToProductResponse(result)

	assert.Equal(t, id.String(), resp.Id)
	assert.Equal(t, "Widget", resp.Name)
	assert.Equal(t, int64(999), resp.PriceMinorUnits)
	assert.Equal(t, "USD", resp.Currency)
	assert.Equal(t, sellerId.String(), resp.SellerId)
	assert.Equal(t, now, resp.CreatedAt)
}

func TestToProductResponse_JsonShape(t *testing.T) {
	result := &common.ProductResult{
		Id:       uuid.New(),
		Name:     "Widget",
		Price:    mustMoney(t, 1234, entities.EUR),
		SellerId: uuid.New(),
	}

	data, err := json.Marshal(ToProductResponse(result))
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(data, &payload))

	assert.Equal(t, float64(1234), payload["price_minor_units"])
	assert.Equal(t, "EUR", payload["currency"])
	assert.Equal(t, result.SellerId.String(), payload["seller_id"])
}

func TestToProductListResponse(t *testing.T) {
	results := []*common.ProductResult{
		{Id: uuid.New(), Name: "Widget", Price: mustMoney(t, 100, entities.USD), SellerId: uuid.New()},
		{Id: uuid.New(), Name: "Gadget", Price: mustMoney(t, 200, entities.EUR), SellerId: uuid.New()},
	}

	resp := ToProductListResponse(results)

	assert.Len(t, resp.Products, 2)
	assert.Equal(t, "Widget", resp.Products[0].Name)
	assert.Equal(t, int64(100), resp.Products[0].PriceMinorUnits)
	assert.Equal(t, "Gadget", resp.Products[1].Name)
	assert.Equal(t, "EUR", resp.Products[1].Currency)
}

func TestToProductListResponse_Empty(t *testing.T) {
	resp := ToProductListResponse(nil)

	assert.NotNil(t, resp.Products)
	assert.Empty(t, resp.Products)
}
