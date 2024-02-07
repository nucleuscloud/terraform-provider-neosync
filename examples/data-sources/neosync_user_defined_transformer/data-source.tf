data "neosync_user_defined_transformer" "my_transformer" {
  id = "3b83d1d3-5ffe-48c6-ac11-7a2e60802864"
}

output "my_id" {
  value = data.neosync_user_defined_transformer.my_transformer.id
}
