# Domain-Driven Design (DDD) and CQRS Principles Guide

This guide explains the DDD and CQRS principles demonstrated in this codebase, presented in a way that allows engineers to apply these patterns to any industry or business domain.

## Table of Contents
1. [Domain-Driven Design Overview](#domain-driven-design-overview)
2. [Architecture Layers](#architecture-layers)
3. [Domain Layer Principles](#domain-layer-principles)
4. [Application Layer Patterns](#application-layer-patterns)
5. [Infrastructure Layer Design](#infrastructure-layer-design)
6. [CQRS Implementation](#cqrs-implementation)
7. [Idempotency Pattern](#idempotency-pattern)
8. [Best Practices](#best-practices)
9. [Applying These Principles](#applying-these-principles)

## Domain-Driven Design Overview

Domain-Driven Design (DDD) is a software development approach that:
- Places the business domain at the heart of the software
- Creates a shared language between developers and domain experts
- Isolates business logic from technical concerns
- Enables scalable and maintainable architectures

### Core Concepts

1. **Ubiquitous Language**: Use the same terminology in code that domain experts use
2. **Bounded Contexts**: Define clear boundaries where specific domain models apply
3. **Entities**: Objects with unique identity that persist over time
4. **Value Objects**: Immutable objects defined by their attributes
5. **Aggregates**: Clusters of entities and value objects with defined boundaries
6. **Repositories**: Abstractions for data persistence

## Architecture Layers

This implementation follows the Onion Architecture pattern with clear separation of concerns:

```
┌─────────────────────────────────────┐
│         Interface Layer             │ ← External APIs, Controllers
├─────────────────────────────────────┤
│       Application Layer             │ ← Use Cases, Commands, Queries
├─────────────────────────────────────┤
│         Domain Layer                │ ← Business Logic, Entities
├─────────────────────────────────────┤
│     Infrastructure Layer            │ ← Database, External Services
└─────────────────────────────────────┘
```

### Layer Dependencies
- Inner layers know nothing about outer layers
- Domain layer has zero dependencies on other layers
- Infrastructure implements interfaces defined by domain
- Application layer orchestrates between domain and infrastructure

## Domain Layer Principles

### 1. Entity Design
```go
type Entity struct {
    Id        uuid.UUID    // Always use unique identifiers
    CreatedAt time.Time    // Set by domain, not database
    UpdatedAt time.Time    // Updated on modifications
    // Business attributes
}
```

**Key Principles:**
- Entities have identity that persists across time
- Factory methods (NewEntity) ensure valid initial state
- Business rules are enforced through methods
- Validation happens at creation and modification

### 2. Validation Pattern
```go
// Private validation method
func (e *Entity) validate() error {
    // Business rule validations
    if e.BusinessAttribute == "" {
        return errors.New("business rule violation")
    }
    return nil
}

// Public modification method with validation
func (e *Entity) UpdateAttribute(value string) error {
    e.BusinessAttribute = value
    e.UpdatedAt = time.Now()
    return e.validate()
}
```

### 3. Validated Entity Pattern
```go
type ValidatedEntity struct {
    Entity
    isValidated bool
}

func NewValidatedEntity(entity *Entity) (*ValidatedEntity, error) {
    if err := entity.validate(); err != nil {
        return nil, err
    }
    return &ValidatedEntity{
        Entity:      *entity,
        isValidated: true,
    }, nil
}
```

**Purpose:** Ensures only valid entities can be persisted

### 4. Repository Interfaces
```go
type EntityRepository interface {
    Create(entity *ValidatedEntity) (*Entity, error)
    FindById(id uuid.UUID) (*Entity, error)
    FindAll() ([]*Entity, error)
    Update(entity *ValidatedEntity) (*Entity, error)
    Delete(id uuid.UUID) error
}
```

**Key Points:**
- Domain defines interfaces, infrastructure implements them
- Methods accept validated entities for writes
- Read methods return regular entities (not validated)
- Always return fresh data after writes

## Application Layer Patterns

### 1. Service Structure
```go
type EntityService struct {
    repo            repositories.EntityRepository
    idempotencyRepo repositories.IdempotencyRepository
}
```

Services orchestrate:
- Command execution
- Query handling
- Transaction boundaries
- Cross-aggregate operations

### 2. Use Case Implementation
- Each use case is a method on the service
- Clear input (commands) and output (results)
- Handles idempotency
- Coordinates between repositories

## Infrastructure Layer Design

### 1. Repository Implementation
```go
type SqlcEntityRepository struct {
    queries *db.Queries
}

func (repo *SqlcEntityRepository) Create(entity *ValidatedEntity) (*Entity, error) {
    ctx := context.Background()
    dbEntity, err := repo.queries.CreateEntity(ctx, db.CreateEntityParams{
        ID:        entity.Id,
        Name:      entity.Name,
        CreatedAt: timestamptzFromTime(entity.CreatedAt),
        UpdatedAt: timestamptzFromTime(entity.UpdatedAt),
    })
    if err != nil {
        return nil, err
    }
    // Always read after write
    return repo.FindById(dbEntity.ID)
}
```

### 2. Mapping Pattern
```go
// Domain to Database
func toDBModel(entity *ValidatedEntity) *DBModel {
    // Map domain entity to database model
}

// Database to Domain
func fromDBModel(dbModel *DBModel) *Entity {
    // Map database model to domain entity
}
```

**Purpose:** Keep domain models pure and database concerns isolated

## CQRS Implementation

### Command Pattern
Commands modify state and are task-oriented:

```go
type CreateEntityCommand struct {
    IdempotencyKey string
    // Business attributes
}

type CreateEntityCommandResult struct {
    Result *EntityResult
}
```

### Query Pattern
Queries retrieve data without side effects:

```go
type EntityQueryResult struct {
    Result *EntityResult
}

type EntityQueryListResult struct {
    Result []*EntityResult
}
```

### Benefits of CQRS
1. **Optimized Read/Write Models**: Different models for different purposes
2. **Scalability**: Scale reads and writes independently
3. **Performance**: Optimize queries without affecting write logic
4. **Clarity**: Clear separation of intentions

## Idempotency Pattern

### Implementation
```go
// Check for existing execution
if command.IdempotencyKey != "" {
    existing, err := idempotencyRepo.FindByKey(ctx, command.IdempotencyKey)
    if existing != nil {
        return cachedResponse, nil
    }
}

// Execute business logic
result := executeBusinessLogic()

// Store result for future requests
if command.IdempotencyKey != "" {
    record := NewIdempotencyRecord(command.IdempotencyKey, request)
    record.SetResponse(response, statusCode)
    idempotencyRepo.Create(ctx, record)
}
```

### Benefits
- Prevents duplicate operations
- Handles network failures gracefully
- Ensures consistency in distributed systems

## Best Practices

### 1. Domain Layer Purity
- No framework dependencies
- No infrastructure concerns
- Business logic only
- Self-contained validation

### 2. Factory Methods
```go
func NewEntity(businessAttribute string) *Entity {
    return &Entity{
        Id:                uuid.New(),
        CreatedAt:         time.Now(),
        UpdatedAt:         time.Now(),
        BusinessAttribute: businessAttribute,
    }
}
```

### 3. Read After Write
Always return fresh data from the database after modifications to ensure consistency.


### 4. Historical Data Compatibility
- Don't validate on read operations
- Allow loading of data created with old business rules
- Validate only on write operations

## Applying These Principles

### Step 1: Identify Your Domain
1. Work with domain experts to understand the business
2. Identify key entities and their relationships
3. Define business rules and invariants
4. Create a ubiquitous language

### Step 2: Design Your Entities
```go
// Example for an e-commerce domain
type Order struct {
    Id         uuid.UUID
    CreatedAt  time.Time
    UpdatedAt  time.Time
    CustomerId uuid.UUID
    Items      []OrderItem
    Status     OrderStatus
    Total      Money
}

func NewOrder(customerId uuid.UUID) *Order {
    return &Order{
        Id:         uuid.New(),
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
        CustomerId: customerId,
        Status:     OrderStatusPending,
        Items:      []OrderItem{},
    }
}

func (o *Order) AddItem(product Product, quantity int) error {
    // Business logic for adding items
    // Validate quantity, calculate prices, etc.
}
```

### Step 3: Define Repository Interfaces
```go
type OrderRepository interface {
    Create(order *ValidatedOrder) (*Order, error)
    FindById(id uuid.UUID) (*Order, error)
    FindByCustomerId(customerId uuid.UUID) ([]*Order, error)
    Update(order *ValidatedOrder) (*Order, error)
}
```

### Step 4: Implement CQRS
Commands:
```go
type PlaceOrderCommand struct {
    IdempotencyKey string
    CustomerId     uuid.UUID
    Items          []OrderItemRequest
}
```

Queries:
```go
type GetCustomerOrdersQuery struct {
    CustomerId uuid.UUID
    Status     *OrderStatus // Optional filter
}
```

### Step 5: Create Application Services
```go
type OrderService struct {
    orderRepo       repositories.OrderRepository
    productRepo     repositories.ProductRepository
    idempotencyRepo repositories.IdempotencyRepository
}

func (s *OrderService) PlaceOrder(cmd *PlaceOrderCommand) (*PlaceOrderResult, error) {
    // Implement idempotency check
    // Validate products exist
    // Create order
    // Calculate totals
    // Save order
    // Return result
}
```

### Industry-Agnostic Guidelines

1. **Healthcare**: Patient (Entity), Appointment (Entity), Diagnosis (Value Object)
2. **Finance**: Account (Entity), Transaction (Entity), Money (Value Object)
3. **Education**: Student (Entity), Course (Entity), Grade (Value Object)
4. **Logistics**: Shipment (Entity), Package (Entity), Address (Value Object)

The patterns remain the same; only the domain concepts change.

## Conclusion

These DDD and CQRS principles provide a robust foundation for building maintainable, scalable applications regardless of your business domain. The key is to:

1. Keep domain logic pure and isolated
2. Use CQRS to separate read and write concerns
3. Implement proper validation and factory patterns
4. Handle idempotency for distributed systems
5. Follow the dependency rules between layers

By applying these principles, you create software that clearly expresses business requirements while remaining flexible for future changes.