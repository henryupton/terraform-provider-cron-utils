# Parse a Unix cron expression
output "unix_schedule" {
  value = provider::cron-utils::parse("*/5 * * * *")
}

# Parse a Quartz cron expression
output "quartz_schedule" {
  value = provider::cron-utils::parse("0 30 9 * * MON-FRI")
}
