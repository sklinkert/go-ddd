package query

import "github.com/sklinkert/go-ddd/internal/application/common"

type GetAllSellersQuery struct {
	// This query takes no parameters, but we define it for consistency
}

type GetAllSellersQueryResult struct {
	Result []*common.SellerResult
}
