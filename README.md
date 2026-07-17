# Go-DDD: Build Domain-Driven Go Services Fast

[![CI](https://github.com/sklinkert/go-ddd/actions/workflows/go.yml/badge.svg)](https://github.com/sklinkert/go-ddd/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/sklinkert/go-ddd/branch/main/graph/badge.svg)](https://codecov.io/gh/sklinkert/go-ddd)
[![Go Reference](https://pkg.go.dev/badge/github.com/sklinkert/go-ddd.svg)](https://pkg.go.dev/github.com/sklinkert/go-ddd)
[![Go Version](https://img.shields.io/github/go-mod/go-version/sklinkert/go-ddd)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[English](README.md) | [简体中文](README.zh-CN.md)

`go-ddd` jump-starts production-grade Go backends that keep business rules, infrastructure, and delivery code cleanly separated. Out of the box you get opinionated DDD building blocks, CQRS command and query flows, idempotent write paths, domain events with a transactional outbox, and tooling to keep schema and code in lockstep.

> 📚 **New to DDD? [Start with the tutorial →](https://sklinkert.github.io/go-ddd/)**
> A nine-chapter walkthrough that teaches Domain-Driven Design from zero, using this codebase as the running example.

## Why This Template

- **Model-first defaults** – Onion architecture keeps the domain pure while application services orchestrate infrastructure concerns.
- **Battle-tested patterns** – Commands, queries, repositories, value objects, domain events, and soft deletes mirror patterns used in real-world enterprise applications.
- **Idempotent pipelines** – Race-safe idempotency keys (atomic reservation, no check-then-write) make every command retry-safe.
- **Domain events + outbox** – State changes and their events are persisted together; a relay publishes them with at-least-once delivery.
- **Migration discipline** – SQL migrations, `migrate.go`, and `sqlc` make schema evolution explicit and reproducible.

## What You Get

- Marketplace example that demonstrates aggregates (Seller, Product), a `Money` value object, cross-module interactions, and validation rules.
- Layered modules under `internal/` for `domain`, `application`, `infrastructure`, `interface`, plus `testhelpers` for fixture reuse.
- Executable entrypoint at `cmd/marketplace/main.go` ready to wire adapters or frameworks of your choice.
- Database assets in `migrations/` and `sql/` plus generated data access via `sqlc`.
- [OpenAPI spec](api/openapi.yaml), `/healthz` + `/readyz` probes, and a one-command Docker Compose stack.

## Tech Stack Essentials

- **Go 1.26** with idiomatic patterns and testify-powered tests.
- **Echo v4** HTTP stack for REST endpoints.
- **pgx/v5** and `sqlc` for type-safe PostgreSQL access.
- **golang-migrate** handling SQL schema migrations
- **Testcontainers** integration to provision disposable Postgres instances during tests.
- **google/uuid** helpers for deterministic ID generation inside the domain.
- **golangci-lint** for static analysis, plus a `Makefile`, `Dockerfile`, and `docker-compose.yml` for a one-command local stack.

## Design Principles in Action

Domain-Driven Design connects implementation to an evolving model. `go-ddd` showcases this by modelling a simple marketplace where `Sellers` manage `Products`, exercising aggregates, value objects, and validation flows.

Anatomy of a write request:

```mermaid
sequenceDiagram
    participant C as Client
    participant R as REST Controller
    participant S as Application Service
    participant D as Domain
    participant P as Postgres

    C->>R: POST /api/v1/products (idempotency_key)
    R->>S: CreateProductCommand
    S->>P: Reserve idempotency key (atomic)
    S->>D: NewProduct(name, Money, ValidatedSeller)
    D->>D: validate() → ValidatedProduct + ProductCreated event
    S->>P: INSERT product + outbox event
    P-->>S: persisted row
    S->>P: store response for idempotency key
    S-->>R: CommandResult
    R-->>C: 201 Created (JSON)
    Note over P: Outbox relay publishes<br/>ProductCreated asynchronously
```

## Documentation

📚 **[DDD from zero: the full tutorial](https://sklinkert.github.io/go-ddd/)** — nine chapters from "why DDD" to entities, value objects, aggregates, repositories, CQRS, the outbox, idempotency, and testing, all anchored to this codebase.

📖 **[Comprehensive DDD & CQRS Principles Guide](DDD_CQRS_PRINCIPLES.md)** - Learn how to apply these patterns to any business domain.

## Repository Structure

![ddd-diagram-onion.png](ddd-diagram-onion.png)


- `domain`: The heart of the software, representing business logic and rules.
    - `entities`: Fundamental objects within our system, like `Product` and `Seller`. Contains basic validation logic.
- `application`: Contains use-case specific operations that interact with the domain layer.
- `infrastructure`: Supports the higher layers with technical capabilities like database access.
    - `db`: Database access and models.
    - `repositories`: Concrete implementations of our storage needs.
- `interface`: The external layer which interacts with the outside world, like API endpoints.
    - `api/rest`: Handlers or controllers for managing HTTP requests and responses.

## Further principles

- Domain
  - Must not depend on other layers.
  - Provides infrastructure with interfaces, but must not access infrastructure.
  - Implements business logic and rules.
  - Executes validations on entities. Validated entities are passed to the infrastructure layer.
  - Domain layer sets defaults of entities (e.g. uuid for ID or creation timestamp). Don't set defaults in the infrastructure layer or even database!
  - Do not leak domain objects to the outside world.
- Application
  - The glue code between the domain and infrastructure layer.
- Infrastructure
   - Repositories are responsible for translating a domain entity to a database model and retrieving it. No business logic is executed here.
   - Implements interfaces defined by the domain layer.
   - Implements persistence logic like accessing a postgres or mysql database.
   - When writing to storage, read written data before returning it. This ensures that the data is written correctly.

## Best Practices

- Don't return validated entities from read methods in the repository. Instead, return the domain entity type directly.
  - Validations will change over time. You don't want to migrate all the data in your database. Instead, you should guarantee you can always load historical data, regardless of how your validation logic has evolved.
  - Otherwise, you won't be able to read data from the database that was written with a different validation logic. You will have to handle errors at runtime.
  - Push validation to the write side-creation (NewX) and update methods - where you must enforce invariants anyway.
- Don't put default values (e.g current timestamp or ID) in the database. Set them in the domain layer (factory!) for several reasons:
  - It's quite dangerous to have two sources of truth.
  - It's easier to test the domain layer.
  - Databases can get replaced, and you don't want to have to change all your default values. 
- Always read the entity after write in the infrastructure layer.
  - This ensures that the data is written correctly, and we are never operating on stale data.
- `find` vs `get`:
  - `find` methods can return null or an empty list.
  - `get` methods must return a value. If the value is not found, throw an error.
- Deletion: Always use soft deletion. Create a `deleted_at` column in your database and set it to the current timestamp when deleting an entity. This way, you can always restore the entity if needed.

## CQRS and Idempotency

### Command Query Responsibility Segregation (CQRS)
CQRS separates read operations (queries) from write operations (commands) in your application. In this codebase:
- **Commands** modify state (CreateSellerCommand, CreateProductCommand, UpdateSellerCommand)
- **Queries** retrieve data without side effects (FindAllSellers, FindSellerById)

This separation enables different optimization strategies:
- **Write optimization**: Commands can use normalized schemas, ACID transactions, and strong consistency
- **Read optimization**: Queries can use denormalized views, caching, read replicas, or even different databases (e.g., PostgreSQL for writes, Elasticsearch for reads)
- **Scalability**: Read and write workloads can be scaled independently based on actual usage patterns
- **Performance**: Complex queries don't impact write performance, and write locks don't block read operations

### Idempotency Keys
Idempotency ensures that multiple identical requests have the same effect as a single request. This is crucial for handling network failures and retries in distributed systems. Implementation:
- Each command accepts an optional `idempotency_key` in the request
- The key is **reserved atomically** (`INSERT ... ON CONFLICT DO NOTHING`), so two concurrent requests with the same key can never both execute — no check-then-write race
- A completed request returns its cached response; a still-running one returns an "in progress" error so the client retries later
- Reusing a key with a **different payload** is rejected instead of silently returning the wrong cached response
- If the command fails, the reservation is released so the client can retry; reservations orphaned by a crash expire after a TTL

This prevents duplicate entities from being created when clients retry failed requests.

### Domain Events and the Transactional Outbox

Aggregates record events (e.g. `ProductCreated`) when something business-relevant happens. Instead of publishing them directly to a broker — which risks losing events when the process crashes between the DB commit and the publish — events are stored in an `outbox_events` table. A relay polls the outbox and publishes unpublished events with at-least-once delivery. See `internal/domain/events/` and `internal/infrastructure/outbox/`.

## Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema management. Migrations are stored in the `migrations/` directory with sequential version numbers.

### Migration Files Structure
```
migrations/
├── 000001_initial_schema.up.sql    # Creates initial tables
├── 000001_initial_schema.down.sql  # Rollback for initial schema
├── 000002_price_as_money.up.sql    # Money as integer cents + currency
├── 000003_outbox.up.sql            # Transactional outbox table
└── ...
```

### Running Migrations

#### Using the built-in utility:
```bash
# Apply all pending migrations
go run migrate.go -database-url "postgres://user:pass@localhost/db?sslmode=disable" -command up

# Rollback last migration
go run migrate.go -database-url "postgres://user:pass@localhost/db?sslmode=disable" -command down -steps 1

# Check current version
go run migrate.go -database-url "postgres://user:pass@localhost/db?sslmode=disable" -command version

# Force to specific version (use with caution)
go run migrate.go -database-url "postgres://user:pass@localhost/db?sslmode=disable" -command force -version 1
```

#### Using the CLI tool directly:
```bash
# Install the CLI tool
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Apply all pending migrations
migrate -path migrations -database "postgres://user:pass@localhost/db?sslmode=disable" up

# Rollback last migration
migrate -path migrations -database "postgres://user:pass@localhost/db?sslmode=disable" down 1
```

### Creating New Migrations
```bash
# Create a new migration
migrate create -ext sql -dir migrations -seq add_user_email_column
```

This will create two files:
- `0000NN_add_user_email_column.up.sql` - Forward migration
- `0000NN_add_user_email_column.down.sql` - Rollback migration

### Migration Best Practices
- Always create both `up` and `down` migrations
- Test migrations on a copy of production data
- Keep migrations small and focused
- Never modify existing migration files once they've been applied in production
- Use descriptive names for migration files

## Getting Started

> Requires **Go 1.26+**. Run `make help` to see all available targets.

### Quickstart with Docker Compose

Bring up Postgres, apply migrations, and start the API in one command:

```bash
make docker-up        # docker compose up --build
# API is now available on http://localhost:8080
make docker-down      # tear everything down
```

### Try the API in 30 seconds

Create a seller:

```bash
curl -s -X POST http://localhost:8080/api/v1/sellers \
  -H 'Content-Type: application/json' \
  -d '{"name": "Acme Corp", "idempotency_key": "create-acme-1"}'
```

```json
{"id":"0197a3c2-...","name":"Acme Corp","created_at":"2026-07-14T09:00:00Z","updated_at":"2026-07-14T09:00:00Z"}
```

Create a product for that seller (prices are integer cents — never floats):

```bash
curl -s -X POST http://localhost:8080/api/v1/products \
  -H 'Content-Type: application/json' \
  -d '{"name": "Wooden Chair", "price_cents": 4999, "currency": "EUR", "seller_id": "<seller-id-from-above>"}'
```

```json
{"id":"0197a3c3-...","name":"Wooden Chair","price_cents":4999,"currency":"EUR","seller_id":"0197a3c2-...","created_at":"...","updated_at":"..."}
```

Replay a request with the same `idempotency_key` — you get the cached response back instead of a duplicate seller:

```bash
curl -s -X POST http://localhost:8080/api/v1/sellers \
  -H 'Content-Type: application/json' \
  -d '{"name": "Acme Corp", "idempotency_key": "create-acme-1"}'
# → identical response, no second row created
```

List products and check service health:

```bash
curl -s http://localhost:8080/api/v1/products
curl -s http://localhost:8080/readyz
```

The full API is described in the [OpenAPI spec](api/openapi.yaml).


### Local development

1. Clone this repository:
```bash
git clone https://github.com/sklinkert/go-ddd.git
cd go-ddd
go mod download
```

2. Install sqlc (for development):
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

3. Generate database code (if you modify SQL queries):
```bash
sqlc generate
```

4. Set up your PostgreSQL database and run migrations:
```bash
# Set your database connection URL
export DATABASE_URL="postgres://user:password@localhost/dbname?sslmode=disable"

# Run migrations using the built-in utility
go run migrate.go -command up

# Or use the CLI tool directly
migrate -path migrations -database $DATABASE_URL up
```

5. Run the application:
```bash
make run            # or: go run ./cmd/marketplace
# Override the database via the DATABASE_URL env var (libpq DSN or postgres:// URL).
```

### Common Make targets

```bash
make build        # build the server binary into ./bin
make test         # run all tests with the race detector (needs Docker)
make test-unit    # run only tests that don't require Docker
make lint         # run golangci-lint
make fmt          # format the code (gofmt + goimports)
make migrate-up   # apply migrations against $DATABASE_URL
make sqlc         # regenerate sqlc code
```

### Contributions
Contributions, issues, and feature requests are welcome! Feel free to check the issues page.

### Use This Template

Click **"Use this template"** on GitHub to bootstrap your own service from this structure, or fork it and replace the marketplace domain with your own. The [DDD & CQRS guide](DDD_CQRS_PRINCIPLES.md) walks you through adapting the patterns to any business domain.

If this template helps you, **give it a ⭐** — it helps others find it.

[![Star History Chart](https://api.star-history.com/svg?repos=sklinkert/go-ddd&type=Date)](https://star-history.com/#sklinkert/go-ddd&Date)

### License
Distributed under the MIT License. See LICENSE for more information.
