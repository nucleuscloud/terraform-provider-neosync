data "neosync_system_transformer" "generate_cc" {
  source = "generate_card_number"
}

output "cc_config" {
  value = data.neosync_system_transformer.generate_cc.config
}
