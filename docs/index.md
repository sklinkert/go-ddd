# Why Domain-Driven Design makes your life easier

I've been writing Go backends for over a decade, and the same three problems show up in almost every codebase I've inherited:

1. **Business logic lives in HTTP handlers.** The rule "a product needs a price above zero" is enforced in a controller, so the importer, the admin CLI, and the message consumer each get their own slightly different copy of it.
2. **Validation is scattered and repeated.** Every function that touches a struct re-checks the same fields, because nobody can be sure whether the value in hand has been validated yet. Half the checks are missing, the other half disagree with each other.
3. **The model is anemic.** Structs are bags of public fields. Any code anywhere can set `product.Price = -5`, and the compiler is fine with it. The "domain model" is really just the database schema with JSON tags.

None of these feel like emergencies on day one. They compound. A year in, nobody can answer "where is the rule that a seller must have a name?" without grepping, and every change means re-discovering which of the five validation paths actually runs.

Domain-Driven Design is the set of patterns I keep coming back to because it attacks all three problems at the root: **put the business rules in one place, make them impossible to bypass, and name things the way the business names them.**

## What DDD actually is (without the jargon wall)

DDD has a reputation for ceremony: fat blue books, event storming workshops, a vocabulary quiz before you're allowed to write code. Ignore that for now. The tactical core is small:

- **Entities** are things with identity and a life cycle (a `Product`, a `Seller`). They protect their own invariants: you can't construct an invalid one.
- **Value objects** are things defined only by their value (`Money`, an email address). Immutable, validated at construction, safe to pass around.
- **Aggregates** are clusters of entities that change together. Each aggregate is a consistency boundary, and other aggregates refer to it by Id, never by embedding it.
- **Repositories** are interfaces the domain defines for loading and storing aggregates. The database is a detail behind them.
- **Domain events** are facts the domain announces ("product created") so other parts of the system can react without being called directly.

That's it. Everything else in this tutorial is those five ideas applied consistently.

## How this tutorial works

This isn't a theory course. Every chapter is anchored to a working codebase: [sklinkert/go-ddd](https://github.com/sklinkert/go-ddd), a production-grade template that models a small marketplace where sellers list products. It ships with a REST API (Echo), PostgreSQL via pgx and sqlc, migrations, testcontainers-based integration tests, race-safe idempotent commands, and a transactional outbox for domain events.

Each chapter shows the real code, explains why it's shaped that way, and tells you which trade-offs I made and where I'd decide differently in your situation. You can clone the repo and run everything locally with one command:

```bash
git clone https://github.com/sklinkert/go-ddd.git
cd go-ddd
docker compose up --build
```

The tutorial builds up in the order I'd introduce these patterns to a colleague:

1. [The domain and its language](tutorial/01-the-domain.md) — what we're modelling and why the words matter
2. [Entities that guard themselves](tutorial/02-entities.md) — constructors, invariants, and the validated-entity pattern
3. [Value objects, starting with Money](tutorial/03-value-objects.md) — why `float64` money is a bug and what to do instead
4. [Aggregates and their boundaries](tutorial/04-aggregates.md) — why `Product` stores a `SellerId`, not a `Seller`
5. [Repositories](tutorial/05-repositories.md) — interfaces in the domain, sqlc in the infrastructure
6. [CQRS, commands and queries](tutorial/06-cqrs.md) — separating writes from reads without going full event sourcing
7. [Domain events and the outbox](tutorial/07-domain-events-outbox.md) — telling the rest of the world, reliably
8. [Idempotent commands](tutorial/08-idempotency.md) — surviving client retries without duplicate writes
9. [Testing a DDD codebase](tutorial/09-testing.md) — fast domain tests, honest integration tests

If you're the reference-manual type, the [Architecture](reference/architecture.md) page has the layer rules and the onion diagram, and the [FAQ](reference/faq.md) answers the questions I get most often, starting with "isn't this overkill?".

## Is it worth it?

Fair question, and the honest answer is: not always. If your service is a thin CRUD layer over a database and the business rules fit in a sentence, DDD is overhead. The [FAQ](reference/faq.md) goes into when I'd skip it.

But the moment real rules show up — money, inventory, permissions, anything where "it depends" is the answer to how the system should behave — the patterns here pay for themselves. Not because they're elegant, but because they give every rule exactly one home, and they make the compiler enforce what code reviews otherwise have to catch.

Start with [chapter 1](tutorial/01-the-domain.md).
