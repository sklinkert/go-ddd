package rest

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/interfaces"
	"github.com/sklinkert/go-ddd/internal/application/query"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest/dto/mapper"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest/dto/request"
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
	e.PUT("/api/v1/products/:id", controller.UpdateProductController)
	e.DELETE("/api/v1/products/:id", controller.DeleteProductController)

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

	result, err := pc.service.CreateProduct(c.Request().Context(), productCommand)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create product",
		})
	}

	response := mapper.ToProductResponse(result.Result)

	return c.JSON(http.StatusCreated, response)
}

func (pc *ProductController) GetAllProductsController(c echo.Context) error {
	products, err := pc.service.FindAllProducts(c.Request().Context())
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

	product, err := pc.service.FindProductById(c.Request().Context(), &query.GetProductByIdQuery{Id: id})
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

func (pc *ProductController) UpdateProductController(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product Id format",
		})
	}

	var updateProductRequest request.UpdateProductRequest
	if err := c.Bind(&updateProductRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to parse request body",
		})
	}

	productCommand, err := updateProductRequest.ToUpdateProductCommand(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid seller Id format",
		})
	}

	result, err := pc.service.UpdateProduct(c.Request().Context(), productCommand)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update product",
		})
	}

	response := mapper.ToProductResponse(result.Result)

	return c.JSON(http.StatusOK, response)
}

func (pc *ProductController) DeleteProductController(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product Id format",
		})
	}

	_, err = pc.service.DeleteProduct(c.Request().Context(), &command.DeleteProductCommand{Id: id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete product",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
