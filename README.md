# terraform-provider-cron-utils

A Terraform provider with utility functions for parsing and converting cron expressions between Unix and Quartz formats.

## Functions

| Function | Description |
|---|---|
| `provider::cron-utils::parse` | Parse a cron expression and return schedule metadata |
| `provider::cron-utils::quartz_to_unix` | Convert a Quartz expression to Unix format |
| `provider::cron-utils::unix_to_quartz` | Convert a Unix expression to Quartz format |

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

output "schedule" {
  value = provider::cron-utils::parse("*/5 * * * *")
}
```

## Supported Formats

- **Unix** — 5-field expressions: `minute hour dom month dow` (e.g. `*/5 * * * *`)
- **Quartz** — 6-7 field expressions: `second minute hour dom month dow [year]` (e.g. `0 */5 * * * ?`)
