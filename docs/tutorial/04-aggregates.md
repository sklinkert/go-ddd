# 4. Aggregates and their boundaries

This is the chapter where DDD stops being "nice structs" and starts being architecture. **Aggregates** are the answer to two questions every growing codebase eventually faces:

1. When I change this object, what else must change *in the same transaction*?
2. When I load this object, how much of the object graph comes with it?

Get these wrong and you end up in one of two failure modes. Either everything loads everything — fetch a product, get its seller, the seller's other products, their reviews, and a 14-join query the ORM wrote for you — or consistency rules quietly span objects that are saved independently, and you get states the business considers impossible.

## The rule: reference other aggregates by Id

An aggregate is a cluster of objects that change together, with one entity as the **aggregate root** — the only entry point. Everything inside the boundary is loaded together, changed together, and saved in one transaction. Everything *outside* is referenced by identity only.

In the marketplace, `Seller` and `Product` are **separate aggregates**. Look at what `Product` stores:

```go
type Product struct {
    Id        uuid.UUID
    CreatedAt time.Time
    UpdatedAt time.Time
    Name      string
    Price     Money
    SellerId  uuid.UUID   // reference by Id — not: Seller Seller

    domainEvents []events.DomainEvent
}
```

`SellerId uuid.UUID`, not an embedded `Seller`. An earlier version of this template *did* embed the full seller, and it caused exactly the problems the rule predicts:

- **Torn ownership.** Two copies of a seller's name in memory — one on the seller aggregate, one inside each product — and no answer to which one is true after an update.
- **Transactional creep.** Renaming a seller technically modified every product embedding it. Does that update products' `UpdatedAt`? Lock their rows? Nobody chose; the ORM chose.
- **Loading bloat.** You can't fetch a product without dragging its seller along, whether the use case needs it or not.

With an Id reference, each aggregate has one owner, one transaction, one loading story. If a use case needs the seller's details alongside a product, the *application layer* loads both aggregates explicitly — visible, intentional, and cheap to review.

## Boundaries come from invariants, not from foreign keys

The mechanical rule ("reference by Id") is easy. The design question is *where to draw the line*, and the criterion is: **what must be immediately consistent?**

An invariant that must hold at every commit belongs inside one aggregate. A rule that may lag a moment can span aggregates and be reconciled through events or checks at the edges.

Work the marketplace examples:

- *"A product's price must be positive."* Involves only the product. Inside the `Product` aggregate — enforced in `validate()`.
- *"A product must belong to a seller that passed validation."* Spans both. The template enforces it **at creation time** through the type system — `NewProduct` demands a `ValidatedSeller`:

```go
// NewProduct requires a ValidatedSeller so a product can only ever be
// created against a seller that passed validation. The product stores just
// the seller's Id: sellers are a separate aggregate and must not be embedded.
func NewProduct(name string, price Money, seller ValidatedSeller) *Product {
```

Note the honesty in that design: it guarantees the seller was valid *when the product was created*. It does not guarantee the seller still exists a week later — that's a cross-aggregate concern, deliberately not a hard invariant. If the business decides deleting a seller must handle their products, that's a use case (delete them? orphan them? block deletion?), implemented in the application layer or reacting to a `SellerDeleted` domain event. Choosing *eventual* consistency between aggregates isn't a compromise; it's the design telling the truth about what the business actually requires.

- *"Product names must be unique per seller."* The marketplace doesn't have this rule today, but it's the classic trap, so let's work it as a hypothetical. The tempting answer is to make Seller a big aggregate containing all its products, so the invariant sits inside one boundary. Resist it: that turns every product operation into a load-all-products operation, and two sellers adding products concurrently now contend on one aggregate. My answer, if the rule arrived: keep the aggregates small and enforce uniqueness with a database unique index on `(seller_id, name)` — a deliberate, documented exception where infrastructure enforces a domain rule, because it does so atomically and cheaply. Dogma loses to a unique index.

## Small aggregates win

The pull toward big aggregates is real — everything *feels* related. Experience says the opposite: **aggregates should be as small as their invariants allow.** Big aggregates serialize writes (every change locks the root), bloat loads, and turn migrations into archaeology. The marketplace's aggregates are one entity each, and that's not because the domain is trivial; most well-factored aggregates in most systems are one entity plus some value objects.

A checklist I actually use when drawing a boundary:

1. List the invariants. Which ones must be true at *every* commit?
2. Group only what those invariants force together. Everything else: Id reference.
3. Check write contention. Will two users routinely modify this aggregate at once? If yes, the boundary is probably too big.
4. Check the load. If the common read pulls megabytes, same conclusion.

## What the repository layer sees

Aggregate boundaries dictate repository shape, which is why [chapter 5](05-repositories.md) comes next. One repository per aggregate root — `ProductRepository`, `SellerRepository` — and no repository for anything inside a boundary. You never load "a product's events" or "half a seller"; you load aggregates, whole, by their root.

The payoff for the ORM-weary: no lazy loading, no N+1 surprises, no accidentally-saved object graphs. Each repository reads and writes one small cluster, and the SQL underneath (sqlc-generated, in this template) is boring and inspectable.

## Try it

1. Find the commit-worthy invariants for `Seller` in [`seller.go`](https://github.com/sklinkert/go-ddd/blob/main/internal/domain/entities/seller.go). It's a short list — that's a feature.
2. Add a `Review` concept: a buyer reviews a product. Own aggregate or inside `Product`? Work the checklist: what invariant would force reviews into the product's transaction? (I can't find one — reviews are their own aggregate with a `ProductId`.)
3. Sketch what "deleting a seller deletes their products" looks like as an application-layer use case versus a domain event handler. Which one survives the introduction of a second delivery mechanism (say, a CLI admin tool) without duplication?

Next: [repositories](05-repositories.md) — the interface lives in the domain, the SQL lives in infrastructure, and the dependency arrow points the way you want.
