package query

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/common"
)

type GetSellerByIdQuery struct {
	Id uuid.UUID
}

type GetSellerByIdQueryResult struct {
	Result *common.SellerResult
}
