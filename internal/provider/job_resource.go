package provider

import (
	"context"
	"errors"
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
	job_model "github.com/nucleuscloud/terraform-provider-neosync/internal/models/jobs"
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

func (r *JobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_job"
}

func (r *JobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Job resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The unique friendly name of the job",
				Required:    true,
			},
			"account_id": schema.StringAttribute{
				Description:   "The unique identifier of the account. Can be pulled from the API Key if present, or must be specified if using a user access token",
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplaceIfConfigured()},
			},

			"source": schema.SingleNestedAttribute{
				Description: "Configuration details about the source data connection",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"postgres": schema.SingleNestedAttribute{
						Description: "Postgres specific connection configurations",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"halt_on_new_column_addition": schema.BoolAttribute{
								Description: "(Deprecated) Whether or not to halt the job if it detects a new column that has been added in the source that has not been defined in the mappings schema",
								Optional:    true,
							},
							"new_column_addition_strategy": schema.SingleNestedAttribute{
								Description: "Strategy for handling new column additions",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"halt_job": schema.SingleNestedAttribute{
										Description: "Halt job if a new column is detected",
										Optional:    true,
										Attributes:  map[string]schema.Attribute{},
									},
									"auto_map": schema.SingleNestedAttribute{
										Description: "Automatically handle unmapped columns using DB defaults or generators",
										Optional:    true,
										Attributes:  map[string]schema.Attribute{},
									},
								},
							},
							"connection_id": schema.StringAttribute{
								Description: "The unique identifier of the connection that is to be used as the source",
								Required:    true,
							},
							"schemas": schema.ListNestedAttribute{
								Description: "A list of schemas and table specific options",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"schema": schema.StringAttribute{
											Description: "The name of the schema",
											Required:    true,
										},
										"tables": schema.ListNestedAttribute{
											Description: "A list of tables and their specific options within the defined schema",
											Required:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"table": schema.StringAttribute{
														Description: "The name of the table",
														Required:    true,
													},
													"where_clause": schema.StringAttribute{
														Description: "A where clause that will be used to subset the table during sync",
														Optional:    true,
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"mysql": schema.SingleNestedAttribute{
						Description: "Mysql specific connection configurations",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"halt_on_new_column_addition": schema.BoolAttribute{
								Description: "Whether or not to halt the job if it detects a new column that has been added in the source that has not been defined in the mappings schema",
								Required:    true,
							},
							"connection_id": schema.StringAttribute{
								Description: "The unique identifier of the connection that is to be used as the source",
								Required:    true,
							},
							"schemas": schema.ListNestedAttribute{
								Description: "A list of schemas and table specific options",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"schema": schema.StringAttribute{
											Description: "The name of the schema",
											Required:    true,
										},
										"tables": schema.ListNestedAttribute{
											Description: "A list of tables and their specific options within the defined schema",
											Required:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"table": schema.StringAttribute{
														Description: "The name of the table",
														Required:    true,
													},
													"where_clause": schema.StringAttribute{
														Description: "A where clause that will be used to subset the table during sync",
														Optional:    true,
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"aws_s3": schema.SingleNestedAttribute{
						Description: "AWS S3 specific connection configurations",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"connection_id": schema.StringAttribute{
								Description: "The unique identifier of the connection that is to be used as the source",
								Required:    true,
							},
						},
					},
					"generate": schema.SingleNestedAttribute{
						Description: "Generate specific connection configurations. Currently only supports single table generation",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"fk_source_connection_id": schema.StringAttribute{
								Description: "The unique connection identifier that is used to generate schema specific details. This is usually set to the destination connectio id if it has been upserted with the schema already",
								Required:    true,
							},
							"schemas": schema.ListNestedAttribute{
								Description: "A list of schemas and table specific options",
								Required:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"schema": schema.StringAttribute{
											Description: "The name of the schema",
											Required:    true,
										},
										"tables": schema.ListNestedAttribute{
											Description: "A list of tables and their specific options within the defined schema",
											Required:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"table": schema.StringAttribute{
														Description: "The name of the table",
														Required:    true,
													},
													"row_count": schema.Int64Attribute{
														Description: "The number of rows to generate into the table",
														Optional:    true,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"destinations": schema.ListNestedAttribute{
				Description: "A list of destination connections and any relevant configurations that are available to them dependent on type",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description:   "The unique identifier of the destination resource. This is set after creation",
							Optional:      true,
							Computed:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"connection_id": schema.StringAttribute{
							Description: "The unique identifier of the connection that will be used during the synchronization process",
							Optional:    true, // required if id is not set
						},
						"postgres": schema.SingleNestedAttribute{
							Description: "Postgres connection specific options",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"truncate_table": schema.SingleNestedAttribute{
									Description: "Details about what truncation should occur",
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										"truncate_before_insert": schema.BoolAttribute{
											Description: "Will truncate the table prior to insertion of any records",
											Optional:    true,
										},
										"cascade": schema.BoolAttribute{
											Description: "Will truncate cascade. This is required if the table holds any foreign key constraints. If this is true, truncate_before_insert must also be true",
											Optional:    true,
										},
									},
								},
								"init_table_schema": schema.BoolAttribute{
									Description: "Whether or not to have neosync init the table schema and constraints it pulled from the source connection",
									Required:    true,
								},
							},
						},
						"mysql": schema.SingleNestedAttribute{
							Description: "Mysql connection specific options",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"truncate_table": schema.SingleNestedAttribute{
									Description: "Details about what truncation should occur",
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										"truncate_before_insert": schema.BoolAttribute{
											Description: "Will truncate the table prior to insertion of any records",
											Optional:    true,
										},
									},
								},
								"init_table_schema": schema.BoolAttribute{
									Description: "Whether or not to have neosync init the table schema and constraints it pulled from the source connection",
									Required:    true,
								},
							},
						},
						"aws_s3": schema.SingleNestedAttribute{
							Description: "AWS S3 connection specific options",
							Optional:    true,
							Attributes:  map[string]schema.Attribute{},
						},
					},
				},
			},
			"mappings": schema.ListNestedAttribute{
				Description: "Details each schema,table,column along with the transformation that will be executed",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"schema": schema.StringAttribute{
							Description: "The database schema",
							Required:    true,
						},
						"table": schema.StringAttribute{
							Description: "The database table",
							Required:    true,
						},
						"column": schema.StringAttribute{
							Description: "The column in the specified table",
							Required:    true,
						},
						"transformer": schema.SingleNestedAttribute{
							Description: "The transformer that will be performed on the column",
							Required:    true,
							Attributes: map[string]schema.Attribute{
								"config": transformerSchema,
							},
						},
					},
				},
			},

			"cron_schedule": schema.StringAttribute{
				Description: "A cron string for how often it's desired to schedule the job to run",
				Optional:    true,
				Computed:    true,
			},

			"sync_options": schema.SingleNestedAttribute{
				Description: "Advanced settings and other options specific to a table sync",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"schedule_to_close_timeout": schema.Int64Attribute{
						Description: "The maximum amount of time allotted for a table sync with retries",
						Optional:    true,
					},
					"start_to_close_timeout": schema.Int64Attribute{
						Description: "The amount of time allotted for a table sync",
						Optional:    true,
					},
					"retry_policy": schema.SingleNestedAttribute{
						Description: "The table sync retry policy",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"maximum_attempts": schema.Int64Attribute{
								Description: "The maximum number of times to retry if there is a failure or timeout",
								Optional:    true,
							},
						},
					},
				},
			},
			"workflow_options": schema.SingleNestedAttribute{
				Description: "Advanced settings and other options specific to a job run",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"run_timeout": schema.Int64Attribute{
						Description: "The max amount of time a job run is allotted",
						Optional:    true,
					},
				},
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
	var data job_model.JobResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountId, err := r.getAccountId(&data)
	if err != nil {
		resp.Diagnostics.AddError("no account id", err.Error())
		return
	}
	data.AccountId = types.StringValue(accountId)

	jobRequest, err := data.ToCreateJobDto()
	if err != nil {
		resp.Diagnostics.AddError("unable to create job request", err.Error())
		return
	}

	jobResp, err := r.client.CreateJob(ctx, connect.NewRequest(jobRequest))
	if err != nil {
		resp.Diagnostics.AddError("create job error", err.Error())
		return
	}

	job := jobResp.Msg.GetJob()
	tflog.Trace(ctx, "created job")

	newModel := job_model.JobResourceModel{}
	err = newModel.FromDto(job)
	if err != nil {
		resp.Diagnostics.AddError("job translate error", err.Error())
		return
	}

	tflog.Trace(ctx, "mapped job to model during creation")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newModel)...)
}

