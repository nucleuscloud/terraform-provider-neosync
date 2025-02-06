resource "neosync_job" "prod_to_stage" {
  name = "prod-to-stage"

  source = {
    postgres = {
      new_column_addition_strategy = {
        halt_job = {}
      }
      connection_id = var.prod_connection_id
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
        config = {
          passthrough = {}
        }
      }
    },
    {
      schema = "public"
      table  = "accounts"
      column = "user_id"
      transformer = {
        config = {
          passthrough = {}
        }
      }
    }
  ]
  virtual_foreign_keys = [
    {
      schema  = "public"
      table   = "accounts"
      columns = ["user_id"]
      foreign_key = {
        schema  = "public"
        table   = "users"
        columns = ["id"]
      }
    }
  ]
}
