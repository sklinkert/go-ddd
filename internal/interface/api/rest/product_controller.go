package rest

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sklinkert/go-ddd/internal/application/interfaces"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest/dto/mapper"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest/dto/request"
	"net/http"
)

type ProductController struct {
	service interfaces.ProductService
}

func NewProductController(e *echo.Echo, service interfaces.ProductService) *ProductController {
	controller := &ProductController{
		service: service,
	}

	e.POST("/api/v1/products", controller.CreateProductController)
	e.GET("/api/v1/products", controller.GetAllProductsController)
	e.GET("/api/v1/products/:id", controller.GetProductByIdController)
	e.Use(middleware.Recover())

	return controller
}

func (pc *ProductController) CreateProductController(c echo.Context) error {
	var createProductRequest request.CreateProductRequest

	if err := c.Bind(&createProductRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to parse request body",
		})
	}

	productCommand, err := createProductRequest.ToCreateProductCommand()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product Id format",
		})
	}

	result, err := pc.service.CreateProduct(productCommand)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create product",
		})
	}

	response := mapper.ToProductResponse(result.Result)

	return c.JSON(http.StatusCreated, response)
}

func (pc *ProductController) GetAllProductsController(c echo.Context) error {
	products, err := pc.service.FindAllProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch products",
		})
	}

	response := mapper.ToProductListResponse(products.Result)

	return c.JSON(http.StatusOK, response)
}

func (pc *ProductController) GetProductByIdController(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product Id format",
		})
	}

	product, err := pc.service.FindProductById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch product",
		})
	}

	if product == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Product not found",
		})
	}

	response := mapper.ToProductResponse(product.Result)

	return c.JSON(http.StatusOK, response)
}
