package query

import "github.com/sklinkert/go-ddd/internal/application/common"

type SellerQueryResult struct {
	Result *common.SellerResult
}

type SellerQueryListResult struct {
	Result []*common.SellerResult
}
