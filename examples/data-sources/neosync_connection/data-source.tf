data "neosync_connection" "foo" {
  id = "3b83d1d3-5ffe-48c6-ac11-7a2e60802864"
}

output "connection_id" {
  value = data.neosync_connection.foo.id
}

output "connection_name" {
  value = data.neosync_connection.foo.name
}
