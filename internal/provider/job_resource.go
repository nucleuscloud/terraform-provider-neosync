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
	Id              types.String     `tfsdk:"id"`
	Name            types.String     `tfsdk:"name"`
	AccountId       types.String     `tfsdk:"account_id"`
	JobSource       *JobSource       `tfsdk:"source"`
	Destinations    []JobDestination `tfsdk:"destinations"`
	Mappings        []JobMapping     `tfsdk:"mappings"`
	CronSchedule    types.String     `tfsdk:"cron_schedule"`
	SyncOptions     *ActivityOptions `tfsdk:"sync_options"`
	WorkflowOptions *WorkflowOptions `tfsdk:"workflow_options"`
}

type JobSource struct {
	Postgres *JobSourcePostgresOptions `tfsdk:"postgres"`
	Mysql    *JobSourceMysqlOptions    `tfsdk:"mysql"`
	Generate *JobSourceGenerateOptions `tfsdk:"generate"`
	AwsS3    *JobSourceAwsS3Options    `tfsdk:"aws_s3"`
}
type JobSourcePostgresOptions struct {
	HaltOnNewColumnAddition types.Bool                            `tfsdk:"halt_on_new_column_addition"`
	ConnectionId            types.String                          `tfsdk:"connection_id"`
	SchemaOptions           []JobSourcePostgresSourceSchemaOption `tfsdk:"schema_options"`
}
type JobSourcePostgresSourceSchemaOption struct {
	Schema types.String                          `tfsdk:"schema"`
	Tables []*JobSourcePostgresSourceTableOption `tfsdk:"tables"`
}
type JobSourcePostgresSourceTableOption struct {
	Table       types.String `tfsdk:"table"`
	WhereClause types.String `tfsdk:"where_clause"`
}

type JobSourceMysqlOptions struct{}
type JobSourceGenerateOptions struct{}
type JobSourceAwsS3Options struct{}

type JobDestination struct {
	Postgres *JobDestinationPostgresOptions `tfsdk:"postgres"`
	Mysql    *JobDestinationMysqlOptions    `tfsdk:"mysql"`
	AwsS3    *JobDestinationAwsS3Options    `tfsdk:"aws_s3"`
}
type JobDestinationPostgresOptions struct {
	TruncateTable   *PostgresDestinationTruncateTable `tfsdk:"truncate_table"`
	InitTableSchema types.Bool                        `tfsdk:"init_table_schema"`
}
type PostgresDestinationTruncateTable struct {
	TruncateBeforeInsert types.Bool `tfsdk:"truncate_before_insert"`
	Cascade              types.Bool `tfsdk:"cascade"`
}
type JobDestinationMysqlOptions struct{}
type JobDestinationAwsS3Options struct{}

type JobMapping struct {
	Schema      types.String `tfsdk:"schema"`
	Table       types.String `tfsdk:"table"`
	Column      types.String `tfsdk:"column"`
	Transformer *Transformer `tfsdk:"transformer"`
}
type Transformer struct {
	Source string             `tfsdk:"source"`
	Config *TransformerConfig `tfsdk:"config"`
}
type TransformerConfig struct{} // todo
type ActivityOptions struct {
	ScheduleToCloseTimeout types.Int64  `tfsdk:"schedule_to_close_timeout"`
	StartToCloseTimeout    types.Int64  `tfsdk:"start_to_close_timeout"`
	RetryPolicy            *RetryPolicy `tfsdk:"retry_policy"`
}
type WorkflowOptions struct {
	RunTimeout types.Int64 `tfsdk:"run_timeout"`
}
type RetryPolicy struct {
	MaximumAttempts types.Int64 `tfsdk:"maximum_attempts"`
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
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},

			"source": schema.SingleNestedAttribute{
				Description: "",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"postgres": schema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"halt_on_new_column_addition": schema.BoolAttribute{
								Description: "",
								Optional:    true,
							},
							"connection_id": schema.StringAttribute{
								Description: "",
								Required:    true,
							},
							"schema_options": schema.ListNestedAttribute{
								Description: "",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"schema": schema.StringAttribute{
											Description: "",
											Required:    true,
										},
										"tables": schema.ListNestedAttribute{
											Description: "",
											Required:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"table": schema.StringAttribute{
														Description: "",
														Required:    true,
													},
													"where_clause": schema.StringAttribute{
														Description: "",
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
						Description: "",
						Optional:    true,
					},
					"aws_s3": schema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
					},
				},
			},
			"destinations": schema.ListNestedAttribute{
				Description: "",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"postgres": schema.SingleNestedAttribute{
							Description: "",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"truncate_table": schema.SingleNestedAttribute{
									Description: "",
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										"truncate_before_insert": schema.BoolAttribute{
											Description: "",
											Optional:    true,
										},
										"cascade": schema.BoolAttribute{
											Description: "",
											Optional:    true,
										},
									},
								},
								"init_table_schema": schema.BoolAttribute{
									Description: "",
									Optional:    true,
								},
							},
						},
						"mysql": schema.SingleNestedAttribute{
							Description: "",
							Optional:    true,
						},
						"aws_s3": schema.SingleNestedAttribute{
							Description: "",
							Optional:    true,
						},
					},
				},
			},
			"mappings": schema.ListNestedAttribute{
				Description: "",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"schema": schema.StringAttribute{
							Description: "",
							Required:    true,
						},
						"table": schema.StringAttribute{
							Description: "",
							Required:    true,
						},
						"column": schema.StringAttribute{
							Description: "",
							Required:    true,
						},
						"transformer": schema.SingleNestedAttribute{
							Description: "",
							Required:    true,
							Attributes: map[string]schema.Attribute{
								"source": schema.StringAttribute{
									Description: "",
									Required:    true,
								},
								"config": schema.SingleNestedAttribute{
									Description: "",
									Required:    true,
									Attributes:  map[string]schema.Attribute{}, // todo
								},
							},
						},
					},
				},
			},

			"cron_schedule": schema.StringAttribute{
				Description: "A cron string for how often it's desired to schedule the job to run",
				Optional:    true,
			},

			"sync_options": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"schedule_to_close_timeout": schema.Int64Attribute{
						Description: "",
						Optional:    true,
					},
					"start_to_close_timeout": schema.Int64Attribute{
						Description: "",
						Optional:    true,
					},
					"retry_policy": schema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"maximum_attempts": schema.Int64Attribute{
								Description: "",
								Optional:    true,
							},
						},
					},
				},
			},
			"workflow_options": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"run_timeout": schema.Int64Attribute{
						Description: "",
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
