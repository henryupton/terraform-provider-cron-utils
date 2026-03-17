# Convert a Unix expression to Quartz format
output "quartz_expression" {
  value = provider::cron-utils::unix_to_quartz("*/5 * * * *")
}
