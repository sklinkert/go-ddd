package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sklinkert/go-ddd/internal/application/services"
	"github.com/sklinkert/go-ddd/internal/infrastructure/rest"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := new(services.MockProductService) // Assuming you have a mock of ProductService
	reqBody := map[string]interface{}{"name": "TestProduct", "price": 9.99}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(reqBodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ctrl := &rest.ProductController{
		Service: mockService,
	}

	// You may want to setup some expected behavior on the mock service if required
	// For example:
	// mockService.On("CreateProduct", mock.Anything).Return(nil)

	// Assertions
	if assert.NoError(t, ctrl.CreateProduct(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		// further assertions on the response, for example checking the returned JSON
	}
}

// Similar tests can be written for GetAllProducts and GetProductByID.
// Just setup the mock service to return appropriate data and make the required call.

// Note: The mock services would be responsible for mocking the behavior of the services layer,
// eliminating the need to connect to an actual database or external services.
