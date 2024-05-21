data "neosync_system_transformer" "generate_cc" {
  # This is the value of the TransformerSource enumeration defined in the proto, but with the TRANSFORMER_SOURCE_ prefix removed.
  source = "generate_card_number"
}

output "cc_config" {
  value = data.neosync_system_transformer.generate_cc.config
}
