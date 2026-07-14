package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest/dto/response"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/labstack/echo/v4"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest"
	"github.com/stretchr/testify/assert"
)

func mustMoney(t *testing.T, cents int64, currency entities.Currency) entities.Money {
	t.Helper()
	money, err := entities.NewMoney(cents, currency)
	require.NoError(t, err)
	return money
}

func TestCreateProduct(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := new(MockProductService)
	sellerId := "123e4567-e89b-12d3-a456-426614174000"
	reqBody := map[string]interface{}{
		"name":            "TestProduct",
		"price_cents":     999,
		"currency":        "USD",
		"seller_id":       sellerId,
		"idempotency_key": "idem-123",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(reqBodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctrl := rest.NewProductController(e, mockService)

	mockService.On("CreateProduct", mock.MatchedBy(func(cmd *command.CreateProductCommand) bool {
		return cmd.Name == "TestProduct" &&
			cmd.PriceCents == 999 &&
			cmd.Currency == entities.USD &&
			cmd.SellerId.String() == sellerId &&
			cmd.IdempotencyKey == "idem-123"
	})).Return((*command.CreateProductCommandResult)(nil), nil)

	// Execute
	err := ctrl.CreateProductController(c)
	assert.NoError(t, err)

	var responseBody map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &responseBody))

	// Assertions
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "TestProduct", responseBody["name"])
	assert.Equal(t, float64(999), responseBody["price_cents"])
	assert.Equal(t, "USD", responseBody["currency"])
	assert.Equal(t, sellerId, responseBody["seller_id"])
	assert.NotEmpty(t, responseBody["id"])
	mockService.AssertExpectations(t)
}

func TestGetAllProducts(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := new(MockProductService)

	sellerId := uuid.New()
	expectedProducts := []*entities.Product{
		{
			Id:       uuid.New(),
			Name:     "TestProduct1",
			Price:    mustMoney(t, 999, entities.USD),
			SellerId: sellerId,
		}, {
			Id:       uuid.New(),
			Name:     "TestProduct2",
			Price:    mustMoney(t, 1499, entities.EUR),
			SellerId: sellerId,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ctrl := rest.NewProductController(e, mockService)
	mockService.On("FindAllProducts").Return(expectedProducts, nil)

	var expectedListResponse response.ListProductsResponse
	for _, product := range expectedProducts {
		expectedListResponse.Products = append(expectedListResponse.Products,
			&response.ProductResponse{
				Id:         product.Id.String(),
				Name:       product.Name,
				PriceCents: product.Price.Cents(),
				Currency:   string(product.Price.Currency()),
				SellerId:   product.SellerId.String(),
			})
	}

	// Assertions
	if assert.NoError(t, ctrl.GetAllProductsController(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var receivedListResponse response.ListProductsResponse
		err := json.Unmarshal(rec.Body.Bytes(), &receivedListResponse)
		if assert.NoError(t, err) {
			assert.ElementsMatch(t, expectedListResponse.Products, receivedListResponse.Products)
		}
	}
}

func TestGetAllProducts_EmptyReturnsEmptyArray(t *testing.T) {
	e := echo.New()
	mockService := new(MockProductService)
	ctrl := rest.NewProductController(e, mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("FindAllProducts").Return([]*entities.Product{}, nil)

	assert.NoError(t, ctrl.GetAllProductsController(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"products":[]}`, rec.Body.String())
}
