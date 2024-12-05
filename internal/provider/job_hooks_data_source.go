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

var _ datasource.DataSource = (*JobHookDataSource)(nil)

func NewJobHookDataSource() datasource.DataSource {
	return &JobHookDataSource{}
}

type JobHookDataSource struct {
	client mgmtv1alpha1connect.JobServiceClient
}

type JobHookDataSourceModel struct {
	Id    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	JobId types.String `tfsdk:"job_id"`
}

func (d *JobHookDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jobhook"
}

func (d *JobHookDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Neosync Job Hook data source",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The unique name of the job hook",
				Computed:    true,
			},
			"job_id": schema.StringAttribute{
				Description: "The unique identifier of the job this hook is associated with",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The unique identifier of the job hook",
				Required:    true,
			},
		},
	}
}

func (d *JobHookDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = providerData.JobClient
}

func (d *JobHookDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data JobHookDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jhResp, err := d.client.GetJobHook(ctx, connect.NewRequest(&mgmtv1alpha1.GetJobHookRequest{
		Id: data.Id.ValueString(),
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to get job hook by id", err.Error())
		return
	}

	jh := jhResp.Msg.GetHook()

	data.Name = types.StringValue(jh.GetName())
	data.JobId = types.StringValue(jh.GetJobId())

	tflog.Trace(ctx, "read job hook", map[string]any{"id": data.Id.ValueString()})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
