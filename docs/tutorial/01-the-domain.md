# 1. The domain and its language

Before a single line of Go, DDD asks one thing of you: agree on what the words mean.

That sounds trivial. It never is. I've sat in meetings where "order" meant three different things to three different teams — the cart before checkout, the paid transaction, and the fulfillment job — and the codebase faithfully reflected the confusion with `Order`, `OrderV2`, and `PurchaseOrder` structs that half-overlapped.

DDD calls the fix the **ubiquitous language**: one vocabulary, shared by developers and domain experts, used *everywhere* — in conversation, in tickets, and crucially in code. If the business says "seller", the struct is `Seller`, the table is `sellers`, the endpoint is `/sellers`. No `Vendor` in the API and `Merchant` in the database because two developers had different tastes.

## The marketplace we're modelling

The template models the smallest domain that still exercises real DDD patterns: a marketplace.

- A **Seller** is someone who lists things for sale. A seller has a name.
- A **Product** is something a seller offers. A product has a name, a price, and belongs to exactly one seller.
- A **Price** is money: an amount *and* a currency. "19.99" is not a price; "19.99 EUR" is.

Deliberately boring. The point of the template is the patterns, not the domain — but even this tiny domain forces the interesting decisions: Can a product exist without a seller? (No.) Can a price be zero? (For a product, no.) What happens to products when a seller is deleted? Every one of those questions is a *business* question, and the code should answer it in exactly one place.

## Where the language lives in the code

Look at the package layout under [`internal/`](https://github.com/sklinkert/go-ddd/tree/main/internal):

```
internal/
├── domain/           # the model: entities, value objects, events, repository interfaces
│   ├── entities/     # Product, Seller, Money, ...
│   ├── events/       # ProductCreated, ...
│   └── repositories/ # ProductRepository, SellerRepository (interfaces only)
├── application/      # use cases: commands, queries, services
├── infrastructure/   # Postgres, sqlc, the outbox relay — the details
└── interface/        # REST controllers, DTOs — the delivery mechanism
```

The `domain` package is the heart, and it's the layer with the fewest imports. Open [`internal/domain/entities/product.go`](https://github.com/sklinkert/go-ddd/blob/main/internal/domain/entities/product.go) and check its import block: standard library, `google/uuid`, and the domain's own `events` package. No Echo, no pgx, no JSON tags. The domain doesn't know HTTP or SQL exist.

That's not an aesthetic preference. It's what makes the language *stable*. HTTP frameworks and database libraries churn; "a product belongs to a seller" doesn't. When the model is free of infrastructure, the code that encodes business knowledge survives every migration — I've moved this template from GORM to sqlc without touching a domain rule.

## Naming is design

A few naming decisions in the template worth noticing, because each encodes a rule:

- `NewProduct(name string, price Money, seller ValidatedSeller)` — the constructor takes a `ValidatedSeller`, not a `Seller`. The type signature says: *you cannot attach a product to a seller that hasn't passed validation.* More on this in [chapter 2](02-entities.md).
- `Money`, not `float64` — a price without a currency is a bug waiting for an exchange rate. [Chapter 3](03-value-objects.md) is all about this.
- `SellerId uuid.UUID` on `Product`, not `Seller Seller` — products reference sellers by identity; they don't own them. That's an aggregate boundary, and it's [chapter 4](04-aggregates.md).

The pattern behind all three: **make the language do work**. Every time a business rule can be expressed as a type instead of a comment, the compiler becomes the reviewer who never gets tired.

## Try it

Clone the repo and find the answers in the code — each should take under a minute, which is itself the point:

```bash
git clone https://github.com/sklinkert/go-ddd.git && cd go-ddd
```

1. What are *all* the rules for a valid product? (One function answers this: `validate()` in `internal/domain/entities/product.go`.)
2. Can any code construct a `Money` with a negative amount? (Try it: write a snippet importing `entities` and see what the compiler and the constructor let you do.)
3. Where would a new rule like "product names must be unique per seller" go? Think about it now; [chapter 4](04-aggregates.md) gives my answer.

Next: [entities that guard themselves](02-entities.md).
