# 8. Idempotent commands

Your client sends `POST /products`. The server creates the product, then the load balancer kills the connection before the response gets out. The client sees a timeout. What does every reasonable HTTP client do on a timeout? It retries. Now you have two products.

This is not an exotic failure mode — it's Tuesday. Mobile clients on flaky networks, rolling deploys, proxy timeouts: anything separating "the server did the work" from "the client learned about it" will eventually duplicate your writes. The standard fix is an **idempotency key**: the client generates a unique key per logical operation, sends it with the request, and the server guarantees the same key never executes the operation twice.

The concept is simple. The implementation is where I got burned, so this chapter walks the version the template actually ships — including the parts I got wrong the first time.

## The naive version (which I shipped)

```go
existing, _ := s.idempotencyRepo.FindByKey(ctx, cmd.IdempotencyKey)
if existing != nil {
    return unmarshalCachedResponse(existing)
}

result, err := s.doCreateProduct(ctx, cmd)
if err != nil {
    return nil, err
}

s.idempotencyRepo.Store(ctx, cmd.IdempotencyKey, result)
return result, nil
```

Check, execute, store. It reads correctly and passes every test you'll write for it, because sequential tests can't see the bug.

The bug: two requests with the same key arrive *at the same time*. Both call `FindByKey`, both get `nil` (storing happens after execution), both execute. Two products, one idempotency key — the mechanism silently did nothing exactly when it mattered. And concurrent duplicates are the **main** case, not a corner case: retries fire because something is slow, and when something is slow, the original request is usually still running.

## Fix one: reserve the key atomically

The check and the claim must be one atomic operation. In Postgres that's a unique constraint plus `ON CONFLICT DO NOTHING` ([`sql/queries/idempotency.sql`](https://github.com/sklinkert/go-ddd/blob/main/sql/queries/idempotency.sql)):

```sql
-- name: ReserveIdempotencyKey :execrows
-- Atomically claims the key. Zero rows means another request already holds it.
INSERT INTO idempotency_records (id, key, request, response, status_code, created_at)
VALUES ($1, $2, $3, '', 0, $4)
ON CONFLICT (key) DO NOTHING;
```

`:execrows` makes sqlc return the affected-row count, and that count is the whole trick: **1 means you won the race and may execute; 0 means someone else holds the key and you must not.** The row is inserted *before* execution as a reservation (`response = ''`, `status_code = 0` meaning "in flight"), not after execution as a cache entry. Postgres's unique index is the arbiter of who executes — no advisory locks, no Redis, no distributed-lock library. The guarantee comes from the database, which is exactly where I want it.

## The losers need answers too

