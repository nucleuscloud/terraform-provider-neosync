package provider

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	mgmtv1alpha1 "github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1"
	"github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1/mgmtv1alpha1connect"
)

var _ datasource.DataSource = &UserDefinedTransformerDataSource{}

func NewUserDefinedTransformerDataSource() datasource.DataSource {
	return &UserDefinedTransformerDataSource{}
}

type UserDefinedTransformerDataSource struct {
	client mgmtv1alpha1connect.TransformersServiceClient
}

type UserDefinedTransformerDataSourceModel struct {
	Name types.String `tfsdk:"name"`
	Id   types.String `tfsdk:"id"`
}

func (d *UserDefinedTransformerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_defined_transformer"
}

func (d *UserDefinedTransformerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Neosync User Defined Transformer data source",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The unique name of the transformer",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The unique identifier of the transformer",
				Required:    true,
			},
		},
	}
}

func (d *UserDefinedTransformerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserDefinedTransformerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserDefinedTransformerDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	transResp, err := d.client.GetUserDefinedTransformerById(ctx, connect.NewRequest(&mgmtv1alpha1.GetUserDefinedTransformerByIdRequest{
		TransformerId: data.Id.ValueString(),
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to get transformer by id", err.Error())
		return
	}

	data.Name = types.StringValue(transResp.Msg.Transformer.Name)
	tflog.Trace(ctx, "read transformer", map[string]any{"id": data.Id.ValueString()})
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
