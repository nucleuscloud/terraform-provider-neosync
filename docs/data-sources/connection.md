---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "neosync_connection Data Source - terraform-provider-neosync"
subcategory: ""
description: |-
  Neosync Connection data source
---

# neosync_connection (Data Source)

Neosync Connection data source

## Example Usage

```terraform
data "neosync_connection" "foo" {
  id = "3b83d1d3-5ffe-48c6-ac11-7a2e60802864"
}

output "connection_id" {
  value = data.neosync_connection.foo.id
}

output "connection_name" {
  value = data.neosync_connection.foo.name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The unique identifier of the connection

### Read-Only

- `name` (String) The unique name of the connection