If you lost the race, what do you tell the client? Depends on the winner ([`internal/application/services/idempotency.go`](https://github.com/sklinkert/go-ddd/blob/main/internal/application/services/idempotency.go)):

```go
reserved := false
for attempt := 0; attempt < 3 && !reserved; attempt++ {
    reserved, err = repo.Reserve(ctx, record)
    if err != nil {
        return nil, err
    }
    if reserved {
        break
    }

    existing, err := repo.FindByKey(ctx, key)
    if err != nil {
        return nil, err
    }
    if existing == nil {
        continue // released between Reserve and FindByKey; try again
    }

    if existing.Request != string(requestJSON) {
        return nil, ErrIdempotencyKeyReuse
    }

    if existing.IsCompleted() {
        var result T
        if err := json.Unmarshal([]byte(existing.Response), &result); err != nil {
            return nil, fmt.Errorf("unmarshal cached idempotency response: %w", err)
        }
        return &result, nil
    }

    if time.Since(existing.CreatedAt) < reservationTTL {
        return nil, ErrRequestInFlight
    }

    // Stale reservation: the previous holder crashed before completing.
    if err := repo.Delete(ctx, key); err != nil {
        return nil, err
    }
}
```

Four outcomes, each deliberate:

- **Winner finished** → serve its stored response. The retry gets the same `201` body the original would have gotten, byte for byte.
- **Winner still running** → `ErrRequestInFlight`, mapped to **409 Conflict**. Back off, retry, hit the cached response.
- **Same key, different payload** → `ErrIdempotencyKeyReuse`, mapped to **422**. Easy to skip, dangerous to skip: if a client bug reuses a key across different requests, silently returning request A's cached response to request B means the caller thinks it created B while holding A. Fail loudly.
- **Reservation older than the TTL with no response** → the holder is dead; take over (next section).

Generics make this one decorator for every command — `withIdempotency[T any]` wraps `CreateProduct`, `UpdateSeller`, and friends identically, which is why [chapter 6's](06-cqrs.md) service body starts with it.

## Fix two: crashes must not brick a key

Reservation-before-execution creates a new failure mode: reserve, then die — OOM kill, deploy — before storing a response. The row says "in flight" forever; every retry gets 409 until a human deletes it. You've traded duplicate writes for a permanently wedged operation.

That's the `reservationTTL` branch above. A reservation past the TTL with no response means the holder is dead; the next retry deletes the stale row and re-reserves. Because re-reserving is the same atomic `INSERT ... ON CONFLICT`, two retries racing to take over still resolve to one winner.

Sizing the TTL is a judgment call: comfortably longer than your slowest *legitimate* execution, or you'll take over reservations that are merely slow — and be back to duplicates. A minute is generous for CRUD; a batch job needs a different mechanism.

## Fix three: release on failure — with a detached context

If execution fails, release the reservation so the client's retry can actually run. Easy. Except: *why* did execution fail? Often because the client disconnected and the request context got cancelled. Call `repo.Delete(ctx, key)` with that cancelled context and the DELETE never reaches Postgres — the reservation leaks, and the disconnecting client (the one most likely to retry!) is locked out for a full TTL. The cleanup fails precisely in the scenario it exists for.

`context.WithoutCancel` (Go 1.21+) keeps the context's values — trace Ids, loggers — but detaches cancellation:

```go
result, err := execute()
if err != nil {
    if deleteErr := repo.Delete(context.WithoutCancel(ctx), key); deleteErr != nil {
        slog.WarnContext(ctx, "failed to release idempotency key",
            slog.String("idempotency_key", key), slog.Any("error", deleteErr))
    }
    return nil, err
}

storeResponse(context.WithoutCancel(ctx), repo, key, result)
```

Same reasoning for `storeResponse` — the business operation already committed; persisting the cached response must not be at the mercy of the request context. And it's best-effort: failing the whole request because the *cache write* failed would lie to the client about work that's done. Log it; the TTL takeover covers the wedged row.

## Proving it under concurrency

Sequential tests can't catch the original bug, so [the test suite](https://github.com/sklinkert/go-ddd/blob/main/internal/application/services/idempotency_test.go) fires concurrent goroutines at one key and asserts the only thing that matters:

```go
assert.Equal(t, 1, executions, "exactly one caller may execute")
assert.Equal(t, callers, successes+inFlight)
```

Successes can exceed one — a loser arriving after the winner finished legitimately gets the cached response — but the business logic runs exactly once, always. Run it with `-race`.

Against the live stack, the behavior in four commands:

```bash
for i in $(seq 5); do
  curl -s -o /dev/null -w "%{http_code}\n" -X POST localhost:8080/api/v1/products \
    -H 'Content-Type: application/json' \
    -d '{"idempotency_key":"k-42","name":"Widget","price_cents":999,"currency":"EUR","seller_id":"..."}' &
done; wait
```

One `201`, four `409`, one row in the database. Retry a moment later: the winner's `201` body back, byte for byte. Same key with `"name":"Gadget"`: `422`.

## The one-sentence version

Idempotency is a **write-side claim, not a read-side check**. If your implementation reads before it writes, it has the race, full stop — make the database's unique constraint do the deciding and branch on rows-affected. Everything else here (TTL takeover, payload comparison, detached-context cleanup) is consequences of taking that sentence seriously.

## Try it

1. Reproduce the naive bug: comment out the `Reserve` call, make the wrapper check-then-store, and run the concurrency test. Watch `executions` climb past 1.
2. Set `reservationTTL` to a millisecond and rerun the tests. Which ones fail, and what duplicate behavior do they demonstrate? (This is why the TTL must exceed legitimate execution time.)
3. Trace the 409 and 422 from sentinel to status code in [`errors.go`](https://github.com/sklinkert/go-ddd/blob/main/internal/interface/api/rest/errors.go) — the same sentinel-error pattern from [chapter 5](05-repositories.md), now covering concurrency semantics.

Next: [testing a DDD codebase](09-testing.md) — why these layers make tests fast where they can be and honest where they must be.
