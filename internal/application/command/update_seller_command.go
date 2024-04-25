package command

import "github.com/google/uuid"

type UpdateSellerCommand struct {
	// TODO: Implement idempotency key

	ID   uuid.UUID
	Name string
}
