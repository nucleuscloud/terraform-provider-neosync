package provider

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

var (
	transformerSchema = schema.SingleNestedAttribute{
		Description: "",
		Required:    true,
		Attributes: map[string]schema.Attribute{
			"generate_email": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"transform_email": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"preserve_domain": schema.BoolAttribute{
						Description: "",
						Required:    true,
					},
					"preserve_length": schema.BoolAttribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"generate_bool": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"generate_card_number": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"valid_luhn": schema.BoolAttribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"generate_city": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"generate_e164_phone_number": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"min": schema.Int64Attribute{
						Description: "",
						Required:    true,
					},
					"max": schema.Int64Attribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"generate_firstname": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"generate_float64": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"randomize_sign": schema.BoolAttribute{
						Description: "",
						Required:    true,
					},
					"min": schema.Float64Attribute{
						Description: "",
						Required:    true,
					},
					"max": schema.Float64Attribute{
						Description: "",
						Required:    true,
					},
					"precision": schema.Int64Attribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"generate_full_address": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"generate_fullname": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"generate_gender": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"abbreviate": schema.BoolAttribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"generate_int64_phone_number": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"generate_int64": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"randomize_sign": schema.BoolAttribute{
						Description: "",
						Required:    true,
					},
					"min": schema.Int64Attribute{
						Description: "",
						Required:    true,
					},
					"max": schema.Int64Attribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"generate_lastname": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"generate_sha256": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"generate_ssn": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"generate_state": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"generate_street_address": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"generate_string_phone_number": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"min": schema.Int64Attribute{
						Description: "",
						Required:    true,
					},
					"max": schema.Int64Attribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"generate_string": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"min": schema.Int64Attribute{
						Description: "",
						Required:    true,
					},
					"max": schema.Int64Attribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"generate_unix_timestamp": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"generate_username": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"generate_utc_timestamp": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"generate_uuid": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"include_hyphens": schema.BoolAttribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"generate_zipcode": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"transform_e164_phone_number": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"preserve_length": schema.BoolAttribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"transform_firstname": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"preserve_length": schema.BoolAttribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"transform_float64": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"randomization_range_min": schema.Float64Attribute{
						Description: "",
						Required:    true,
					},
					"randomization_range_max": schema.Float64Attribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"transform_fullname": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"preserve_length": schema.BoolAttribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"transform_int64_phone_number": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"preserve_length": schema.BoolAttribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"transform_int64": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"randomization_range_min": schema.Int64Attribute{
						Description: "",
						Required:    true,
					},
					"randomization_range_max": schema.Int64Attribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"transform_lastname": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"preserve_length": schema.BoolAttribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"transform_phone_number": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"preserve_length": schema.BoolAttribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"transform_string": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"preserve_length": schema.BoolAttribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"passthrough": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"null": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"user_defined_transformer": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"generate_default": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
			"transform_javascript": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"code": schema.StringAttribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"generate_categorical": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"categories": schema.StringAttribute{
						Description: "",
						Required:    true,
					},
				},
			},
			"transform_character_scramble": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
			},
		},
	}
)
