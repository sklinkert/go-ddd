package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest/dto/request"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest/dto/response"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSeller(t *testing.T) {
	// Arrange
	mockService := NewMockSellerService()
	controller := rest.NewSellerController(echo.New(), mockService)

	// Create a seller for testing
	seller := entities.NewSeller("TestSeller")

	sellerJSON, _ := json.Marshal(seller)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/sellers", bytes.NewReader(sellerJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Act
	if err := controller.CreateSellerController(c); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("rec: %s\n", rec.Body.String())

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdSeller entities.Seller
	_ = json.Unmarshal(rec.Body.Bytes(), &createdSeller)
	assert.Equal(t, seller.Name, createdSeller.Name)
}

func TestPutSeller(t *testing.T) {
	// Arrange
	mockService := NewMockSellerService()
	controller := rest.NewSellerController(echo.New(), mockService)

	createdSeller, err := mockService.CreateSeller(&command.CreateSellerCommand{Name: "TestSeller"})
	assert.NoError(t, err)

	updateRequest := request.UpdateSellerRequest{
		Id:   createdSeller.Result.Id,
		Name: "updatedName",
	}

	sellerJSON, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/sellers", bytes.NewReader(sellerJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Act
	if err := controller.PutSellerController(c); err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var receivedResponse response.SellerResponse
	err = json.Unmarshal(rec.Body.Bytes(), &receivedResponse)
	assert.NoError(t, err)

	assert.Equal(t, updateRequest.Name, receivedResponse.Name)
}

func TestDeleteSeller(t *testing.T) {
	// Arrange
	mockService := NewMockSellerService()
	controller := rest.NewSellerController(echo.New(), mockService)

	createdSeller, err := mockService.CreateSeller(&command.CreateSellerCommand{Name: "TestSeller"})
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/sellers/%s", createdSeller.Result.Id), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Act
	if err := controller.DeleteSellerController(c); err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestGetSellerById(t *testing.T) {
	// Arrange
	mockService := NewMockSellerService()
	controller := rest.NewSellerController(echo.New(), mockService)

	createdSeller, err := mockService.CreateSeller(&command.CreateSellerCommand{Name: "TestSeller"})
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/sellers/%s", createdSeller.Result.Id), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Act
	if err := controller.GetSellerByIdController(c); err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var fetchedSeller response.SellerResponse
	err = json.Unmarshal(rec.Body.Bytes(), &fetchedSeller)
	assert.NoError(t, err)

	assert.Equal(t, createdSeller.Result.Id.String(), fetchedSeller.Id)
	assert.Equal(t, createdSeller.Result.Name, fetchedSeller.Name)
}

func TestGetAllSellers(t *testing.T) {
	// Arrange
	mockService := NewMockSellerService()
	controller := rest.NewSellerController(echo.New(), mockService)

	_, err := mockService.CreateSeller(&command.CreateSellerCommand{Name: "TestSeller1"})
	assert.NoError(t, err)

	_, err = mockService.CreateSeller(&command.CreateSellerCommand{Name: "TestSeller2"})
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sellers", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Act
	if err := controller.GetAllSellersController(c); err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)

	var sellers response.ListSellersResponse
	err = json.Unmarshal(rec.Body.Bytes(), &sellers)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(sellers.Sellers))
}
