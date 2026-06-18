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

	body, _ := json.Marshal(map[string]any{"Name": "X", "Price": 1.0, "SellerId": uuid.NewString()})
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

	body, _ := json.Marshal(map[string]any{"Name": "X", "Price": 1.0, "SellerId": "not-a-uuid"})
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
	body, _ := json.Marshal(map[string]any{"Name": "Widget v2", "Price": 19.99, "SellerId": uuid.NewString()})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/products/"+id.String(), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	mockService.On("UpdateProduct", mock.Anything).Return(&command.UpdateProductCommandResult{
		Result: &common.ProductResult{Id: id, Name: "Widget v2", Price: 19.99},
	}, nil)

	assert.NoError(t, ctrl.UpdateProductController(c))
	assert.Equal(t, http.StatusOK, rec.Code)
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
