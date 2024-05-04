package mapper

import (
	"github.com/sklinkert/go-ddd/internal/application/common"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest/dto/response"
)

func ToSellerResponse(product *common.SellerResult) *response.SellerResponse {
	return &response.SellerResponse{
		Id:        product.Id.String(),
		Name:      product.Name,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}

func ToSellerListResponse(products []*common.SellerResult) *response.ListSellersResponse {
	var responseList []*response.SellerResponse

	for _, product := range products {
		responseList = append(responseList, ToSellerResponse(product))
	}

	return &response.ListSellersResponse{Sellers: responseList}
}
