package query

import "github.com/sklinkert/go-ddd/internal/application/common"

type GetAllProductsQueryResult struct {
	Result []*common.ProductResult
}
