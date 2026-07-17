# 5. Repositories

A **repository** gives the domain the illusion of an in-memory collection of aggregates: get one by Id, get them all, add, update, remove. Behind the illusion sits a database — but the domain never sees it, because of one structural decision that this chapter is really about:

**The interface lives in the domain. The implementation lives in infrastructure.**

## The dependency arrow points inward

Here is the entire [`ProductRepository`](https://github.com/sklinkert/go-ddd/blob/main/internal/domain/repositories/product_repository.go) as the domain knows it:

```go
package repositories

type ProductRepository interface {
    Create(ctx context.Context, product *entities.ValidatedProduct) (*entities.Product, error)
    FindById(ctx context.Context, id uuid.UUID) (*entities.Product, error)
    FindAll(ctx context.Context) ([]*entities.Product, error)
    Update(ctx context.Context, product *entities.ValidatedProduct) (*entities.Product, error)
    Delete(ctx context.Context, id uuid.UUID) error
}
```

Notice what's absent: no `*sql.DB`, no pgx types, no query strings, no ORM session. And notice what's present: `ValidatedProduct` in the write signatures — the type-level guarantee from [chapter 2](02-entities.md) that nothing unvalidated reaches storage.

Naively, the domain would depend on the database package. The interface flips the arrow: **infrastructure depends on the domain**, by implementing an interface the domain owns. Postgres becomes a plug-in. This is dependency inversion doing real work, not ceremony:

- Application services are testable with an in-memory fake — no Docker, no mocks-of-mocks. Milliseconds per test.
- The storage engine is swappable. Not hypothetically: this template migrated from GORM to sqlc/pgx, and the domain and application layers didn't change.
- Every capability the domain grants to persistence is enumerated in one small interface. There's no way to "just run a quick query" from a service — if a use case needs a new access pattern, the interface grows, visibly, in review.

## One repository per aggregate root

[Chapter 4](04-aggregates.md) drew the boundaries; repositories respect them. There's a `ProductRepository` and a `SellerRepository`, and there will never be a `ProductEventRepository` or a `MoneyRepository` — you load and store *aggregates*, whole, by their root. This is the discipline that kills the N+1/lazy-loading class of problems: nothing is ever half-loaded, and nothing inside a boundary is saved independently.

## The implementation: sqlc, not an ORM

The concrete side, [`SqlcProductRepository`](https://github.com/sklinkert/go-ddd/blob/main/internal/infrastructure/db/postgres/sqlc_product_repository.go), is built on [sqlc](https://sqlc.dev): you write real SQL in `.sql` files, sqlc generates type-safe Go from it at build time. The queries are plain enough to review:

```sql
-- name: GetProductById :one
SELECT p.id, p.name, p.price_cents, p.currency, p.seller_id, p.created_at, p.updated_at
FROM products p
JOIN sellers s ON p.seller_id = s.id
WHERE p.id = $1 AND p.deleted_at IS NULL AND s.deleted_at IS NULL;
```

I prefer this over an ORM in a DDD codebase for a specific reason: **the mapping between rows and domain objects is where invariants leak**, and I want that mapping explicit and dumb. No reflection deciding what loads when, no save() cascading through an object graph the aggregate design says shouldn't exist. The SQL does exactly what it says, at compile time.

(The soft-delete filter `deleted_at IS NULL` also shows why queries belong to the implementation: "deleted products don't exist" is enforced in every read path, in one reviewed place.)

## Reconstruction goes through the constructors

Reading a row back is a boundary crossing, and the same rule from [chapter 3](03-value-objects.md) applies — reconstruction routes through the validating constructors:

```go
func productFromRow(id uuid.UUID, name string, priceCents int64, currency string, /* ... */) (*entities.Product, error) {
    price, err := entities.NewMoney(priceCents, entities.Currency(currency))
    if err != nil {
        return nil, err
    }

    return &entities.Product{
        Id:    id,
        Name:  name,
        Price: price,
        // ...
    }, nil
}
```

If someone hand-edits a row to `currency = 'XXX'`, the repository refuses to materialize it rather than letting corrupt money flow into the domain. The database is *inside* the trust boundary for most teams; treating it as slightly outside costs two comparisons and has saved me real incidents.

## Errors: sentinels, not strings, not nils-with-vibes

Two deliberate choices in the read/write paths:

```go
if errors.Is(err, pgx.ErrNoRows) {
    return nil, nil   // FindById: absence is not an error
}
```

```go
if rows == 0 {
    return nil, entities.ErrProductNotFound   // Update: absence IS an error
}
```

They look inconsistent; they're not. *Finding* nothing is a normal outcome the caller asked about — `(nil, nil)` lets the application layer decide it's a 404. *Updating* nothing means the caller acted on an aggregate that doesn't exist — that's `ErrProductNotFound`, a **domain** sentinel, so the REST layer can map it with `errors.Is` and no layer above infrastructure ever sees a pgx error type. Infrastructure errors don't leak; domain errors do the traveling.

## Where transactions live

The repository owns the transaction, not the service. `Create` opens one, writes the product *and* its domain events (the outbox — [chapter 7](07-domain-events-outbox.md)), reads the row back, and commits:

```go
tx, err := repo.pool.Begin(ctx)
defer func() { _ = tx.Rollback(ctx) }()

qtx := repo.queries.WithTx(tx)
// insert product + insert outbox events + read back
return created, tx.Commit(ctx)
```

Why not let the application service compose transactions? Because [chapter 4](04-aggregates.md) already decided the unit of consistency: one aggregate, one transaction. The repository is exactly that unit, so the transaction hides behind the interface. The day a use case genuinely needs to commit two aggregates atomically, that's a design smell to revisit at the aggregate level first — not a reason to hand every service a transaction handle.

## Try it

1. Write an in-memory `ProductRepository` backed by a `map[uuid.UUID]*entities.Product` and a mutex. Note the service tests in `internal/application/services/` already do this — read their fakes and see how little code the interface demands.
2. Add `FindBySellerId(ctx, sellerId)`. Feel the shape of the change: one interface method, one SQL query, one generated function, one implementation method. Nothing else moves.
3. Delete the `deleted_at IS NULL` clause from one query and run the integration tests (`make test`). The soft-delete tests catch it — the query file is covered by tests against a real Postgres in a container, which is [chapter 9's](09-testing.md) topic.

Next: [CQRS, commands and queries](06-cqrs.md) — why the write path and the read path stopped sharing code.
