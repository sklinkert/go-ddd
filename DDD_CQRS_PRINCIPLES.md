# Domain-Driven Design (DDD) and CQRS Principles Guide

> **This guide has grown into a full tutorial site: [sklinkert.github.io/go-ddd](https://sklinkert.github.io/go-ddd/)**
>
> Everything that used to live in this file is now covered there in more depth, with every pattern anchored to the actual code in this repository.

Where to find what:

| Topic | Now lives at |
|---|---|
| Why DDD, core concepts | [Why DDD makes your life easier](https://sklinkert.github.io/go-ddd/) |
| Ubiquitous language, the marketplace domain | [Tutorial 1: The domain and its language](https://sklinkert.github.io/go-ddd/tutorial/01-the-domain/) |
| Entity design, validation, the validated-entity pattern | [Tutorial 2: Entities that guard themselves](https://sklinkert.github.io/go-ddd/tutorial/02-entities/) |
| Value objects, the `Money` implementation | [Tutorial 3: Value objects, starting with Money](https://sklinkert.github.io/go-ddd/tutorial/03-value-objects/) |
| Aggregates, boundaries, reference by Id | [Tutorial 4: Aggregates and their boundaries](https://sklinkert.github.io/go-ddd/tutorial/04-aggregates/) |
| Repository interfaces and the sqlc implementation | [Tutorial 5: Repositories](https://sklinkert.github.io/go-ddd/tutorial/05-repositories/) |
| CQRS: commands, queries, results | [Tutorial 6: CQRS, commands and queries](https://sklinkert.github.io/go-ddd/tutorial/06-cqrs/) |
| Domain events and the transactional outbox | [Tutorial 7: Domain events and the outbox](https://sklinkert.github.io/go-ddd/tutorial/07-domain-events-outbox/) |
| Race-safe idempotency keys | [Tutorial 8: Idempotent commands](https://sklinkert.github.io/go-ddd/tutorial/08-idempotency/) |
| Testing strategy per layer | [Tutorial 9: Testing a DDD codebase](https://sklinkert.github.io/go-ddd/tutorial/09-testing/) |
| Architecture layers, conventions, best practices | [Architecture reference](https://sklinkert.github.io/go-ddd/reference/architecture/) |
| Applying the patterns to your own domain, trade-offs | [FAQ](https://sklinkert.github.io/go-ddd/reference/faq/) |

Prefer reading in the repo? The same pages are plain markdown under [`docs/`](docs/).

For the outbox pattern as a standalone, broker-agnostic library, see [go-outbox](https://github.com/sklinkert/go-outbox).
