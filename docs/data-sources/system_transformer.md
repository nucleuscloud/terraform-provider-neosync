---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "neosync_system_transformer Data Source - terraform-provider-neosync"
subcategory: ""
description: |-
  Neosync System Transformer data source
---

# neosync_system_transformer (Data Source)

Neosync System Transformer data source



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `source` (String) The unique key that is used by the system to determine which transformer to use. This is the same value that is specified as the key in the config

### Read-Only

- `config` (Attributes) Default values for each system transformer. Can be used to feed into any values intended to remain unchanged for any user defined transformer (see [below for nested schema](#nestedatt--config))
- `datatype` (String) The datatype of the transformer
- `description` (String) The description of the transformer
- `name` (String) The unique name of the transformer

<a id="nestedatt--config"></a>
### Nested Schema for `config`

Optional:

- `generate_bool` (Attributes) (see [below for nested schema](#nestedatt--config--generate_bool))
- `generate_card_number` (Attributes) (see [below for nested schema](#nestedatt--config--generate_card_number))
- `generate_categorical` (Attributes) (see [below for nested schema](#nestedatt--config--generate_categorical))
- `generate_city` (Attributes) (see [below for nested schema](#nestedatt--config--generate_city))
- `generate_default` (Attributes) (see [below for nested schema](#nestedatt--config--generate_default))
- `generate_e164_phone_number` (Attributes) (see [below for nested schema](#nestedatt--config--generate_e164_phone_number))
- `generate_email` (Attributes) (see [below for nested schema](#nestedatt--config--generate_email))
- `generate_firstname` (Attributes) (see [below for nested schema](#nestedatt--config--generate_firstname))
- `generate_float64` (Attributes) (see [below for nested schema](#nestedatt--config--generate_float64))
- `generate_full_address` (Attributes) (see [below for nested schema](#nestedatt--config--generate_full_address))
- `generate_fullname` (Attributes) (see [below for nested schema](#nestedatt--config--generate_fullname))
- `generate_gender` (Attributes) (see [below for nested schema](#nestedatt--config--generate_gender))
- `generate_int64` (Attributes) (see [below for nested schema](#nestedatt--config--generate_int64))
- `generate_int64_phone_number` (Attributes) (see [below for nested schema](#nestedatt--config--generate_int64_phone_number))
- `generate_lastname` (Attributes) (see [below for nested schema](#nestedatt--config--generate_lastname))
- `generate_sha256` (Attributes) (see [below for nested schema](#nestedatt--config--generate_sha256))
- `generate_ssn` (Attributes) (see [below for nested schema](#nestedatt--config--generate_ssn))
- `generate_state` (Attributes) (see [below for nested schema](#nestedatt--config--generate_state))
- `generate_street_address` (Attributes) (see [below for nested schema](#nestedatt--config--generate_street_address))
- `generate_string` (Attributes) (see [below for nested schema](#nestedatt--config--generate_string))
- `generate_string_phone_number` (Attributes) (see [below for nested schema](#nestedatt--config--generate_string_phone_number))
- `generate_unix_timestamp` (Attributes) (see [below for nested schema](#nestedatt--config--generate_unix_timestamp))
- `generate_username` (Attributes) (see [below for nested schema](#nestedatt--config--generate_username))
- `generate_utc_timestamp` (Attributes) (see [below for nested schema](#nestedatt--config--generate_utc_timestamp))
- `generate_uuid` (Attributes) (see [below for nested schema](#nestedatt--config--generate_uuid))
- `generate_zipcode` (Attributes) (see [below for nested schema](#nestedatt--config--generate_zipcode))
- `null` (Attributes) (see [below for nested schema](#nestedatt--config--null))
- `passthrough` (Attributes) (see [below for nested schema](#nestedatt--config--passthrough))
- `transform_character_scramble` (Attributes) (see [below for nested schema](#nestedatt--config--transform_character_scramble))
- `transform_e164_phone_number` (Attributes) (see [below for nested schema](#nestedatt--config--transform_e164_phone_number))
- `transform_email` (Attributes) (see [below for nested schema](#nestedatt--config--transform_email))
- `transform_firstname` (Attributes) (see [below for nested schema](#nestedatt--config--transform_firstname))
- `transform_float64` (Attributes) (see [below for nested schema](#nestedatt--config--transform_float64))
- `transform_fullname` (Attributes) (see [below for nested schema](#nestedatt--config--transform_fullname))
- `transform_int64` (Attributes) (see [below for nested schema](#nestedatt--config--transform_int64))
- `transform_int64_phone_number` (Attributes) (see [below for nested schema](#nestedatt--config--transform_int64_phone_number))
- `transform_javascript` (Attributes) (see [below for nested schema](#nestedatt--config--transform_javascript))
- `transform_lastname` (Attributes) (see [below for nested schema](#nestedatt--config--transform_lastname))
- `transform_phone_number` (Attributes) (see [below for nested schema](#nestedatt--config--transform_phone_number))
- `transform_string` (Attributes) (see [below for nested schema](#nestedatt--config--transform_string))
- `user_defined_transformer` (Attributes) (see [below for nested schema](#nestedatt--config--user_defined_transformer))

<a id="nestedatt--config--generate_bool"></a>
### Nested Schema for `config.generate_bool`


<a id="nestedatt--config--generate_card_number"></a>
### Nested Schema for `config.generate_card_number`

Required:

- `valid_luhn` (Boolean)


<a id="nestedatt--config--generate_categorical"></a>
### Nested Schema for `config.generate_categorical`

Required:

- `categories` (String)


<a id="nestedatt--config--generate_city"></a>
### Nested Schema for `config.generate_city`


<a id="nestedatt--config--generate_default"></a>
### Nested Schema for `config.generate_default`


<a id="nestedatt--config--generate_e164_phone_number"></a>
### Nested Schema for `config.generate_e164_phone_number`

Required:

- `max` (Number)
- `min` (Number)


<a id="nestedatt--config--generate_email"></a>
### Nested Schema for `config.generate_email`


<a id="nestedatt--config--generate_firstname"></a>
### Nested Schema for `config.generate_firstname`


<a id="nestedatt--config--generate_float64"></a>
### Nested Schema for `config.generate_float64`

Required:

- `max` (Number)
- `min` (Number)
- `precision` (Number)
- `randomize_sign` (Boolean)


<a id="nestedatt--config--generate_full_address"></a>
### Nested Schema for `config.generate_full_address`


<a id="nestedatt--config--generate_fullname"></a>
### Nested Schema for `config.generate_fullname`


<a id="nestedatt--config--generate_gender"></a>
### Nested Schema for `config.generate_gender`

Required:

- `abbreviate` (Boolean)


<a id="nestedatt--config--generate_int64"></a>
### Nested Schema for `config.generate_int64`

Required:

- `max` (Number)
- `min` (Number)
- `randomize_sign` (Boolean)


<a id="nestedatt--config--generate_int64_phone_number"></a>
### Nested Schema for `config.generate_int64_phone_number`


<a id="nestedatt--config--generate_lastname"></a>
### Nested Schema for `config.generate_lastname`


<a id="nestedatt--config--generate_sha256"></a>
### Nested Schema for `config.generate_sha256`


<a id="nestedatt--config--generate_ssn"></a>
### Nested Schema for `config.generate_ssn`


<a id="nestedatt--config--generate_state"></a>
### Nested Schema for `config.generate_state`


<a id="nestedatt--config--generate_street_address"></a>
### Nested Schema for `config.generate_street_address`


<a id="nestedatt--config--generate_string"></a>
### Nested Schema for `config.generate_string`

Required:

- `max` (Number)
- `min` (Number)


<a id="nestedatt--config--generate_string_phone_number"></a>
### Nested Schema for `config.generate_string_phone_number`

Required:

- `max` (Number)
- `min` (Number)


<a id="nestedatt--config--generate_unix_timestamp"></a>
### Nested Schema for `config.generate_unix_timestamp`


<a id="nestedatt--config--generate_username"></a>
### Nested Schema for `config.generate_username`


<a id="nestedatt--config--generate_utc_timestamp"></a>
### Nested Schema for `config.generate_utc_timestamp`


<a id="nestedatt--config--generate_uuid"></a>
### Nested Schema for `config.generate_uuid`

Required:

- `include_hyphens` (Boolean)


<a id="nestedatt--config--generate_zipcode"></a>
### Nested Schema for `config.generate_zipcode`


<a id="nestedatt--config--null"></a>
### Nested Schema for `config.null`


<a id="nestedatt--config--passthrough"></a>
### Nested Schema for `config.passthrough`


<a id="nestedatt--config--transform_character_scramble"></a>
### Nested Schema for `config.transform_character_scramble`


<a id="nestedatt--config--transform_e164_phone_number"></a>
### Nested Schema for `config.transform_e164_phone_number`

Required:

- `preserve_length` (Boolean)


<a id="nestedatt--config--transform_email"></a>
### Nested Schema for `config.transform_email`

Required:

- `preserve_domain` (Boolean)
- `preserve_length` (Boolean)


<a id="nestedatt--config--transform_firstname"></a>
### Nested Schema for `config.transform_firstname`

Required:

- `preserve_length` (Boolean)


<a id="nestedatt--config--transform_float64"></a>
### Nested Schema for `config.transform_float64`

Required:

- `randomization_range_max` (Number)
- `randomization_range_min` (Number)


<a id="nestedatt--config--transform_fullname"></a>
### Nested Schema for `config.transform_fullname`

Required:

- `preserve_length` (Boolean)


<a id="nestedatt--config--transform_int64"></a>
### Nested Schema for `config.transform_int64`

Required:

- `randomization_range_max` (Number)
- `randomization_range_min` (Number)


<a id="nestedatt--config--transform_int64_phone_number"></a>
### Nested Schema for `config.transform_int64_phone_number`

Required:

- `preserve_length` (Boolean)


<a id="nestedatt--config--transform_javascript"></a>
### Nested Schema for `config.transform_javascript`

Required:

- `code` (String)


<a id="nestedatt--config--transform_lastname"></a>
### Nested Schema for `config.transform_lastname`

Required:

- `preserve_length` (Boolean)


<a id="nestedatt--config--transform_phone_number"></a>
### Nested Schema for `config.transform_phone_number`

Required:

- `preserve_length` (Boolean)


<a id="nestedatt--config--transform_string"></a>
### Nested Schema for `config.transform_string`

Required:

- `preserve_length` (Boolean)


<a id="nestedatt--config--user_defined_transformer"></a>
### Nested Schema for `config.user_defined_transformer`

Required:

- `id` (String)