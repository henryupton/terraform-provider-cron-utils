---
page_title: "unix_to_quartz function - cron-utils"
description: |-
  Converts a Unix 5-field cron expression to a Quartz 6-field expression.
---

# Function: `unix_to_quartz`

Converts a Unix 5-field cron expression to a Quartz 6-field expression by prepending a seconds field of `0`.

## Example Usage

```hcl
output "quartz_expression" {
  value = provider::cron-utils::unix_to_quartz("*/5 * * * *")
  # Returns: "0 */5 * * * *"
}
```

```hcl
output "quartz_expression_weekday" {
  value = provider::cron-utils::unix_to_quartz("30 9 * * MON-FRI")
  # Returns: "0 30 9 * * MON-FRI"
}
```

## Signature

```text
unix_to_quartz(expression string) string
```

## Arguments

1. `expression` (String, Required) — A valid Unix cron expression with exactly 5 fields (`minute hour dom month dow`). Supports standard cron syntax including `@yearly`, `@monthly`, `@weekly`, `@daily`, and `@hourly` descriptors.

## Return Value

A String containing the equivalent Quartz 6-field cron expression (`second minute hour dom month dow`), with the seconds field set to `0`.

Returns an error if:
- The expression does not have exactly 5 fields
- The expression is not a valid Unix cron expression
