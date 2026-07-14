package rest_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/common"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProduct_ServiceError(t *testing.T) {
	e := echo.New()
	mockService := new(MockProductService)
	ctrl := rest.NewProductController(e, mockService)

	body, _ := json.Marshal(map[string]any{"name": "X", "price_cents": 100, "currency": "EUR", "seller_id": uuid.NewString()})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("CreateProduct", mock.Anything).Return((*command.CreateProductCommandResult)(nil), errors.New("boom"))

	assert.NoError(t, ctrl.CreateProductController(c))
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}

func TestCreateProduct_InvalidSellerId(t *testing.T) {
	e := echo.New()
	ctrl := rest.NewProductController(e, new(MockProductService))

	body, _ := json.Marshal(map[string]any{"name": "X", "price_cents": 100, "currency": "EUR", "seller_id": "not-a-uuid"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, ctrl.CreateProductController(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetProductById_NotFound(t *testing.T) {
	e := echo.New()
	mockService := new(MockProductService)
	ctrl := rest.NewProductController(e, mockService)

	id := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+id.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	// nil entity + nil error => not found
	mockService.On("FindProductById", mock.Anything).Return((*entities.Product)(nil), nil)

	assert.NoError(t, ctrl.GetProductByIdController(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockService.AssertExpectations(t)
}

func TestGetProductById_InvalidId(t *testing.T) {
	e := echo.New()
	ctrl := rest.NewProductController(e, new(MockProductService))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("not-a-uuid")

	assert.NoError(t, ctrl.GetProductByIdController(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateProduct_Success(t *testing.T) {
	e := echo.New()
	mockService := new(MockProductService)
	ctrl := rest.NewProductController(e, mockService)

	id := uuid.New()
	sellerId := uuid.New()
	body, _ := json.Marshal(map[string]any{"name": "Widget v2", "price_cents": 1999, "currency": "USD", "seller_id": sellerId.String()})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/products/"+id.String(), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	mockService.On("UpdateProduct", mock.MatchedBy(func(cmd *command.UpdateProductCommand) bool {
		return cmd.Id == id && cmd.PriceCents == 1999 && cmd.Currency == entities.USD && cmd.SellerId == sellerId
	})).Return(&command.UpdateProductCommandResult{
		Result: &common.ProductResult{
			Id:       id,
			Name:     "Widget v2",
			Price:    mustMoney(t, 1999, entities.USD),
			SellerId: sellerId,
		},
	}, nil)

	assert.NoError(t, ctrl.UpdateProductController(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	var responseBody map[string]any
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &responseBody))
	assert.Equal(t, "Widget v2", responseBody["name"])
	assert.Equal(t, float64(1999), responseBody["price_cents"])
	assert.Equal(t, "USD", responseBody["currency"])
	assert.Equal(t, sellerId.String(), responseBody["seller_id"])
	mockService.AssertExpectations(t)
}

func TestDeleteProduct_Success(t *testing.T) {
	e := echo.New()
	mockService := new(MockProductService)
	ctrl := rest.NewProductController(e, mockService)

	id := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+id.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	mockService.On("DeleteProduct", mock.Anything).Return(&command.DeleteProductCommandResult{Success: true}, nil)

	assert.NoError(t, ctrl.DeleteProductController(c))
	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockService.AssertExpectations(t)
}
