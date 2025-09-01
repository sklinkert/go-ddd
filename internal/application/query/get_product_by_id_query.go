package query

import (
	"github.com/google/uuid"
	"github.com/sklinkert/go-ddd/internal/application/common"
)

type GetProductByIdQuery struct {
	Id uuid.UUID
}

type GetProductByIdQueryResult struct {
	Result *common.ProductResult
}
