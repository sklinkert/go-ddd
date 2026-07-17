# FAQ

The questions I get most, answered the way I'd answer them over coffee.

## Isn't DDD overkill for my project?

Sometimes, yes. If your service is a thin CRUD layer — forms in, rows out, no rules beyond "required field" — then entities, value objects, and repository indirection are ceremony without payoff. Use a plain handler-and-database structure and be happy.

The switch flips when *rules* arrive: money, inventory, state machines, permissions, anything where a business person answers "well, it depends" when you ask how it should behave. At that point the question isn't whether the rules get encoded — they will be, somewhere — but whether they live in one guarded place or are smeared across handlers and jobs. My rule of thumb: **if you've written the same `if` in two places, or you can't answer "where is the rule about X?" in ten seconds, you've crossed the threshold.**

And it's not all-or-nothing. A `Money` value object pays for itself in any codebase that touches prices, even one that adopts nothing else from this template.

## DDD vs Clean Architecture vs Hexagonal — which one is this?

They're the same picture at the resolution that matters: dependencies point inward, the domain doesn't know about infrastructure, and the boundary between them is interfaces the inner layer owns. Hexagonal calls them ports and adapters, Clean Architecture calls them use cases and gateways, DDD adds the *modelling* vocabulary — entities, value objects, aggregates — that the other two are largely silent about.

This template is onion-layered in structure and DDD in modelling. If you arrived here fluent in hexagonal: `internal/domain/repositories` are ports, `internal/infrastructure/db/postgres` are adapters, and the DDD content is everything the ports carry.

## Why is there so much code for a two-entity marketplace?

Because the template optimizes for *copying the shape*, not for minimal line count. The marketplace is deliberately trivial so the patterns stay visible; the value is that every pattern is carried to production standard — validation you can't bypass, idempotency that survives races, an outbox that's actually transactional, tests against a real database.

A tutorial that shows 20% of a pattern teaches you to ship 20% of it. The gaps between demo-grade and production-grade (the TTL takeover in idempotency, the partial index on the outbox, `context.WithoutCancel` in cleanup paths) are exactly the parts you can't google your way to when the incident hits.

## Why sqlc instead of GORM or ent?

This template originally used GORM, and the migration to sqlc was a deliberate downgrade in magic. In a DDD codebase the row-to-aggregate mapping is a *boundary* where invariants can leak, so I want it explicit, dumb, and reviewable. sqlc gives compile-time-checked SQL with zero runtime reflection deciding what loads when — no lazy loading, no cascades, no surprise queries. ORMs solve "I don't want to write SQL"; DDD's repositories need the opposite: SQL that does exactly what the aggregate design says and nothing more.

Nothing in the architecture *requires* sqlc, though. The repository interfaces don't know it exists — swapping it out means touching `internal/infrastructure/` only. That's the point of the layering.

## Why do entities have exported fields? Isn't that "anemic"?

The fields are exported because the infrastructure layer needs to read them for persistence, and Go has no `friend` access. What keeps the model non-anemic isn't field privacy — it's that **behavior lives on the entity** (`UpdatePrice`, `AssignSeller`, event recording) and every mutation path re-validates. The `ValidatedProduct` wrapper then makes "went through validation" a compile-time fact at the repository boundary.

Could a colleague assign a field directly and skip validation? Yes, and review should catch it — the same way review catches someone ignoring an error. The pattern makes the safe path the convenient path; it doesn't try to make Go into a language it isn't.

## Where are the domain services / factories / specifications / …?

Not every DDD pattern earns its place in a small domain. The template includes a pattern when the marketplace genuinely exercises it, and omits it when it would be decorative. A domain service (logic spanning multiple aggregates that belongs to no single one) would appear the day a rule like "a seller's total listed value may not exceed X" shows up — and it would be a plain function in the domain package, not a framework.

The same restraint applies to CQRS: no bus, no separate read store, no event sourcing. See [chapter 6](../tutorial/06-cqrs.md) for what's deliberately skipped and why the seams for adding it later are already in place.

## How do I add my own aggregate?

The mechanical checklist, in the order that keeps the compiler helping you:

1. **Domain**: entity + constructor + `validate()` + `ValidatedX` wrapper in `internal/domain/entities/`; events if other systems care; repository interface in `internal/domain/repositories/`.
2. **Schema**: migration pair in `migrations/`, queries in `sql/queries/`, `make sqlc`.
3. **Infrastructure**: `SqlcXRepository` implementing the interface, mapping rows through the constructors.
4. **Application**: commands/queries + service, wrapping writes in `withIdempotency`.
5. **Interface**: DTOs, controller, error mapping entries, routes in `main.go`.
6. **Tests at each layer as you go** — domain tests first; they're the cheapest and they pin the rules before any plumbing exists.

Chapter-by-chapter, that's the whole [tutorial](../tutorial/01-the-domain.md) in reverse.

## Does this scale to microservices / multiple bounded contexts?

The template is one bounded context in one deployable, and that's the honest starting point for almost everyone. The DDD concept that governs the split — the **bounded context** — is about language: when "product" starts meaning different things to different parts of the business, you have two contexts, whether they deploy together or not.

What the template already gives you for that future: aggregates that reference by Id (no shared object graphs to untangle), domain events with an outbox (the integration mechanism between contexts), and a domain layer with no infrastructure entanglement (the part you'd lift out). Extracting a context from here is moving directories, not rewriting models.

## Why `Id` and not `ID`?

House style, chosen and enforced on purpose (the linter rules that would complain are disabled in `.golangci.yml`). `SellerId` composes better than `SellerID` in a codebase full of generated code and JSON tags, and consistency beats the style guide's preference. Adopt the template, keep the convention or flip it — just do it everywhere.

## Something's wrong / missing / could be better

Open an [issue](https://github.com/sklinkert/go-ddd/issues) or a [discussion](https://github.com/sklinkert/go-ddd/discussions) — pattern questions are as welcome as bug reports; there's an issue template specifically for them. And if the template or this tutorial taught you something, [a star](https://github.com/sklinkert/go-ddd) helps the next person find it.
