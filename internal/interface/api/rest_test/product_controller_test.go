package rest_test

import (
	"bytes"
	"encoding/json"
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
	reqBody := map[string]interface{}{"Name": "TestProduct", "Price": 9.99}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(reqBodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctrl := rest.NewProductController(e, mockService)
	mockService.On("CreateProduct", mock.Anything).Return(nil)

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
	delete(responseBody, "ID")
	delete(responseBody, "Seller")

	// Assertions
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, reqBody, responseBody)
	mockService.AssertExpectations(t)
}
