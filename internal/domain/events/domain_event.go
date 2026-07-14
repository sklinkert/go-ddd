package events

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent is implemented by every event raised in the domain layer.
// Events are facts: named in past tense and immutable once created.
type DomainEvent interface {
	EventId() uuid.UUID
	EventName() string
	OccurredAt() time.Time
	// AggregateId identifies the aggregate the event belongs to.
	AggregateId() uuid.UUID
}

// BaseEvent carries the fields shared by all domain events. Embed it and
// implement EventName on the concrete event.
type BaseEvent struct {
	Id          uuid.UUID
	Aggregate   uuid.UUID
	OccurredAtT time.Time
}

func NewBaseEvent(aggregateId uuid.UUID) BaseEvent {
	return BaseEvent{
		Id:          uuid.Must(uuid.NewV7()),
		Aggregate:   aggregateId,
		OccurredAtT: time.Now(),
	}
}

func (e BaseEvent) EventId() uuid.UUID     { return e.Id }
func (e BaseEvent) AggregateId() uuid.UUID { return e.Aggregate }
func (e BaseEvent) OccurredAt() time.Time  { return e.OccurredAtT }
