package provider

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	mgmtv1alpha1 "github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1"
	"github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1/mgmtv1alpha1connect"
	"github.com/nucleuscloud/terraform-provider-neosync/internal/models"
)

var _ resource.Resource = (*JobHookResource)(nil)
var _ resource.ResourceWithImportState = (*JobHookResource)(nil)

func NewJobHookResource() resource.Resource {
	return &JobHookResource{}
}

type JobHookResource struct {
	client mgmtv1alpha1connect.JobServiceClient
}

func (r *JobHookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jobhook"
}

func (r *JobHookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Job Hook",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the job hook.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The unique friendly name of the job hook",
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "What this hook does.",
			},
			"job_id": schema.StringAttribute{
				Description: "The unique identifier of the job this hook is associated with",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Whether or not this hook is enabled",
			},
			"priority": schema.Int32Attribute{
				Description: "The priority of this hook. 0-100, lower values are higher priority",
				Required:    true,
			},
			"config": schema.SingleNestedAttribute{
				Description: "The configuration for the hook itself",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"sql": schema.SingleNestedAttribute{
						Description: "A hook that will execute SQL on the specified connection",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"query": schema.StringAttribute{
								Description: "The SQL query that will be invoked",
								Required:    true,
							},
							"connection_id": schema.StringAttribute{
								Description: "The Neosync connection id that the query will be invoked on",
								Required:    true,
							},
							"timing": schema.SingleNestedAttribute{
								Description: "The timing of when in the run lifecycle the hook will be invoked",
								Required:    true,
								Attributes: map[string]schema.Attribute{
									"pre_sync": schema.SingleNestedAttribute{
										Description: "Will run before the first table sync (also truncation and schema init, if enabled)",
										Optional:    true,
										Attributes:  map[string]schema.Attribute{},
									},
									"post_sync": schema.SingleNestedAttribute{
										Description: "Will run after the last table is synced",
										Optional:    true,
										Attributes:  map[string]schema.Attribute{},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *JobHookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
}

func (r *JobHookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.JobHookResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jhResp, err := r.client.GetJobHook(ctx, connect.NewRequest(&mgmtv1alpha1.GetJobHookRequest{Id: data.Id.ValueString()}))
	if err != nil {
		resp.Diagnostics.AddError("unable to get job hook", err.Error())
		return
	}

	tflog.Trace(ctx, "retrieved job hook")

	hook := jhResp.Msg.GetHook()

	newModel := models.JobHookResourceModel{}
	err = newModel.FromDto(hook)
	if err != nil {
		resp.Diagnostics.AddError("read job hook model map error", err.Error())
		return
	}
	tflog.Trace(ctx, "mapped job hook to model during read")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newModel)...)
}

func (r *JobHookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.JobHookResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest, err := data.ToCreateJobHookDto()
	if err != nil {
		resp.Diagnostics.AddError("unable to create job request from planned state", err.Error())
		return
	}

	createResp, err := r.client.CreateJobHook(ctx, connect.NewRequest(createRequest))
	if err != nil {
		resp.Diagnostics.AddError("create job hook error", err.Error())
		return
	}

	hook := createResp.Msg.GetHook()
	tflog.Trace(ctx, "created job hook")

	newModel := models.JobHookResourceModel{}
	err = newModel.FromDto(hook)
	if err != nil {
		resp.Diagnostics.AddError("create job hook model map error", err.Error())
		return
	}
	tflog.Trace(ctx, "mapped job hook to model during creation")

	resp.Diagnostics.Append(resp.State.Set(ctx, &newModel)...)
}

func (r *JobHookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel models.JobHookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "read in planned model during update")

	var stateModel models.JobHookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "read in state model during update")

	updateRequest, err := planModel.ToUpdateJobHookDto()
	if err != nil {
		resp.Diagnostics.AddError("unable to map job hook model to update request", err.Error())
		return
	}

	updateResp, err := r.client.UpdateJobHook(ctx, connect.NewRequest(updateRequest))
	if err != nil {
		resp.Diagnostics.AddError("unable to update job hook", err.Error())
		return
	}
	tflog.Trace(ctx, "updated job hook")

	updatedHook := updateResp.Msg.GetHook()

	newModel := models.JobHookResourceModel{}
	err = newModel.FromDto(updatedHook)
	if err != nil {
		resp.Diagnostics.AddError("unable to update job hook resp to new state model", err.Error())
		return
	}
	tflog.Trace(ctx, "mapped job hook to new model during update")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newModel)...)
}

func (r *JobHookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.JobHookResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteJobHook(ctx, connect.NewRequest(&mgmtv1alpha1.DeleteJobHookRequest{Id: data.Id.ValueString()}))
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete job hook", err.Error())
		return
	}

	tflog.Trace(ctx, "deleted job hook")
}

func (r *JobHookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" {
		resp.Diagnostics.AddError("unable to import job hook", "must provide ID")
		return
	}

	jhResp, err := r.client.GetJobHook(ctx, connect.NewRequest(&mgmtv1alpha1.GetJobHookRequest{Id: req.ID}))
	if err != nil {
		resp.Diagnostics.AddError("unable to get job hook during import", err.Error())
		return
	}
	tflog.Trace(ctx, "retrieved job hook during import")

	hook := jhResp.Msg.GetHook()

	model := models.JobHookResourceModel{}
	err = model.FromDto(hook)
	if err != nil {
		resp.Diagnostics.AddError("unable to map job hook to model during import", err.Error())
		return
	}
	tflog.Trace(ctx, "mapped job hook to resource model during import")
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
