package provider

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	mgmtv1alpha1 "github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1"
	"github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1/mgmtv1alpha1connect"
	transformer_model "github.com/nucleuscloud/terraform-provider-neosync/internal/models/transformers"
)

var _ datasource.DataSource = &SystemTransformerDataSource{}

func NewSystemTransformerDataSource() datasource.DataSource {
	return &SystemTransformerDataSource{}
}

type SystemTransformerDataSource struct {
	client mgmtv1alpha1connect.TransformersServiceClient
}

type SystemTransformerDataSourceModel struct {
	Name        types.String                         `tfsdk:"name"`
	Description types.String                         `tfsdk:"description"`
	Datatype    types.Int64                          `tfsdk:"datatype"`
	Source      types.String                         `tfsdk:"source"`
	Config      *transformer_model.TransformerConfig `tfsdk:"config"`
}

func (d *SystemTransformerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_transformer"
}

func (d *SystemTransformerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	configSchema := transformerSchema
	configSchema.Description = "Default values for each system transformer. Can be used to feed into any values intended to remain unchanged for any user defined transformer"
	configSchema.Required = false
	configSchema.Computed = true

	resp.Schema = schema.Schema{
		Description: "Neosync System Transformer data source",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The unique name of the transformer",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the transformer",
				Computed:    true,
			},
			"datatype": schema.Int64Attribute{
				Description: "The datatype of the transformer",
				Computed:    true,
			},
			"source": schema.StringAttribute{
				Description: "The unique key that is used by the system to determine which transformer to use. This is the same value that is specified as the key in the config",
				Required:    true,
			},
			"config": configSchema,
		},
	}
}

func (d *SystemTransformerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*ConfigData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ConfigData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = providerData.TransformerClient
}

func (d *SystemTransformerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemTransformerDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	transResp, err := d.client.GetSystemTransformerBySource(ctx, connect.NewRequest(&mgmtv1alpha1.GetSystemTransformerBySourceRequest{
		Source: stateSourceToTransformerSource(data.Source.ValueString()),
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to get system transformer", err.Error())
		return
	}

	transformer := transResp.Msg.GetTransformer()
	modelConfig := &transformer_model.Transformer{}
	err = modelConfig.FromDto(transformer.Config)
	if err != nil {
		resp.Diagnostics.AddError("unable to convert dto transformer config to model", err.Error())
		return
	}

	data.Name = types.StringValue(transformer.Name)
	data.Description = types.StringValue(transformer.Description)
	data.Datatype = types.Int64Value(int64(transformer.DataType)) //nolint staticcheck
	data.Source = types.StringValue(transformerSourceToStateSource(transformer.Source))
	data.Config = modelConfig.Config
	tflog.Trace(ctx, "read system transformer", map[string]any{"source": data.Source.ValueString()})
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func stateSourceToTransformerSource(source string) mgmtv1alpha1.TransformerSource {
	if source == "null" {
		source = "generate_null"
	}
	if source == "custom" {
		source = "user_defined"
	}
	key := fmt.Sprintf("TRANSFORMER_SOURCE_%s", strings.ToUpper(source))

	value, ok := mgmtv1alpha1.TransformerSource_value[key]
	if !ok {
		return mgmtv1alpha1.TransformerSource_TRANSFORMER_SOURCE_UNSPECIFIED
	}
	return mgmtv1alpha1.TransformerSource(value)
}

func transformerSourceToStateSource(source mgmtv1alpha1.TransformerSource) string {
	name, ok := mgmtv1alpha1.TransformerSource_name[int32(source)]
	if !ok {
		return "unspecified"
	}
	return strings.ToLower(strings.TrimPrefix(name, "TRANSFORMER_SOURCE_"))
}
