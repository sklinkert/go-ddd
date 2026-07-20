package mapper

import (
	"github.com/sklinkert/go-ddd/internal/application/common"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest/dto/response"
)

func ToProductResponse(product *common.ProductResult) *response.ProductResponse {
	return &response.ProductResponse{
		Id:              product.Id.String(),
		Name:            product.Name,
		PriceMinorUnits: product.Price.MinorUnits(),
		Currency:        string(product.Price.Currency()),
		SellerId:        product.SellerId.String(),
		CreatedAt:       product.CreatedAt,
		UpdatedAt:       product.UpdatedAt,
	}
}

func ToProductListResponse(products []*common.ProductResult) *response.ListProductsResponse {
	responseList := make([]*response.ProductResponse, 0, len(products))
	for _, product := range products {
		responseList = append(responseList, ToProductResponse(product))
	}
	return &response.ListProductsResponse{Products: responseList}
}
