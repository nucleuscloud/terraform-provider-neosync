resource "neosync_job" "prod_to_stage" {
  name = "prod-to-stage"

  source = {
    postgres = {
      halt_on_new_column_additon = false
      connection_id              = var.prod_connection_id
    }
  }
  destinations = [
    {
      connection_id = var.stage_connection_id
      postgres = {
        init_table_schema = false
        truncate_table = {
          truncate_before_insert = true
          cascade                = true
        }
      }
    }
  ]

  mappings = [
    {
      schema = "public"
      table  = "users"
      column = "id"
      transformer = {
        source = "passthrough"
        config = {
          passthrough = {}
        }
      }
    }
  ]
}
