---
page_title: "cron-utils Provider"
description: |-
  Utility functions for parsing and converting cron expressions between Unix and Quartz formats.
---

# cron-utils Provider

The `cron-utils` provider exposes utility functions for working with cron expressions. It supports both Unix (5-field) and Quartz (6-7 field) formats, and provides parsing, conversion, and schedule metadata.

## Usage

```hcl
terraform {
  required_providers {
    cron-utils = {
      source = "henryupton/cron-utils"
    }
  }
}

provider "cron-utils" {}
```

## Functions

- [`provider::cron-utils::parse`](functions/parse.md) — Parse a cron expression and return schedule metadata
- [`provider::cron-utils::quartz_to_unix`](functions/quartz_to_unix.md) — Convert a Quartz expression to Unix format
- [`provider::cron-utils::unix_to_quartz`](functions/unix_to_quartz.md) — Convert a Unix expression to Quartz format
