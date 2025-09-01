# Go-DDD: Domain Driven Design Template in Golang

Welcome to `go-ddd`, a reference implementation/template repository demonstrating the [Domain Driven Design (DDD)](https://en.wikipedia.org/wiki/Domain-driven_design) and CQRS (Command Query Responsibility Segregation) approach in Golang. This project aims to help developers and architects understand the DDD structure, especially in the context of Go, and how it can lead to cleaner, enterprise-ready, maintainable, and scalable codebases.

## Overview

Domain-Driven Design is a methodology and design pattern used to build complex enterprise software by connecting the implementation to an evolving model. `go-ddd` showcases this by setting up a simple marketplace where `Sellers` can sell `Products`.

### Why DDD?

- **Ubiquitous Language**: Promotes a common language between developers and stakeholders.
- **Isolation of Domain Logic**: The domain logic is separate from the infrastructure and application layers, promoting SOLID principles.
- **Scalability**: Allows for easier microservices architecture transitions.

## Documentation

📚 **[Comprehensive DDD & CQRS Principles Guide](DDD_CQRS_PRINCIPLES.md)** - Learn how to apply these patterns to any business domain.

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
- The application layer checks if this key has been processed before
- If yes, it returns the cached response without re-executing business logic
- If no, it executes the command and stores the response for future requests

This prevents duplicate entities from being created when clients retry failed requests.

## Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema management. Migrations are stored in the `migrations/` directory with sequential version numbers.

### Migration Files Structure
```
migrations/
├── 000001_initial_schema.up.sql    # Creates initial tables
└── 000001_initial_schema.down.sql  # Rollback for initial schema
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
- `000002_add_user_email_column.up.sql` - Forward migration
- `000002_add_user_email_column.down.sql` - Rollback migration

### Migration Best Practices
- Always create both `up` and `down` migrations
- Test migrations on a copy of production data
- Keep migrations small and focused
- Never modify existing migration files once they've been applied in production
- Use descriptive names for migration files

## Tech Stack

This project uses the following key dependencies:

- **golang-migrate** (`github.com/golang-migrate/migrate/v4`) - Database migration tool
- **sqlc** - Type-safe SQL code generation for PostgreSQL
- **pgx/v5** (`github.com/jackc/pgx/v5`) - PostgreSQL driver and toolkit
- **Echo** (`github.com/labstack/echo/v4`) - HTTP web framework
- **UUID** (`github.com/google/uuid`) - UUID generation
- **Testify** (`github.com/stretchr/testify`) - Testing toolkit

## Getting Started

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
go run ./cmd/marketplace
```

### Contributions
Contributions, issues, and feature requests are welcome! Feel free to check the issues page.

### License
Distributed under the MIT License. See LICENSE for more information.
