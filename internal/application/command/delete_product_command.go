package command

import (
	"github.com/google/uuid"
)

type DeleteProductCommand struct {
	IdempotencyKey string
	Id             uuid.UUID
}

type DeleteProductCommandResult struct {
	Success bool
}
