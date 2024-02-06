---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "neosync_job Resource - terraform-provider-neosync"
subcategory: ""
description: |-
  Job resource
---

# neosync_job (Resource)

Job resource



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `destinations` (Attributes List) (see [below for nested schema](#nestedatt--destinations))
- `name` (String) The unique friendly name of the job
- `source` (Attributes) (see [below for nested schema](#nestedatt--source))

### Optional

- `account_id` (String) The unique identifier of the account. Can be pulled from the API Key if present, or must be specified if using a user access token
- `cron_schedule` (String) A cron string for how often it's desired to schedule the job to run
- `mappings` (Attributes List) (see [below for nested schema](#nestedatt--mappings))
- `sync_options` (Attributes) (see [below for nested schema](#nestedatt--sync_options))
- `workflow_options` (Attributes) (see [below for nested schema](#nestedatt--workflow_options))

### Read-Only

- `id` (String) The unique identifier of the connection

<a id="nestedatt--destinations"></a>
### Nested Schema for `destinations`

Optional:

- `aws_s3` (Attributes) (see [below for nested schema](#nestedatt--destinations--aws_s3))
- `connection_id` (String)
- `id` (String)
- `mysql` (Attributes) (see [below for nested schema](#nestedatt--destinations--mysql))
- `postgres` (Attributes) (see [below for nested schema](#nestedatt--destinations--postgres))

<a id="nestedatt--destinations--aws_s3"></a>
### Nested Schema for `destinations.aws_s3`


<a id="nestedatt--destinations--mysql"></a>
### Nested Schema for `destinations.mysql`

Required:

- `init_table_schema` (Boolean)

Optional:

- `truncate_table` (Attributes) (see [below for nested schema](#nestedatt--destinations--mysql--truncate_table))

<a id="nestedatt--destinations--mysql--truncate_table"></a>
### Nested Schema for `destinations.mysql.truncate_table`

Optional:

- `truncate_before_insert` (Boolean)



<a id="nestedatt--destinations--postgres"></a>
### Nested Schema for `destinations.postgres`

Required:

- `init_table_schema` (Boolean)

Optional:

- `truncate_table` (Attributes) (see [below for nested schema](#nestedatt--destinations--postgres--truncate_table))

<a id="nestedatt--destinations--postgres--truncate_table"></a>
### Nested Schema for `destinations.postgres.truncate_table`

Optional:

- `cascade` (Boolean)
- `truncate_before_insert` (Boolean)




<a id="nestedatt--source"></a>
### Nested Schema for `source`

Optional:

- `aws_s3` (Attributes) (see [below for nested schema](#nestedatt--source--aws_s3))
- `generate` (Attributes) (see [below for nested schema](#nestedatt--source--generate))
- `mysql` (Attributes) (see [below for nested schema](#nestedatt--source--mysql))
- `postgres` (Attributes) (see [below for nested schema](#nestedatt--source--postgres))

<a id="nestedatt--source--aws_s3"></a>
### Nested Schema for `source.aws_s3`

Required:

- `connection_id` (String)


<a id="nestedatt--source--generate"></a>
### Nested Schema for `source.generate`

Required:

- `fk_source_connection_id` (String)
- `schemas` (Attributes List) (see [below for nested schema](#nestedatt--source--generate--schemas))

<a id="nestedatt--source--generate--schemas"></a>
### Nested Schema for `source.generate.schemas`

Required:

- `schema` (String)
- `tables` (Attributes List) (see [below for nested schema](#nestedatt--source--generate--schemas--tables))

<a id="nestedatt--source--generate--schemas--tables"></a>
### Nested Schema for `source.generate.schemas.tables`

Required:

- `table` (String)

Optional:

- `row_count` (Number)




<a id="nestedatt--source--mysql"></a>
### Nested Schema for `source.mysql`

Required:

- `connection_id` (String)
- `halt_on_new_column_addition` (Boolean)

Optional:

- `schemas` (Attributes List) (see [below for nested schema](#nestedatt--source--mysql--schemas))

<a id="nestedatt--source--mysql--schemas"></a>
### Nested Schema for `source.mysql.schemas`

Required:

- `schema` (String)
- `tables` (Attributes List) (see [below for nested schema](#nestedatt--source--mysql--schemas--tables))

<a id="nestedatt--source--mysql--schemas--tables"></a>
### Nested Schema for `source.mysql.schemas.tables`

Required:

- `table` (String)

Optional:

- `where_clause` (String)




<a id="nestedatt--source--postgres"></a>
### Nested Schema for `source.postgres`

Required:

- `connection_id` (String)
- `halt_on_new_column_addition` (Boolean)

Optional:

- `schemas` (Attributes List) (see [below for nested schema](#nestedatt--source--postgres--schemas))

<a id="nestedatt--source--postgres--schemas"></a>
### Nested Schema for `source.postgres.schemas`

Required:

- `schema` (String)
- `tables` (Attributes List) (see [below for nested schema](#nestedatt--source--postgres--schemas--tables))

<a id="nestedatt--source--postgres--schemas--tables"></a>
### Nested Schema for `source.postgres.schemas.tables`

Required:

- `table` (String)

Optional:

- `where_clause` (String)





<a id="nestedatt--mappings"></a>
### Nested Schema for `mappings`

Required:

- `column` (String)
- `schema` (String)
- `table` (String)
- `transformer` (Attributes) (see [below for nested schema](#nestedatt--mappings--transformer))

<a id="nestedatt--mappings--transformer"></a>
### Nested Schema for `mappings.transformer`

Required:

- `config` (Attributes) (see [below for nested schema](#nestedatt--mappings--transformer--config))
- `source` (String)

<a id="nestedatt--mappings--transformer--config"></a>
### Nested Schema for `mappings.transformer.config`

Optional:

- `passthrough` (Attributes) (see [below for nested schema](#nestedatt--mappings--transformer--config--passthrough))

<a id="nestedatt--mappings--transformer--config--passthrough"></a>
### Nested Schema for `mappings.transformer.config.passthrough`





<a id="nestedatt--sync_options"></a>
### Nested Schema for `sync_options`

Optional:

- `retry_policy` (Attributes) (see [below for nested schema](#nestedatt--sync_options--retry_policy))
- `schedule_to_close_timeout` (Number)
- `start_to_close_timeout` (Number)

<a id="nestedatt--sync_options--retry_policy"></a>
### Nested Schema for `sync_options.retry_policy`

Optional:

- `maximum_attempts` (Number)



<a id="nestedatt--workflow_options"></a>
### Nested Schema for `workflow_options`

Optional:

- `run_timeout` (Number)