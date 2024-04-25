package rest_test

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := new(MockProductService)
	reqBody := map[string]interface{}{"Name": "TestProduct", "Price": 9.99, "SellerId": "123e4567-e89b-12d3-a456-426614174000"}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(reqBodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctrl := rest.NewProductController(e, mockService)

	createProductCommandResult := &command.CreateProductCommandResult{
		Result: command.ProductResult{
			Id:    uuid.New(),
			Name:  "TestProduct",
			Price: 9.99,
		},
	}
	mockService.On("CreateProduct", mock.Anything).Return(createProductCommandResult, nil)

	// Execute
	err := ctrl.CreateProduct(c)
	assert.NoError(t, err)

	// Deserialize the response body
	var responseBody map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Remove fields from responseBody that are not present in reqBody
	// For example, remove ID and Seller fields
	delete(responseBody, "Id")
	delete(responseBody, "Seller")
	delete(reqBody, "SellerId")

	// Assertions
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, reqBody, responseBody)
	mockService.AssertExpectations(t)
}

func TestGetAllProducts(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := new(MockProductService) // Assuming you have a mock of ProductService

	expectedProducts := []*entities.Product{
		{
			ID:    uuid.New(),
			Name:  "TestProduct1",
			Price: 9.99,
		}, {
			ID:    uuid.New(),
			Name:  "TestProduct2",
			Price: 14.99,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/expectedProducts", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ctrl := rest.NewProductController(e, mockService)
	mockService.On("FindAllProducts").Return(expectedProducts, nil)

	// Assertions
	if assert.NoError(t, ctrl.GetAllProducts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var responseProducts []*entities.Product
		err := json.Unmarshal(rec.Body.Bytes(), &responseProducts)
		if assert.NoError(t, err) {
			assert.ElementsMatch(t, expectedProducts, responseProducts)
		}
	}
}
