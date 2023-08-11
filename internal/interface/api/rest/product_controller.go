package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/sklinkert/go-ddd/internal/application/services"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"net/http"
	"strconv"
)

type ProductController struct {
	service *services.ProductService
}

func NewProductController(e *echo.Echo, service *services.ProductService) {
	controller := &ProductController{
		service: service,
	}

	e.POST("/products", controller.CreateProduct)
	e.GET("/products", controller.GetAllProducts)
	e.GET("/products/:id", controller.GetProductByID)
}

func (pc *ProductController) CreateProduct(c echo.Context) error {
	product := &entities.Product{}

	if err := c.Bind(product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to parse request body",
		})
	}

	err := pc.service.CreateProduct(product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create product",
		})
	}

	return c.JSON(http.StatusCreated, product)
}

func (pc *ProductController) GetAllProducts(c echo.Context) error {
	products, err := pc.service.GetAllProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch products",
		})
	}

	return c.JSON(http.StatusOK, products)
}

func (pc *ProductController) GetProductByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product ID format",
		})
	}

	product, err := pc.service.FindProductByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch product",
		})
	}

	return c.JSON(http.StatusOK, product)
}
