# 7. Domain events and the outbox

Here's a piece of code I've written, in some form, at three different companies:

```go
if err := s.repo.Create(ctx, product); err != nil {
    return err
}

// Tell the rest of the world
return s.kafka.Publish(ctx, "product.created", toEvent(product))
```

Save to Postgres, then publish to Kafka. It works in every demo, every test, and roughly 99.9% of the time in production. The remaining 0.1% is where the fun is: the database and the broker are two separate systems, and **no transaction spans both**.

- Crash after the DB commit, before the publish → the product exists, but the search index, pricing service, and notifications never hear about it. Silent data drift.
- Flip the order — publish first — and a failed insert produces a *ghost event*: consumers react to a product that was never created.

You can't "just be careful" your way out. A retry loop shrinks the window; it doesn't close it. This is the **dual-write problem**, and the fix is old and boring, which is exactly what you want: the **transactional outbox**. Don't publish the event — *store* it, in the same database, in the same transaction as the state change. A separate process does the actual publishing.

## First decision: events belong to the aggregate

Before any infrastructure: who creates the event? In many codebases the service layer builds it, right next to the publish call. That's backwards. "A product was created" is a domain fact, and the aggregate is the thing that knows its own state changed. So the aggregate records events as part of the change:

```go
func NewProduct(name string, price Money, seller ValidatedSeller) *Product {
    product := &Product{ /* ... */ }

    product.recordEvent(events.NewProductCreated(
        product.Id, name, price.MinorUnits(), string(price.Currency()), seller.Id))

    return product
}

// PullEvents returns the recorded domain events and clears them.
func (p *Product) PullEvents() []events.DomainEvent {
    pulled := p.domainEvents
    p.domainEvents = nil
    return pulled
}
```

The events are dumb structs — past-tense names, immutable, no behavior ([`internal/domain/events/`](https://github.com/sklinkert/go-ddd/blob/main/internal/domain/events/product_events.go)):

```go
type DomainEvent interface {
    EventId() uuid.UUID
    EventName() string
    OccurredAt() time.Time
    AggregateId() uuid.UUID
}

func (e ProductCreated) EventName() string { return "product.created" }
```

Two details worth noticing. The event Id is a UUIDv7 — time-ordered, so it sorts nicely and doubles as a deduplication key for consumers. And `PullEvents` *clears* the slice, so the repository pulls exactly once per save and a retried save can't double-insert the same events.

## One transaction or it didn't happen

The outbox table ([`migrations/000003_outbox.up.sql`](https://github.com/sklinkert/go-ddd/blob/main/migrations/000003_outbox.up.sql)) is deliberately simple:

```sql
CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    aggregate_id UUID NOT NULL,
    event_name TEXT NOT NULL,
    payload JSONB NOT NULL,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL,
    published_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_outbox_events_unpublished
    ON outbox_events(occurred_at) WHERE published_at IS NULL;
```

Note the **partial index**. The relay only ever asks one question — "unpublished events, oldest first" — while the table grows forever. An index on `WHERE published_at IS NULL` stays tiny no matter how many millions of published rows accumulate, because rows drop out of it the moment they're marked published.

The payoff happens in the repository ([chapter 5](05-repositories.md) showed the transaction): aggregate insert and outbox insert share one `pgx` transaction.

```go
qtx := repo.queries.WithTx(tx)

if _, err := qtx.CreateProduct(ctx, /* ... */); err != nil {
    return nil, err
}

if err := insertOutboxEvents(ctx, qtx, product.PullEvents()); err != nil {
    return nil, err
}
// ... read-back, then tx.Commit(ctx)
```

Either the product row *and* its events commit, or neither does. No orphaned product, no ghost event. And the service layer knows nothing about any of this — it calls `repo.Create` and the events ride along. **You can't forget to publish, because there is no publish step to forget.**

## The relay: dumb on purpose

Something still has to move events from Postgres to the broker. That's the [relay](https://github.com/sklinkert/go-ddd/blob/main/internal/infrastructure/outbox/relay.go) — a loop that polls unpublished rows and hands them to a `Publisher`:

```go
type Publisher interface {
    Publish(ctx context.Context, eventName string, payload []byte) error
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
```

In the template the publisher just logs via `slog`; in production you swap in Kafka, NATS, SQS. The interesting property is the failure mode: publish succeeds, then `MarkOutboxEventPublished` fails — crash, network blip, deploy. Next tick, the row is still unpublished, so it publishes **again**.

That's not a bug; it's the contract: the outbox gives you **at-least-once** delivery, never exactly-once. Every consumer must be idempotent — handling `product.created` twice must equal handling it once. The event Id is the dedup key. If that sounds like a burden: you needed idempotent consumers anyway. Kafka redelivers on consumer-group rebalances all by itself. At-least-once is the honest default of distributed messaging; the outbox just stops pretending otherwise.

## The caveats nobody puts in the diagram

**Ordering.** Poll order is `occurred_at`; with a single relay, one aggregate's events come out in recording order. That's per-relay ordering, not global. Publishing to a partitioned topic? Partition by `aggregate_id`. Cross-aggregate ordering is a promise you should never make.

**Scaling the relay.** Two naive relay instances grab the same batch and double-publish everything. The standard fix is `FOR UPDATE SKIP LOCKED` in the poll query — each relay locks the rows it's working; competitors skip them. The template ships the single-instance version deliberately: add SKIP LOCKED when you measure the need, not before.

**Polling vs. CDC.** Polling every few seconds is fine for a huge range of workloads and needs zero extra infrastructure. Debezium tailing the WAL is lower-latency and much more machinery. Start with polling.

**Cleanup.** Published rows pile up; a nightly `DELETE ... WHERE published_at < now() - interval '30 days'` keeps the table sane. The partial index doesn't care either way.

## When to skip all of this

Often, honestly. If the "event handler" lives in the same process — creating a product should warm a cache — call the function; in-process, in-transaction, done. And if losing the occasional event is tolerable (analytics pings), fire-and-forget with a retry is genuinely fine.

The outbox earns its keep exactly when a state change in *your* database must reliably reach *another* system. The moment someone says "when X happens here, Y must happen over there" — reach for it. It's maybe 200 lines including the migration, and it turns a distributed-systems problem into a table and a for loop.

!!! tip "Standalone version"
    I extracted a broker-agnostic implementation of this pattern into [go-outbox](https://github.com/sklinkert/go-outbox) if you want the outbox without the template.

## Try it

1. Run the stack (`make docker-up`), create a product, and watch the relay log `publishing domain event` with `event_name=product.created`. Then check the row: `SELECT event_name, published_at FROM outbox_events;`
2. Kill the app between insert and relay tick (set a long poll interval), restart, and confirm the event still goes out. That's the whole pattern in one experiment.
3. Add a `ProductPriceChanged` event: record it in `UpdatePrice`, and check `Update` in the repository — does it insert outbox events today? (Look. This is a real extension point.)

Next: [idempotent commands](08-idempotency.md) — the other half of surviving retries, this time on the way *in*.
