# Convert a Quartz expression to Unix format
output "unix_expression" {
  value = provider::cron-utils::quartz_to_unix("0 */5 * * * ?")
}
