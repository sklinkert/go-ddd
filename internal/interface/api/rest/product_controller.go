package rest

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sklinkert/go-ddd/internal/application/command"
	"github.com/sklinkert/go-ddd/internal/application/interfaces"
	"net/http"
)

type ProductController struct {
	service interfaces.ProductService
}

func NewProductController(e *echo.Echo, service interfaces.ProductService) *ProductController {
	controller := &ProductController{
		service: service,
	}

	e.POST("/products", controller.CreateProduct)
	e.GET("/products", controller.GetAllProducts)
	e.GET("/products/:id", controller.GetProductByID)

	return controller
}

func (pc *ProductController) CreateProduct(c echo.Context) error {
	var createProductRequest CreateProductRequest

	if err := c.Bind(&createProductRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to parse request body",
		})
	}

	productCommand, err := getCreateProductCommand(&createProductRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product ID format",
		})
	}

	result, err := pc.service.CreateProduct(productCommand)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create product",
		})
	}

	return c.JSON(http.StatusCreated, result.Result)
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
	id, err := uuid.Parse(c.Param("id"))
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

func getCreateProductCommand(createProductRequest *CreateProductRequest) (*command.CreateProductCommand, error) {
	sellerId, err := uuid.Parse(createProductRequest.SellerID)
	if err != nil {
		return nil, err
	}

	return &command.CreateProductCommand{
		Name:     createProductRequest.Name,
		Price:    createProductRequest.Price,
		SellerID: sellerId,
	}, nil
}
