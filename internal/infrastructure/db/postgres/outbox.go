package postgres

import (
	"context"
	"encoding/json"

	"github.com/sklinkert/go-ddd/internal/domain/events"
	db "github.com/sklinkert/go-ddd/internal/infrastructure/db/sqlc"
)

// insertOutboxEvents stores domain events in the outbox table. Call it with
// queries bound to the same transaction as the aggregate write so the state
// change and its events commit or roll back together.
func insertOutboxEvents(ctx context.Context, queries *db.Queries, domainEvents []events.DomainEvent) error {
	for _, event := range domainEvents {
		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		if err := queries.InsertOutboxEvent(ctx, db.InsertOutboxEventParams{
			ID:          event.EventId(),
			AggregateID: event.AggregateId(),
			EventName:   event.EventName(),
			Payload:     payload,
			OccurredAt:  timestamptzFromTime(event.OccurredAt()),
		}); err != nil {
			return err
		}
	}

	return nil
}
