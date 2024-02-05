package provider

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	mgmtv1alpha1 "github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1"
	"github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1/mgmtv1alpha1connect"
)

var _ resource.Resource = &JobResource{}
var _ resource.ResourceWithImportState = &JobResource{}

func NewJobResource() resource.Resource {
	return &JobResource{}
}

type JobResource struct {
	client    mgmtv1alpha1connect.JobServiceClient
	accountId *string
}

type JobResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	AccountId types.String `tfsdk:"account_id"`
}

func (r *JobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_job"
}

func (r *JobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Example resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The unique friendly name of the connection",
				Required:    true,
			},
			"account_id": schema.StringAttribute{
				Description:   "The unique identifier of the account. Can be pulled from the API Key if present, or must be specified if using a user access token",
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},

			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the connection",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
		},
	}
}

func (r *JobResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*ConfigData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ConfigData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = providerData.JobClient
	r.accountId = providerData.AccountId
}

func (r *JobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data JobResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var accountId string
	if data.AccountId.ValueString() == "" {
		if r.accountId != nil {
			accountId = *r.accountId
		}
	} else {
		accountId = data.AccountId.ValueString()
	}
	if accountId == "" {
		resp.Diagnostics.AddError("no account id", "must provide account id either on the resource or provide through environment configuration")
		return
	}

	jobResp, err := r.client.CreateJob(ctx, connect.NewRequest(&mgmtv1alpha1.CreateJobRequest{
		JobName:   data.Name.ValueString(),
		AccountId: accountId,
	}))
	if err != nil {
		resp.Diagnostics.AddError("create job error", err.Error())
		return
	}

	job := jobResp.Msg.Job

	data.Id = types.StringValue(job.Id)
	data.Name = types.StringValue(job.Name)
	data.AccountId = types.StringValue(job.AccountId)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created job resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data JobResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	jobResp, err := r.client.GetJob(ctx, connect.NewRequest(&mgmtv1alpha1.GetJobRequest{
		Id: data.Id.ValueString(),
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to get job", err.Error())
		return
	}

	job := jobResp.Msg.Job

	data.Id = types.StringValue(job.Id)
	data.Name = types.StringValue(job.Name)
	data.AccountId = types.StringValue(job.AccountId)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data JobResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// connResp, err := r.client.UpdateConnection(ctx, connect.NewRequest(&mgmtv1alpha1.UpdateConnectionRequest{
	// 	Id:               data.Id.ValueString(),
	// 	Name:             data.Name.ValueString(),
	// 	ConnectionConfig: cc,
	// }))
	// if err != nil {
	// 	resp.Diagnostics.AddError("Unable to update connection", err.Error())
	// 	return
	// }

	// connection := connResp.Msg.Connection

	// data.Id = types.StringValue(connection.Id)
	// data.Name = types.StringValue(connection.Name)
	// data.AccountId = types.StringValue(connection.AccountId)
	// err = hydrateResourceModelFromConnectionConfig(connection.ConnectionConfig, &data)
	// if err != nil {
	// 	resp.Diagnostics.AddError("connection config hydration error", err.Error())
	// 	return
	// }

	tflog.Trace(ctx, "updated job")
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data JobResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteJob(ctx, connect.NewRequest(&mgmtv1alpha1.DeleteJobRequest{
		Id: data.Id.ValueString(),
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete job", err.Error())
		return
	}

	tflog.Trace(ctx, "deleted job")
}

func (r *JobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError("Unable to import", "must provide ID")
		return
	}

	jobResp, err := r.client.GetJob(ctx, connect.NewRequest(&mgmtv1alpha1.GetJobRequest{
		Id: req.ID,
	}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to get job", err.Error())
		return
	}

	job := jobResp.Msg.Job

	var data JobResourceModel
	data.Id = types.StringValue(job.Id)
	data.Name = types.StringValue(job.Name)
	data.AccountId = types.StringValue(job.AccountId)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
