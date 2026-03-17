---
page_title: "parse function - cron-utils"
description: |-
  Parses a Unix or Quartz cron expression and returns schedule metadata including next execution times and regularity.
---

# Function: `parse`

Parses a Unix (5-field) or Quartz (6-7 field) cron expression and returns schedule metadata. The format is auto-detected by field count.

## Example Usage

```hcl
output "schedule" {
  value = provider::cron-utils::parse("*/5 * * * *")
}

# Returns:
# {
#   expression_type     = "unix"
#   next_execution      = "2026-03-17T10:05:00Z"
#   next_execution_unix = 1742205900
#   is_regular          = true
#   interval_seconds    = 300
#   next_executions     = ["2026-03-17T10:05:00Z", "2026-03-17T10:10:00Z", ...]
# }
```

```hcl
output "quartz_schedule" {
  value = provider::cron-utils::parse("0 */15 9-17 * * MON-FRI")
}
```

## Signature

```text
parse(expression string) object
```

## Arguments

1. `expression` (String, Required) — A cron expression in Unix 5-field or Quartz 6-7 field format.
   - **Unix**: `"*/5 * * * *"` — fields are `minute hour dom month dow`
   - **Quartz**: `"0 */5 * * * ?"` — fields are `second minute hour dom month dow [year]`

## Return Value

An object with the following attributes:

| Attribute | Type | Description |
|---|---|---|
| `expression_type` | String | Format detected: `"unix"` or `"quartz"` |
| `next_execution` | String | RFC3339 timestamp of the next scheduled execution |
| `next_execution_unix` | Number | Unix timestamp (seconds) of the next scheduled execution |
| `is_regular` | Bool | `true` if all intervals between the next 13 executions are equal |
| `interval_seconds` | Number | Seconds between executions when `is_regular` is `true`, otherwise `0` |
| `next_executions` | List of String | RFC3339 timestamps of the next 5 scheduled executions |
