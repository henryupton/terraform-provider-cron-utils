---
page_title: "quartz_to_unix function - cron-utils"
description: |-
  Converts a Quartz 6-7 field cron expression to a Unix 5-field expression.
---

# Function: `quartz_to_unix`

Converts a Quartz 6-7 field cron expression to a Unix 5-field expression by dropping the seconds field (and year field if present). Returns an error if the expression uses Quartz-specific syntax (`L`, `W`, `#`) that has no Unix equivalent.

## Example Usage

```hcl
output "unix_expression" {
  value = provider::cron-utils::quartz_to_unix("0 */5 * * * ?")
  # Returns: "*/5 * * * *"
}
```

```hcl
output "unix_expression_with_year" {
  value = provider::cron-utils::quartz_to_unix("0 30 9 * * MON-FRI 2026")
  # Returns: "30 9 * * MON-FRI"
}
```

## Signature

```text
quartz_to_unix(expression string) string
```

## Arguments

1. `expression` (String, Required) — A valid Quartz cron expression with 6 or 7 fields (`second minute hour dom month dow [year]`).

## Return Value

A String containing the equivalent Unix 5-field cron expression (`minute hour dom month dow`).

Returns an error if:
- The expression does not have 6 or 7 fields
- The expression is not a valid Quartz cron expression
- Any field uses Quartz-specific syntax (`L`, `W`, `#`) that cannot be represented in Unix format
