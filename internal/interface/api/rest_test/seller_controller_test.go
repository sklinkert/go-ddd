package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest"
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
	req := httptest.NewRequest(http.MethodPost, "/sellers", bytes.NewReader(sellerJSON))
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
