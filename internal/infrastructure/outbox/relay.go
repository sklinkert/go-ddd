package outbox

import (
	"context"
	"log/slog"
	"time"

	db "github.com/sklinkert/go-ddd/internal/infrastructure/db/sqlc"
)

// Publisher forwards an event payload to the outside world (message broker,
// webhook, ...). Implementations must be safe to call repeatedly with the
// same event: the outbox guarantees at-least-once, not exactly-once.
type Publisher interface {
	Publish(ctx context.Context, eventName string, payload []byte) error
}

// SlogPublisher logs events instead of sending them anywhere. Swap it for a
// Kafka/NATS/SQS publisher in a real deployment.
type SlogPublisher struct{}

func (SlogPublisher) Publish(ctx context.Context, eventName string, payload []byte) error {
	slog.InfoContext(ctx, "publishing domain event",
		slog.String("event", eventName), slog.String("payload", string(payload)))
	return nil
}

// Relay polls the outbox table and publishes unpublished events. Run exactly
// one instance (or guard with FOR UPDATE SKIP LOCKED when scaling out).
type Relay struct {
	queries   *db.Queries
	publisher Publisher
	interval  time.Duration
	batchSize int32
}

func NewRelay(queries *db.Queries, publisher Publisher, interval time.Duration) *Relay {
	return &Relay{
		queries:   queries,
		publisher: publisher,
		interval:  interval,
		batchSize: 100,
	}
}

// Start blocks until ctx is cancelled.
func (r *Relay) Start(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.relayBatch(ctx); err != nil {
				slog.ErrorContext(ctx, "outbox relay batch failed", slog.Any("error", err))
			}
		}
	}
}

func (r *Relay) relayBatch(ctx context.Context) error {
	events, err := r.queries.GetUnpublishedOutboxEvents(ctx, r.batchSize)
	if err != nil {
		return err
	}

	for _, event := range events {
		if err := r.publisher.Publish(ctx, event.EventName, event.Payload); err != nil {
			// Stop the batch; unpublished events are retried next tick.
			return err
		}

		if err := r.queries.MarkOutboxEventPublished(ctx, event.ID); err != nil {
			return err
		}
	}

	return nil
}
