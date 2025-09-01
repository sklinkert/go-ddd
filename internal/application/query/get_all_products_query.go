package query

import "github.com/sklinkert/go-ddd/internal/application/common"

type GetAllProductsQuery struct {
	// This query takes no parameters, but we define it for consistency
}

type GetAllProductsQueryResult struct {
	Result []*common.ProductResult
}
