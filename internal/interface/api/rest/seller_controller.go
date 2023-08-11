package rest

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sklinkert/go-ddd/internal/application/services"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"net/http"
)

type SellerController struct {
	service *services.SellerService
}

func NewSellerController(e *echo.Echo, service *services.SellerService) {
	controller := &SellerController{
		service: service,
	}

	e.POST("/sellers", controller.CreateSeller)
	e.GET("/sellers", controller.GetAllSellers)
	e.GET("/sellers/:id", controller.GetSellerByID)
}

func (sc *SellerController) CreateSeller(c echo.Context) error {
	seller := &entities.Seller{}

	if err := c.Bind(seller); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to parse request body",
		})
	}

	err := sc.service.CreateSeller(seller)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create seller",
		})
	}

	return c.JSON(http.StatusCreated, seller)
}

func (sc *SellerController) GetAllSellers(c echo.Context) error {
	sellers, err := sc.service.GetAllSellers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch sellers",
		})
	}

	return c.JSON(http.StatusOK, sellers)
}

func (sc *SellerController) GetSellerByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid seller ID format",
		})
	}

	seller, err := sc.service.GetSellerByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch seller",
		})
	}

	return c.JSON(http.StatusOK, seller)
}
