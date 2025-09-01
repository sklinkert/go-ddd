package query

import "github.com/sklinkert/go-ddd/internal/application/common"

type GetAllSellersQueryResult struct {
	Result []*common.SellerResult
}
