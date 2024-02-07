# Pre-configure a system transformer with settings defined.
resource "neosync_user_defined_transformer" "my_transformer" {
  name        = "custom-gen-cc"
  description = "Generate CC with valid luhn turned off"
  datatype    = "int64"
  source      = "generate_card_number"
  config = {
    "generate_card_number" = {
      valid_luhn = false
    }
  }
}

# Utilize a system transformer to ensure values are properly set and exist
data "neosync_system_transformer" "generate_cc" {
  source = "generate_card_number"
}

resource "neosync_user_defined_transformer" "my_transformer" {
  name        = "custom-gen-cc"
  description = "Generate CC with valid luhn turned off"
  datatype    = data.neosync_system_transformer.generate_cc.datatype
  source      = data.neosync_system_transformer.generate_cc.source
  config = {
    [data.neosync_system_transformer.generate_cc.source] = {
      valid_luhn = false
    }
  }
}