func (r *JobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data job_model.JobResourceModel

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
	tflog.Trace(ctx, "got job")

	job := jobResp.Msg.GetJob()

	newModel := job_model.JobResourceModel{}
	err = newModel.FromDto(job)
	if err != nil {
		resp.Diagnostics.AddError("unable to convert dto to state", err.Error())
		return
	}

	tflog.Trace(ctx, "mapped job to model")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newModel)...)
}

func (r *JobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel job_model.JobResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateModel job_model.JobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateJobRequest, err := planModel.ToUpdateJobDto(&stateModel, planModel.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to create update job request", err.Error())
		return
	}

	if updateJobRequest.UpdateJobScheduleRequest != nil {
		_, err = r.client.UpdateJobSchedule(ctx, connect.NewRequest(updateJobRequest.UpdateJobScheduleRequest))
		if err != nil {
			resp.Diagnostics.AddError("unable to update job schedule", err.Error())
			return
		}
	}

	if updateJobRequest.UpdateJobSourceConnectionRequest != nil {
		_, err = r.client.UpdateJobSourceConnection(ctx, connect.NewRequest(updateJobRequest.UpdateJobSourceConnectionRequest))
		if err != nil {
			resp.Diagnostics.AddError("unable to update job source connection", err.Error())
			return
		}
	}

	if updateJobRequest.CreateJobDestinationConnectionsRequest != nil {
		_, err = r.client.CreateJobDestinationConnections(ctx, connect.NewRequest(updateJobRequest.CreateJobDestinationConnectionsRequest))
		if err != nil {
			resp.Diagnostics.AddError("unable to create job destination connections", err.Error())
			return
		}
	}

	if len(updateJobRequest.UpdateJobDestinationConnectionRequests) > 0 {
		for _, req := range updateJobRequest.UpdateJobDestinationConnectionRequests {
			_, err = r.client.UpdateJobDestinationConnection(ctx, connect.NewRequest(req))
			if err != nil {
				resp.Diagnostics.AddError("unable to update job destination connection", err.Error())
				return
			}
		}
	}

	if len(updateJobRequest.DeleteJobDestinationConnectionRequests) > 0 {
		for _, req := range updateJobRequest.DeleteJobDestinationConnectionRequests {
			_, err = r.client.DeleteJobDestinationConnection(ctx, connect.NewRequest(req))
			if err != nil {
				resp.Diagnostics.AddError("unable to delete job destination connection", err.Error())
				return
			}
		}
	}

	getResp, err := r.client.GetJob(ctx, connect.NewRequest(&mgmtv1alpha1.GetJobRequest{
		Id: planModel.Id.ValueString(),
	}))
	if err != nil {
		resp.Diagnostics.AddError("unable to get updated job model", err.Error())
		return
	}
	tflog.Trace(ctx, "got updated job model")
	job := getResp.Msg.GetJob()

	updatedModel := job_model.JobResourceModel{}
	err = updatedModel.FromDto(job)
	if err != nil {
		resp.Diagnostics.AddError("unable to model latest job dto", err.Error())
		return
	}

	tflog.Trace(ctx, "updated job")
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedModel)...)
}

func (r *JobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data job_model.JobResourceModel

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
	tflog.Trace(ctx, "retrieved job during import")
	job := jobResp.Msg.GetJob()

	var data job_model.JobResourceModel
	err = data.FromDto(job)
	if err != nil {
		resp.Diagnostics.AddError("unable to map job to model", err.Error())
		return
	}

	tflog.Trace(ctx, "mapped job to model during import")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobResource) getAccountId(data *job_model.JobResourceModel) (string, error) {
	var accountId string
	if data.AccountId.ValueString() == "" {
		if r.accountId != nil {
			accountId = *r.accountId
		}
	} else {
		accountId = data.AccountId.ValueString()
	}
	if accountId == "" {
		return "", errors.New("must provide account id either on the resource or provide through environment configuration")
	}
	return accountId, nil
}
