package request

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateProductRequest_ToCreateProductCommand(t *testing.T) {
	sellerId := uuid.New()
	req := &CreateProductRequest{
		IdempotencyKey: "key-1",
		Name:           "Widget",
		Price:          9.99,
		SellerId:       sellerId.String(),
	}

	cmd, err := req.ToCreateProductCommand()

	assert.NoError(t, err)
	assert.Equal(t, "key-1", cmd.IdempotencyKey)
	assert.Equal(t, "Widget", cmd.Name)
	assert.Equal(t, 9.99, cmd.Price)
	assert.Equal(t, sellerId, cmd.SellerId)
}

func TestCreateProductRequest_ToCreateProductCommand_InvalidSellerId(t *testing.T) {
	req := &CreateProductRequest{Name: "Widget", Price: 9.99, SellerId: "not-a-uuid"}

	cmd, err := req.ToCreateProductCommand()

	assert.Error(t, err)
	assert.Nil(t, cmd)
}

func TestUpdateProductRequest_ToUpdateProductCommand(t *testing.T) {
	sellerId := uuid.New()
	productId := uuid.New()
	req := &UpdateProductRequest{
		IdempotencyKey: "key-2",
		Name:           "Widget v2",
		Price:          19.99,
		SellerId:       sellerId.String(),
	}

	cmd, err := req.ToUpdateProductCommand(productId)

	assert.NoError(t, err)
	assert.Equal(t, productId, cmd.Id)
	assert.Equal(t, "key-2", cmd.IdempotencyKey)
	assert.Equal(t, "Widget v2", cmd.Name)
	assert.Equal(t, 19.99, cmd.Price)
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
