package mapper

import (
	"github.com/sklinkert/go-ddd/internal/application/common"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest/dto/response"
)

func ToProductResponse(product *common.ProductResult) *response.ProductResponse {
	return &response.ProductResponse{
		Id:        product.Id.String(),
		Name:      product.Name,
		Price:     product.Price,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

func ToProductListResponse(products []*common.ProductResult) *response.ListProductsResponse {
	var responseList []*response.ProductResponse
	for _, product := range products {
		responseList = append(responseList, ToProductResponse(product))
	}
	return &response.ListProductsResponse{Products: responseList}
}
