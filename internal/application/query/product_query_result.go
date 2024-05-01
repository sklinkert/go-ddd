package query

import "github.com/sklinkert/go-ddd/internal/application/common"

type ProductQueryResult struct {
	Result *common.ProductResult
}

type ProductQueryListResult struct {
	Result []*common.ProductResult
}
