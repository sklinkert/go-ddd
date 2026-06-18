package mapper

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/common"
	"github.com/stretchr/testify/assert"
)

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
	result := &common.ProductResult{Id: id, Name: "Widget", Price: 9.99, CreatedAt: now, UpdatedAt: now}

	resp := ToProductResponse(result)

	assert.Equal(t, id.String(), resp.Id)
	assert.Equal(t, "Widget", resp.Name)
	assert.Equal(t, 9.99, resp.Price)
	assert.Equal(t, now, resp.CreatedAt)
}

func TestToProductListResponse(t *testing.T) {
	results := []*common.ProductResult{
		{Id: uuid.New(), Name: "Widget", Price: 1.0},
		{Id: uuid.New(), Name: "Gadget", Price: 2.0},
	}

	resp := ToProductListResponse(results)

	assert.Len(t, resp.Products, 2)
	assert.Equal(t, "Widget", resp.Products[0].Name)
	assert.Equal(t, "Gadget", resp.Products[1].Name)
}

func TestToProductListResponse_Empty(t *testing.T) {
	resp := ToProductListResponse(nil)
	assert.Empty(t, resp.Products)
}
