---
name: moneybags
description: Use when querying the user's personal finance records through the moneybags MCP server — "what did I spend on coffee", "list my reimbursements", "find transactions at Starbucks", "how much was my tithe". Covers the search_transactions tool, the transaction type values, and the data quirks that make naive queries silently return the wrong answer.
---

# Moneybags MCP

Moneybags is a personal finance tracker for a Singapore-based user's spending. It is
populated mostly by the user typing free-text notes into a Telegram bot, which is
why the data is free-form and inconsistent.

## Scope

`search_transactions` is currently the **only** tool exposed over MCP. Portfolio,
income, mortgage, stock trades, assets and car data are **not** reachable through
this server.

If a question needs that data, say the data is not available here rather than
guessing or inventing figures. Do not use a different tool to approximate it.

## Connection

The server speaks streamable HTTP at:

```
http://<moneybags-host>:6000/mcp
```

Sessions are stateless. There is **no authentication**, so the port is expected to
sit on a private network, VPN, or SSH tunnel.

## Picking arguments

Both arguments are optional. Omitting both returns the 100 most recent
transactions.

### `description`

A **case-insensitive substring** match against the transaction's free-text note. It
is not a category, a merchant field, or a keyword index — the notes are whatever
the user happened to type, so the same merchant can appear as `Starbucks Coffee`,
`starbucks`, or `coffee @ starbucks`.

**Prefer short, distinctive, lowercase substrings.** `coffee` will find far more
than `starbucks coffee`, because a longer phrase has to appear verbatim and
contiguously.

`%` and `_` are matched literally, not as wildcards. Searching `100%` finds notes
containing `100%` and will not match `100X`.

### `type`

Restricts to a single transaction type. Matching is lenient: casing, `_`, `-`,
repeated whitespace and the Telegram shorthands all normalise, so `shared cc reim`,
`SHARED_CC_REIMBURSE` and `Shared CC Reimburse` are the same type.

Pass `type` whenever you know it. It is far more reliable than guessing at note
wording.

### The four types a description search can never find

`CREDIT CARD`, `INSURANCE`, `TITHE` and `TAX` **always have an empty description**.
The write path never accepts a note for these types, so their transactions have
nothing to match against.

A description search therefore **silently misses all of them**. Always reach for
`type` when the question concerns these:

| Instead of | Use |
|---|---|
| `{"description": "insurance"}` | `{"type": "insurance"}` |
| `{"description": "tithe"}` | `{"type": "tithe"}` |
| `{"description": "credit card"}` | `{"type": "credit card"}` |
| `{"description": "tax"}` | `{"type": "tax"}` |

## Reading the response

```json
{
  "transactions": [
    {
      "id": 3,
      "date": "2023-04-05",
      "type": "SHARED_CC_REIMBURSE",
      "description": "COFFEE beans",
      "amount": 12
    }
  ],
  "count": 1,
  "truncated": false
}
```

- `date` is `yyyy-mm-dd` **in Singapore time**, already converted server-side. Do
  not shift it again.
- `amount` is stored exactly as it was entered; moneybags does not normalise the
  sign. Read the actual values rather than assuming one.
- `id` is what the Telegram bot uses to delete a transaction.

### `truncated`

Results are capped at **100 transactions**, newest first.

When `truncated` is `true`, more transactions matched than were returned, and the
result is only the most recent slice of them. On `truncated: true`:

- Do **not** sum the amounts and present the result as a total.
- Do **not** describe the result as "all" or "every" matching transaction.
- Say "at least N" or narrow the query with a `type` filter.

## Known limitation: no date filter

There is no date argument, so **"what did I spend in March?" is not answerable by
the server.**

You can filter client-side on the returned `date` values, but only when
`truncated` is `false` — otherwise the window is incomplete and the answer would be
wrong. When the result is truncated and the user asked about a period, say the
search returned more matches than could be shown rather than presenting a partial
figure as the answer.

## Worked examples

| Question | Call | Note |
|---|---|---|
| "How much did I spend on coffee?" | `{"description": "coffee"}` | Sum `amount`; only call it a total if `truncated` is false. |
| "Find my Starbucks transactions" | `{"description": "starbucks"}` | Try `coffee` as a fallback if empty. |
| "List my reimbursements" | `{"type": "reimburse"}` | Type filter, no description. |
| "What did I pay for insurance?" | `{"type": "insurance"}` | **Never** a description search; the note is always empty. |
| "What's my biggest shared cc reim expense?" | `{"type": "shared cc reim"}` | Sort by `amount` client-side. |
| "Anything about groceries?" | `{"description": "grocery"}` | Try `groceries` too — note wording is inconsistent. |

## Failure modes

**Unknown type.** The call fails with `isError: true` and a message listing every
valid type. Read that list and retry. Do not guess type strings repeatedly.

**Empty result.** Usually a wording mismatch, not proof that nothing was spent.
Before concluding there were no matching transactions:

1. Retry with a shorter substring (`starbucks coffee` → `coffee`).
2. If the question concerns spending category, drop the description entirely and
   filter by `type`.
3. If the question concerns one of the four always-empty types, you almost
   certainly forgot the `type` filter.

Only report "no matching transactions" once the shorter and type-filtered variants
have also come back empty.

## Type reference

| Type | Also accepted | Meaning |
|---|---|---|
| `OWN` | `own` | Regular spending for the user. |
| `REIM` | `reimburse` | Paid upfront on a credit card, to be reimbursed. |
| `SHARED` | `shared` | Shared expense the other party paid first. |
| `SHARED REIM` | | Paid upfront, reimbursed from the shared account. |
| `SPECIAL SHARED` | | One-off shared expense the other party paid first. |
| `SPECIAL SHARED REIM` | | As `SHARED REIM`, but excluded from regular spend totals. |
| `SPECIAL OWN` | | One-off spending for the user. |
| `CREDIT CARD` | `cc` | Paid using a credit card. Note always empty. |
| `SHARED CC REIMBURSE` | `shared cc reim` | Shared credit card used personally, to be reimbursed. |
| `INSURANCE` | `insurance` | Insurance payment. Note always empty. |
| `TITHE` | `tithe` | Given to parents. Note always empty. |
| `TAX` | `tax` | Paid to the tax authority. Note always empty. |