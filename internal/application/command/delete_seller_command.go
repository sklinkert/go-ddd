package command

import (
	"github.com/google/uuid"
)

type DeleteSellerCommand struct {
	IdempotencyKey string
	Id             uuid.UUID
}

type DeleteSellerCommandResult struct {
	Success bool
}
