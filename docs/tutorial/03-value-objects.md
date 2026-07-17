# 3. Value objects, starting with Money

Entities have identity; **value objects** have only value. Two `Money` values of 19.99 EUR are not "two moneys that happen to be equal" — they are the same value, the way two 7s are the same number. That definition sounds philosophical until you see what it buys you in code: value objects are immutable, validated at construction, compared by value, and safe to copy and pass anywhere.

Money is the canonical value object because getting it wrong is expensive in the most literal sense. It's also the bug I've seen most often in the wild: a `Price float64` field, and a team that swears the numbers "look fine".

## Why float64 money is broken

Not risky — broken. Binary floating point can't represent most decimal fractions exactly. The party trick is `0.1 + 0.2 == 0.3` being `false`, but the version that hurts is drift under accumulation:

```go
var total float64
for i := 0; i < 1000; i++ {
    total += 0.10 // a 10-cent fee, a thousand times
}
fmt.Println(total)        // 99.9999999999986
fmt.Println(total == 100) // false
```

A thousand 10-cent charges and exact comparison is already gone. Feed that into a `>= threshold` check or a reconciliation diff against your payment provider and you're in epsilon-comparison whack-a-mole.

And `float64` has a second, quieter problem: **it's a bare number**. `19.99` of what? EUR? USD? Cents already? I've debugged a production incident where one service sent euros and another read the same field as cents. The type system waved the 100x pricing error through, because `float64 == float64`.

## The Money value object

The fix, from [`internal/domain/entities/money.go`](https://github.com/sklinkert/go-ddd/blob/main/internal/domain/entities/money.go): integer minor units, currency attached, both fields unexported.

```go
type Currency string

const (
    EUR Currency = "EUR"
    USD Currency = "USD"
)

var supportedCurrencies = map[Currency]struct{}{
    EUR: {},
    USD: {},
}

// Money is an immutable value object storing an amount in minor units
// (cents) to avoid floating-point rounding errors.
type Money struct {
    cents    int64
    currency Currency
}

func NewMoney(cents int64, currency Currency) (Money, error) {
    if cents < 0 {
        return Money{}, fmt.Errorf("%w: amount must not be negative", ErrValidation)
    }
    if _, ok := supportedCurrencies[currency]; !ok {
        return Money{}, fmt.Errorf("%w: unsupported currency %q", ErrValidation, currency)
    }

    return Money{cents: cents, currency: currency}, nil
}

func (m Money) Cents() int64       { return m.cents }
func (m Money) Currency() Currency { return m.currency }
```

The design decisions, spelled out:

**Unexported fields, one constructor.** Because `cents` and `currency` are unexported, the only way to build a `Money` outside the package is `NewMoney`. Every `Money` in the entire system has passed the negative check and the currency whitelist. There is no "construct raw, validate later" path to forget. This is [chapter 2's](02-entities.md) idea taken to its logical end: for a value object, we *can* make invalid states fully unrepresentable, because nothing needs to mutate it.

**Integer minor units, not `big.Rat` or a decimal library.** `int64` cents gives exact addition and comparison for free, sorts and indexes trivially in the database, and covers about ±92 quadrillion dollars. For prices and balances, minor units are the boring answer that works. (For FX rates or interest math, reach for a decimal type — different problem.)

**Currency is part of equality.** `Money` is a comparable struct, so `==` compares amount *and* currency: `NewMoney(1000, EUR) != NewMoney(1000, USD)`. The euros-versus-cents class of bug now fails at the type level instead of on an accountant's spreadsheet.

**A whitelist of currencies.** Two entries looks restrictive; that's the point. The *domain* decides which currencies the business supports. Adding one is a one-line change that forces a moment of thought — which is what you want when the alternative is silently accepting `"BTC"` or `"EURO"` from a client.

## Immutability changes how arithmetic looks

There's no `SetCents`. If the template grew an `Add`, it would return a *new* value:

```go
func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, fmt.Errorf("%w: cannot add %s to %s", ErrValidation, other.currency, m.currency)
    }
    return NewMoney(m.cents+other.cents, m.currency)
}
```

Note what the signature forces you to decide: what does `EUR + USD` mean? My answer is "an error — conversion is an explicit domain operation with a rate and a timestamp, never implicit". You may answer differently, but the value object made you answer *once*, in one place, instead of everywhere an addition happens.

## Surviving the edges: JSON

A value object is only as good as its boundaries. Money constantly crosses process edges — API, database, outbox messages — and every crossing is a chance to smuggle in an invalid value. The subtle one is JSON: `encoding/json` normally writes straight into struct fields via reflection, bypassing your constructor. With unexported fields it can't, so the template implements the interfaces and routes deserialization *back through* `NewMoney`:

```go
// UnmarshalJSON goes through NewMoney so a Money can never be deserialized
// into an invalid state.
func (m *Money) UnmarshalJSON(data []byte) error {
    var raw moneyJSON
    if err := json.Unmarshal(data, &raw); err != nil {
        return err
    }

    money, err := NewMoney(raw.Cents, raw.Currency)
    if err != nil {
        return err
    }

    *m = money
    return nil
}
```

Why bother, when the value was valid at serialization time? Because JSON doesn't only come from you. It comes from an outbox row written by last year's code, a queue message from another service, a fixture someone hand-edited. Revalidating on the way in costs two comparisons and closes the whole category.

At the REST edge, the DTO doesn't expose the domain type at all — it carries `price_cents` and `currency` as explicit primitive fields. The name `price_cents` does real work: no client developer will ever wonder whether to send `19.99` or `1999`.

## The sharp edges I'll tell you about myself

- **`String()` assumes two decimal places.** Correct for EUR and USD, wrong for JPY (0 decimals) or KWD (3). The whitelist keeps the assumption safe *by construction* — but if you go properly multi-currency, "cents" itself becomes a misleading word. The accurate term is *minor units*, and you'll need per-currency exponents from ISO 4217.
- **Division still rounds.** Integer cents make addition exact, but 100 cents split three ways is 33+33+34 no matter what. You need an allocation strategy (largest-remainder works), and if you allocate, *persist the allocation* — a refund must mirror the original split, not recompute it.
- **Parsing at ingestion boundaries is where money bugs actually live.** A bank API sending `"5000"` JPY means 5000 yen, not 50.00 of anything. Parse currency-aware, before the value ever becomes an integer in your system.

## Beyond Money

The pattern generalizes to anything defined by its value and burdened with rules: email addresses, percentages, date ranges, country codes, quantities-with-units. The test for "should this be a value object?" is simple: *do I keep validating this same shape in multiple places?* If yes, give it a constructor and a type, and delete the scattered checks.

## Try it

1. Add `GBP` to the whitelist and write the test. One line plus one assertion — feel how cheap extending a value object is.
2. Write `Sub(other Money)` and decide what happens when the result would be negative. (There's no universally right answer: an account balance may go negative, a price may not. Your domain decides.)
3. Round-trip a `Money` through `json.Marshal`/`json.Unmarshal`, then hand-edit the JSON to `"currency": "XXX"` and unmarshal again. Watch the constructor reject it.

Next: [aggregates and their boundaries](04-aggregates.md) — why `Product` holds a `SellerId` and not a `Seller`.
