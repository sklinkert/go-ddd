# 2. Entities that guard themselves

An **entity** is a domain object with identity and a life cycle. Two products with identical names and prices are still two different products; what makes each one *itself* is its Id, not its attributes. Entities get created, change over time, and eventually get deleted — and at every step, certain things must hold true. Those are its **invariants**.

The anemic version — the one most codebases ship — looks like this:

```go
// The struct anyone can corrupt.
type Product struct {
    Id    uuid.UUID
    Name  string
    Price float64
}
```

Every field public, no constructor, no rules. The invariants exist only in the heads of the developers and in scattered `if` statements across handlers, services, and jobs. Any code can produce a product with an empty name and a negative price, and the compiler will help it do so.

## Rule one: constructors, not struct literals

The template's [`Product`](https://github.com/sklinkert/go-ddd/blob/main/internal/domain/entities/product.go) is built through a constructor that establishes every invariant at birth:

```go
func NewProduct(name string, price Money, seller ValidatedSeller) *Product {
    product := &Product{
        Id:        uuid.Must(uuid.NewV7()),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name:      name,
        Price:     price,
        SellerId:  seller.Id,
    }

    product.recordEvent(events.NewProductCreated(
        product.Id, name, price.MinorUnits(), string(price.Currency()), seller.Id))

    return product
}
```

Three things are decided here, once, for the whole system:

1. **Identity is assigned by the domain.** A UUIDv7 (time-ordered, so it indexes nicely) is generated in the constructor — not by the database, not by the caller. The product has its identity before it ever touches Postgres.
2. **The price is a `Money`**, which as we'll see in [chapter 3](03-value-objects.md) cannot exist in an invalid state.
3. **The seller parameter is a `ValidatedSeller`.** Not a `Seller` — a `ValidatedSeller`. You literally cannot call this function with an unvalidated one; the type doesn't fit.

That third point is the template's signature pattern, so let's take it apart.

## The validated-entity pattern

The problem it solves: in most codebases, "has this struct been validated?" is a question with no answer. A `*Product` in your hand might have come from the constructor, from a JSON unmarshal, from a half-updated cache — you don't know, so defensive code re-validates everywhere, and the checks drift apart.

The template makes validation a *type*:

```go
type ValidatedProduct struct {
    Product
    isValidated bool
}

func NewValidatedProduct(product *Product) (*ValidatedProduct, error) {
    if err := product.validate(); err != nil {
        return nil, err
    }

    return &ValidatedProduct{
        Product:     *product,
        isValidated: true,
    }, nil
}
```

`isValidated` is unexported, so the only way to obtain a `ValidatedProduct` outside the `entities` package is to pass a `Product` through `NewValidatedProduct` — and through `validate()`:

```go
func (p *Product) validate() error {
    if p.Name == "" {
        return fmt.Errorf("%w: name must not be empty", ErrValidation)
    }
    if p.Price.MinorUnits() == 0 {
        return fmt.Errorf("%w: price must be greater than 0", ErrValidation)
    }
    if p.SellerId == uuid.Nil {
        return fmt.Errorf("%w: seller id must not be empty", ErrValidation)
    }
    if p.CreatedAt.After(p.UpdatedAt) {
        return fmt.Errorf("%w: created_at must be before updated_at", ErrValidation)
    }

    return nil
}
```

Now look at what downstream code can demand. The repository interface takes the validated type:

```go
type ProductRepository interface {
    Create(ctx context.Context, product *entities.ValidatedProduct) (*entities.Product, error)
    // ...
}
```

The signature *is* the guarantee: nothing reaches the database without passing validation, and the compiler enforces it. No code review needed, no "did you remember to call Validate()?" comment. Forgetting is a type error.

Note also every check wraps `ErrValidation`, a sentinel from [`errors.go`](https://github.com/sklinkert/go-ddd/blob/main/internal/domain/entities/errors.go). The REST layer maps `errors.Is(err, ErrValidation)` to a 400 without parsing message strings — the domain speaks in errors, the edge translates them.

## Mutation goes through methods

Entities change, and changes must re-establish invariants. So fields are mutated through methods that end in `validate()`:

```go
func (p *Product) UpdatePrice(price Money) error {
    p.Price = price
    p.UpdatedAt = time.Now()

    return p.validate()
}
```

Is this bulletproof? No — `Product` fields are exported (the persistence layer needs them), so a determined colleague can still do `p.Price = Money{}` directly. Go doesn't give us the access control to prevent that entirely without heavy ceremony. The pattern's claim is more modest and, in practice, enough: the *convenient* path and the *reviewed* path are the safe one, and the repository boundary demands the validated type. In several years of running this pattern in production code, "someone bypassed the constructor" has not been the bug. The bug was always in codebases where there was no constructor to bypass.

## What about validation at the API edge?

You might object: my HTTP layer already validates requests. Keep it! Edge validation and domain validation answer different questions:

- The edge asks: *is this request well-formed?* (Is `price_minor_units` a number? Is `seller_id` a UUID?)
- The domain asks: *is this a valid product?* (Is the price positive? Does the seller exist and pass validation?)

The edge check is about protocol; it produces friendly 400s fast. The domain check is about business truth, and it runs no matter where the call came from — HTTP today, a message consumer tomorrow, a backfill script at 2am. The scattered-validation problem isn't solved by choosing one location; it's solved by giving each rule its *one correct* location.

## Try it

1. Delete the `Price.MinorUnits() == 0` check from `validate()` and run `go test ./internal/...` — watch which tests fail and read what they assert. The test suite documents the invariants.
2. Add a new rule: product names must be at most 200 characters. Notice you touch exactly two files — the entity and its test.
3. Try to call `productRepository.Create` with a plain `*Product`. Enjoy the compile error; that error is the pattern working.

Next: [value objects, starting with Money](03-value-objects.md) — the same "invalid states are unconstructible" idea, applied to values instead of identities.
