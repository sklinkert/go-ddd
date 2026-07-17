# 9. Testing a DDD codebase

I don't trust my code — I test it. But *what* to test at *which* level is where teams burn time, and it's also where this architecture quietly pays out its biggest dividend: **the layering decides the testing strategy for you.**

The rule of thumb the template follows:

| Layer | What's tested | Against | Speed |
|---|---|---|---|
| Domain | Business rules, invariants, events | Nothing — pure Go | microseconds |
| Application | Orchestration, idempotency, error flows | In-memory fakes | milliseconds |
| Infrastructure | SQL, mapping, transactions | **Real Postgres** (testcontainers) | seconds |

Each layer's tests answer a different question, and none of them re-answers a lower layer's question.

## Domain tests: pure functions in disguise

The domain has no dependencies — no database, no framework, no clock injection ceremony. So its tests are plain table stakes Go, and they read like the business rules they verify ([`internal/domain/entities/`](https://github.com/sklinkert/go-ddd/tree/main/internal/domain/entities)):

```go
func TestNewMoney_RejectsNegativeAmount(t *testing.T) {
    _, err := entities.NewMoney(-1, entities.EUR)

    assert.ErrorIs(t, err, entities.ErrValidation)
}

func TestProduct_PullEventsClearsEvents(t *testing.T) {
    product := entities.NewProduct("Widget", price, validatedSeller)

    first := product.PullEvents()
    second := product.PullEvents()

    assert.Len(t, first, 1)
    assert.Empty(t, second)
}
```

Two things to notice. Assertions target **sentinel errors** (`assert.ErrorIs(t, err, ErrValidation)`), not message strings — messages are documentation and may be reworded; sentinels are contract. And the *whole invariant surface* of an entity is testable without starting anything: this is the payoff of [chapter 2's](02-entities.md) insistence that rules live in the entity. When a colleague asks "what are the rules for a product?", the entity test file is a readable, executable answer.

These tests run in microseconds, so they run constantly — on save, in pre-commit, wherever. There's no excuse for an untested business rule when testing it costs this little, which is precisely the incentive structure you want.

## Application tests: fakes, not mock frameworks

Application services orchestrate; their tests verify the orchestration — commands flow, errors propagate, idempotency holds. Because repositories are interfaces owned by the domain ([chapter 5](05-repositories.md)), the tests hand services simple in-memory fakes:

```go
type MockProductRepository struct {
    products []*entities.ValidatedProduct
}

func (m *MockProductRepository) Create(ctx context.Context, product *entities.ValidatedProduct) (*entities.Product, error) {
    m.products = append(m.products, product)
    return &product.Product, nil
}
```

A slice and some loops — no mock framework, no `EXPECT().Times(1)` choreography. I prefer hand-written fakes in this layer for a concrete reason: expectation-style mocks test *how* the service talks to its dependencies, which welds tests to the implementation. Fakes test *what ends up true*, which survives refactors. When I moved this template's persistence from GORM to sqlc, the application tests didn't change — that's the property to protect.

The star of this layer is the idempotency suite ([chapter 8](08-idempotency.md)): eight goroutines hammering one key, asserting `executions == 1`. Concurrency tests belong here, against fast fakes, where you can afford to run them thousands of times with `-race` — not against a container where each attempt costs seconds.

## Infrastructure tests: a real database or it doesn't count

Here's the strong opinion: **mocking the database in repository tests is worthless.** The repository's entire job is SQL, transactions, and row mapping. A test that mocks the SQL away verifies that the code calls the mock — nothing else. Every real repository bug I've seen lived exactly in the parts a mock skips: a query that ignores `deleted_at`, a mapping that flips two columns, a transaction that doesn't actually cover the outbox insert.

So the template's repository tests run against real Postgres via [testcontainers](https://testcontainers.com) ([`internal/testhelpers/postgres_test_container.go`](https://github.com/sklinkert/go-ddd/blob/main/internal/testhelpers/postgres_test_container.go)):

```go
func SetupTestDB(t *testing.T) *PostgresTestContainer {
    postgresContainer, err := postgres.Run(ctx,
        "postgres:17-alpine",
        postgres.WithDatabase(dbName),
        postgres.WithUsername(dbUser),
        postgres.WithPassword(dbPassword),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2),
        ),
    )
    require.NoError(t, err, "Failed to start postgres container")
    // connect, apply schema, return pool + queries
}
```

Each test package gets a disposable Postgres 17 in Docker: real unique constraints (the idempotency reservation *depends* on one), real transaction semantics (the outbox atomicity claim is *tested*, not asserted), real `ON CONFLICT` behavior. When the test passes, it means what it says.

The cost is seconds of startup per package, and it's paid honestly: these tests run in CI on every push (the GitHub Actions runner ships with Docker), and locally via `make test`. What you get for those seconds is the class of confidence mocks can't sell you — for example, [`outbox_test.go`](https://github.com/sklinkert/go-ddd/blob/main/internal/infrastructure/db/postgres/outbox_test.go) proves that a failed product insert rolls back the outbox rows too. Try proving that with a mock.

## What deliberately isn't tested

Symmetry demands the other list:

- **No tests that a fake returns what the fake was told to return.** If a test can only fail when the test's own setup changes, delete it.
- **No E2E test pyramid-tip for every feature.** The template wires everything together in `main.go`; one smoke path (bring the stack up, create a seller and product, see the event relay fire) covers the wiring. Feature behavior is already covered below, faster.
- **No coverage worship.** Generated code (sqlc's output) is excluded from the metric. Chasing a number through generated files is how teams end up with impressive dashboards and untested invariants.

## The feedback loop is the feature

The reason to care about this pyramid isn't ideology, it's *iteration speed*. Rule change? Domain test, microseconds. New orchestration? Fake-backed test, milliseconds. New query? One container-backed test, seconds. Each change lands in the cheapest layer that can catch its bugs — and the architecture is what made those layers separable in the first place.

That's the quiet argument for DDD I'd make to a skeptical team: forget the vocabulary — a codebase where business rules are testable in microseconds *without Docker* is a codebase where the rules actually get tested.

## Try it

1. `make test` — watch the domain and application packages finish before the first container is even up.
2. Break a query on purpose (drop the `deleted_at IS NULL` from `GetProductById`) and see which layer catches it. Then break a business rule (allow zero prices) and see which layer catches *that*. Different layers, by design.
3. Write the test for exercise 3 of [chapter 7](07-domain-events-outbox.md): does `Update` insert outbox events for a `ProductPriceChanged`? You now have all three tools — decide which layer this test belongs in, and why. (My answer: infrastructure — it's a claim about a transaction.)

---

That's the tutorial. For the layer rules in one page, see [Architecture](../reference/architecture.md); for the questions everyone asks next, the [FAQ](../reference/faq.md). And if this codebase is the way you like to learn, [star the repo](https://github.com/sklinkert/go-ddd) — it genuinely helps more people find it.
