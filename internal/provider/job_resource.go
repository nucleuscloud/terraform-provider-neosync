package provider

import (
	"context"
	"errors"
	"fmt"
	"math"

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
	Id              types.String      `tfsdk:"id"`
	Name            types.String      `tfsdk:"name"`
	AccountId       types.String      `tfsdk:"account_id"`
	JobSource       *JobSource        `tfsdk:"source"`
	Destinations    []*JobDestination `tfsdk:"destinations"`
	Mappings        []*JobMapping     `tfsdk:"mappings"`
	CronSchedule    types.String      `tfsdk:"cron_schedule"`
	SyncOptions     *ActivityOptions  `tfsdk:"sync_options"`
	WorkflowOptions *WorkflowOptions  `tfsdk:"workflow_options"`
}

func toCreateJobDto(model *JobResourceModel) (*mgmtv1alpha1.CreateJobRequest, error) {
	if model == nil {
		return nil, errors.New("model was nil")
	}

	mappings, err := fromModelJobMappings(model.Mappings)
	if err != nil {
		return nil, err
	}
	source, err := fromModelJobSource(model.JobSource)
	if err != nil {
		return nil, err
	}

	workflowOpts, err := fromModelWorkflowOptions(model.WorkflowOptions)
	if err != nil {
		return nil, err
	}
	syncOpts, err := fromModelSyncOptions(model.SyncOptions)
	if err != nil {
		return nil, err
	}
	destinations, err := fromModelDestinationsToCreate(model.Destinations)
	if err != nil {
		return nil, err
	}

	return &mgmtv1alpha1.CreateJobRequest{
		JobName:         model.Name.ValueString(),
		AccountId:       model.AccountId.ValueString(),
		CronSchedule:    model.CronSchedule.ValueStringPointer(),
		Mappings:        mappings,
		Source:          source,
		Destinations:    destinations,
		InitiateJobRun:  false,
		WorkflowOptions: workflowOpts,
		SyncOptions:     syncOpts,
	}, nil
}

func fromModelDestinationsToCreate(input []*JobDestination) ([]*mgmtv1alpha1.CreateJobDestination, error) {
	output := []*mgmtv1alpha1.CreateJobDestination{}
	for _, jd := range input {
		cjd := &mgmtv1alpha1.CreateJobDestination{
			ConnectionId: jd.ConnectionId.ValueString(),
			Options:      &mgmtv1alpha1.JobDestinationOptions{},
		}
		if jd.Postgres != nil {
			var truncateTable *mgmtv1alpha1.PostgresTruncateTableConfig
			if jd.Postgres.TruncateTable != nil {
				truncateTable = &mgmtv1alpha1.PostgresTruncateTableConfig{
					TruncateBeforeInsert: jd.Postgres.TruncateTable.TruncateBeforeInsert.ValueBool(),
					Cascade:              jd.Postgres.TruncateTable.Cascade.ValueBool(),
				}
			}
			cjd.Options.Config = &mgmtv1alpha1.JobDestinationOptions_PostgresOptions{
				PostgresOptions: &mgmtv1alpha1.PostgresDestinationConnectionOptions{
					InitTableSchema: jd.Postgres.InitTableSchema.ValueBool(),
					TruncateTable:   truncateTable,
				},
			}
		} else if jd.AwsS3 != nil {
			cjd.Options.Config = &mgmtv1alpha1.JobDestinationOptions_AwsS3Options{
				AwsS3Options: &mgmtv1alpha1.AwsS3DestinationConnectionOptions{},
			}
		} else if jd.Mysql != nil {
			var truncateTable *mgmtv1alpha1.MysqlTruncateTableConfig
			if jd.Mysql.TruncateTable != nil {
				truncateTable = &mgmtv1alpha1.MysqlTruncateTableConfig{
					TruncateBeforeInsert: jd.Mysql.TruncateTable.TruncateBeforeInsert.ValueBool(),
				}
			}
			cjd.Options.Config = &mgmtv1alpha1.JobDestinationOptions_MysqlOptions{
				MysqlOptions: &mgmtv1alpha1.MysqlDestinationConnectionOptions{
					InitTableSchema: jd.Mysql.InitTableSchema.ValueBool(),
					TruncateTable:   truncateTable,
				},
			}
		} else {
			return nil, fmt.Errorf("the provided job destination type is not currently supported by this provider: %w", errors.ErrUnsupported)
		}

		output = append(output, cjd)
	}
	return output, nil
}
func fromModelDestinations(input []*JobDestination) ([]*mgmtv1alpha1.JobDestination, error) {
	output := []*mgmtv1alpha1.JobDestination{}

	cjds, err := fromModelDestinationsToCreate(input)
	if err != nil {
		return nil, err
	}

	for idx, jd := range input {
		output = append(output, &mgmtv1alpha1.JobDestination{
			Id:           jd.Id.ValueString(),
			ConnectionId: cjds[idx].ConnectionId,
			Options:      cjds[idx].Options,
		})
	}
	return output, nil
}

func fromModelSyncOptions(input *ActivityOptions) (*mgmtv1alpha1.ActivityOptions, error) { //nolint
	if input == nil {
		return nil, nil
	}

	if input.ScheduleToCloseTimeout.IsUnknown() && input.StartToCloseTimeout.IsUnknown() && input.RetryPolicy == nil {
		return nil, nil
	}

	output := &mgmtv1alpha1.ActivityOptions{
		ScheduleToCloseTimeout: input.ScheduleToCloseTimeout.ValueInt64Pointer(),
		StartToCloseTimeout:    input.StartToCloseTimeout.ValueInt64Pointer(),
		RetryPolicy:            &mgmtv1alpha1.RetryPolicy{},
	}
	if input.RetryPolicy != nil {
		output.RetryPolicy.MaximumAttempts = i64Toi32(input.RetryPolicy.MaximumAttempts.ValueInt64Pointer())
	}
	return output, nil
}

func fromModelWorkflowOptions(input *WorkflowOptions) (*mgmtv1alpha1.WorkflowOptions, error) { //nolint
	if input == nil {
		return nil, nil
	}
	if input.RunTimeout.IsUnknown() {
		return nil, nil
	}

	output := &mgmtv1alpha1.WorkflowOptions{
		RunTimeout: input.RunTimeout.ValueInt64Pointer(),
	}

	return output, nil
}

func fromModelJobSource(input *JobSource) (*mgmtv1alpha1.JobSource, error) {
	if input == nil {
		return nil, errors.New("input job source was nil")
	}

	output := &mgmtv1alpha1.JobSource{
		Options: &mgmtv1alpha1.JobSourceOptions{},
	}
	if input.Postgres != nil {
		schemas := []*mgmtv1alpha1.PostgresSourceSchemaOption{}
		for _, schemaOpt := range input.Postgres.SchemaOptions {
			tables := []*mgmtv1alpha1.PostgresSourceTableOption{}
			for _, tableOpt := range schemaOpt.Tables {
				tables = append(tables, &mgmtv1alpha1.PostgresSourceTableOption{
					Table:       tableOpt.Table.ValueString(),
					WhereClause: tableOpt.WhereClause.ValueStringPointer(),
				})
			}
			schemas = append(schemas, &mgmtv1alpha1.PostgresSourceSchemaOption{
				Schema: schemaOpt.Schema.ValueString(),
				Tables: tables,
			})
		}
		output.Options.Config = &mgmtv1alpha1.JobSourceOptions_Postgres{
			Postgres: &mgmtv1alpha1.PostgresSourceConnectionOptions{
				HaltOnNewColumnAddition: input.Postgres.HaltOnNewColumnAddition.ValueBool(),
				ConnectionId:            input.Postgres.ConnectionId.ValueString(),
				Schemas:                 schemas,
			},
		}
	} else if input.Mysql != nil {
		schemas := []*mgmtv1alpha1.MysqlSourceSchemaOption{}
		for _, schemaOpt := range input.Mysql.SchemaOptions {
			tables := []*mgmtv1alpha1.MysqlSourceTableOption{}
			for _, tableOpt := range schemaOpt.Tables {
				tables = append(tables, &mgmtv1alpha1.MysqlSourceTableOption{
					Table:       tableOpt.Table.ValueString(),
					WhereClause: tableOpt.WhereClause.ValueStringPointer(),
				})
			}
			schemas = append(schemas, &mgmtv1alpha1.MysqlSourceSchemaOption{
				Schema: schemaOpt.Schema.ValueString(),
				Tables: tables,
			})
		}
		output.Options.Config = &mgmtv1alpha1.JobSourceOptions_Mysql{
			Mysql: &mgmtv1alpha1.MysqlSourceConnectionOptions{
				HaltOnNewColumnAddition: input.Mysql.HaltOnNewColumnAddition.ValueBool(),
				ConnectionId:            input.Mysql.ConnectionId.ValueString(),
				Schemas:                 schemas,
			},
		}
	} else if input.AwsS3 != nil {
		output.Options.Config = &mgmtv1alpha1.JobSourceOptions_AwsS3{
			AwsS3: &mgmtv1alpha1.AwsS3SourceConnectionOptions{
				ConnectionId: input.AwsS3.ConnectionId.ValueString(),
			},
		}
	} else if input.Generate != nil {
		schemas := []*mgmtv1alpha1.GenerateSourceSchemaOption{}
		for _, schemaOpt := range input.Generate.Schemas {
			tables := []*mgmtv1alpha1.GenerateSourceTableOption{}
			for _, tableOpt := range schemaOpt.Tables {
				tables = append(tables, &mgmtv1alpha1.GenerateSourceTableOption{
					Table:    tableOpt.Table.ValueString(),
					RowCount: tableOpt.RowCount.ValueInt64(),
				})
			}
			schemas = append(schemas, &mgmtv1alpha1.GenerateSourceSchemaOption{
				Schema: schemaOpt.Schema.ValueString(),
				Tables: tables,
			})
		}
		output.Options.Config = &mgmtv1alpha1.JobSourceOptions_Generate{
			Generate: &mgmtv1alpha1.GenerateSourceOptions{
				FkSourceConnectionId: input.Generate.FkSourceConnectionId.ValueStringPointer(),
				Schemas:              schemas,
			},
		}
	} else {
		return nil, fmt.Errorf("the provided job source input is not currently supported by this provider: %w", errors.ErrUnsupported)
	}
	return output, nil
}

func fromModelJobMappings(input []*JobMapping) ([]*mgmtv1alpha1.JobMapping, error) {
	output := []*mgmtv1alpha1.JobMapping{}

	for _, inputMapping := range input {
		if inputMapping.Transformer == nil || inputMapping.Transformer.Config == nil {
			return nil, errors.New("transformer on input mapping was nil")
		}
		mapping := &mgmtv1alpha1.JobMapping{
			Schema: inputMapping.Schema.ValueString(),
			Table:  inputMapping.Table.ValueString(),
			Column: inputMapping.Column.ValueString(),
			Transformer: &mgmtv1alpha1.JobMappingTransformer{
				Source: inputMapping.Transformer.Source.ValueString(),
				Config: &mgmtv1alpha1.TransformerConfig{},
			},
		}
		config := &mgmtv1alpha1.TransformerConfig{}
		if inputMapping.Transformer.Config.Passthrough != nil {
			config.Config = &mgmtv1alpha1.TransformerConfig_PassthroughConfig{}
		} else {
			return nil, fmt.Errorf("the provided transformer config is not supported: %w", errors.ErrUnsupported)
		}

		mapping.Transformer.Config = config
		output = append(output, mapping)
	}

	return output, nil
}

func fromJobDto(dto *mgmtv1alpha1.Job) (*JobResourceModel, error) {
	if dto == nil {
		return nil, errors.New("dto was nil")
	}
	model := &JobResourceModel{
		Id:           types.StringValue(dto.Id),
		Name:         types.StringValue(dto.Name),
		AccountId:    types.StringValue(dto.AccountId),
		JobSource:    &JobSource{},
		Destinations: []*JobDestination{},
		Mappings:     []*JobMapping{},
		CronSchedule: types.StringPointerValue(dto.CronSchedule),
	}

	switch source := dto.Source.Options.Config.(type) {
	case *mgmtv1alpha1.JobSourceOptions_Postgres:
		model.JobSource.Postgres = &JobSourcePostgresOptions{
			HaltOnNewColumnAddition: types.BoolValue(source.Postgres.HaltOnNewColumnAddition),
			ConnectionId:            types.StringValue(source.Postgres.ConnectionId),
		}
		schemaOpts := []*JobSourcePostgresSourceSchemaOption{}
		for _, dtoopt := range source.Postgres.Schemas {
			opt := &JobSourcePostgresSourceSchemaOption{
				Schema: types.StringValue(dtoopt.Schema),
				Tables: []*JobSourcePostgresSourceTableOption{},
			}
			for _, schemaOpt := range dtoopt.Tables {
				opt.Tables = append(opt.Tables, &JobSourcePostgresSourceTableOption{
					Table:       types.StringValue(schemaOpt.Table),
					WhereClause: types.StringPointerValue(schemaOpt.WhereClause),
				})
			}
			schemaOpts = append(model.JobSource.Postgres.SchemaOptions, opt)
		}
		if len(schemaOpts) > 0 {
			model.JobSource.Postgres.SchemaOptions = schemaOpts
		}
	case *mgmtv1alpha1.JobSourceOptions_Mysql:
		model.JobSource.Mysql = &JobSourceMysqlOptions{
			HaltOnNewColumnAddition: types.BoolValue(source.Mysql.HaltOnNewColumnAddition),
			ConnectionId:            types.StringValue(source.Mysql.ConnectionId),
		}
		schemaOpts := []*JobSourceMysqlSourceSchemaOption{}
		for _, dtoopt := range source.Mysql.Schemas {
			opt := &JobSourceMysqlSourceSchemaOption{
				Schema: types.StringValue(dtoopt.Schema),
				Tables: []*JobSourceMysqlSourceTableOption{},
			}
			for _, schemaOpt := range dtoopt.Tables {
				opt.Tables = append(opt.Tables, &JobSourceMysqlSourceTableOption{
					Table:       types.StringValue(schemaOpt.Table),
					WhereClause: types.StringPointerValue(schemaOpt.WhereClause),
				})
			}
			schemaOpts = append(model.JobSource.Mysql.SchemaOptions, opt)
		}
		if len(schemaOpts) > 0 {
			model.JobSource.Mysql.SchemaOptions = schemaOpts
		}
	case *mgmtv1alpha1.JobSourceOptions_Generate:
		model.JobSource.Generate = &JobSourceGenerateOptions{
			FkSourceConnectionId: types.StringPointerValue(source.Generate.FkSourceConnectionId),
		}
		schemaOpts := []*JobSourceGenerateSchemaOption{}
		for _, dtoopt := range source.Generate.Schemas {
			opt := &JobSourceGenerateSchemaOption{
				Schema: types.StringValue(dtoopt.Schema),
				Tables: []*JobSourceGenerateTableOption{},
			}
			for _, schemaOpt := range dtoopt.Tables {
				opt.Tables = append(opt.Tables, &JobSourceGenerateTableOption{
					Table:    types.StringValue(schemaOpt.Table),
					RowCount: types.Int64Value(schemaOpt.RowCount),
				})
			}
			schemaOpts = append(model.JobSource.Generate.Schemas, opt)
		}
		if len(schemaOpts) > 0 {
			model.JobSource.Generate.Schemas = schemaOpts
		}
	case *mgmtv1alpha1.JobSourceOptions_AwsS3:
		model.JobSource.AwsS3 = &JobSourceAwsS3Options{
			ConnectionId: types.StringValue(source.AwsS3.ConnectionId),
		}
	default:
		return nil, fmt.Errorf("this job source is not currently supported by this provider: %w", errors.ErrUnsupported)
	}
	for _, dtoDest := range dto.Destinations {
		dest := &JobDestination{
			Id:           types.StringValue(dtoDest.Id),
			ConnectionId: types.StringValue(dtoDest.ConnectionId),
		}

		switch opt := dtoDest.Options.Config.(type) {
		case *mgmtv1alpha1.JobDestinationOptions_PostgresOptions:
			dest.Postgres = &JobDestinationPostgresOptions{
				InitTableSchema: types.BoolValue(opt.PostgresOptions.InitTableSchema),
			}
			if opt.PostgresOptions.TruncateTable != nil {
				dest.Postgres.TruncateTable = &PostgresDestinationTruncateTable{
					TruncateBeforeInsert: types.BoolValue(opt.PostgresOptions.TruncateTable.TruncateBeforeInsert),
					Cascade:              types.BoolValue(opt.PostgresOptions.TruncateTable.Cascade),
				}
			}
		case *mgmtv1alpha1.JobDestinationOptions_AwsS3Options:
			dest.AwsS3 = &JobDestinationAwsS3Options{}
		case *mgmtv1alpha1.JobDestinationOptions_MysqlOptions:
			dest.Mysql = &JobDestinationMysqlOptions{
				InitTableSchema: types.BoolValue(opt.MysqlOptions.InitTableSchema),
			}
			if opt.MysqlOptions.TruncateTable != nil {
				dest.Mysql.TruncateTable = &MysqlDestinationTruncateTable{
					TruncateBeforeInsert: types.BoolValue(opt.MysqlOptions.TruncateTable.TruncateBeforeInsert),
				}
			}
		default:
			return nil, fmt.Errorf("this job dest is not currently supported by this provider: %w", errors.ErrUnsupported)
		}

		model.Destinations = append(model.Destinations, dest)
	}
	for _, dtoMapping := range dto.Mappings {
		tconfig := &TransformerConfig{}
		switch dtoMapping.Transformer.Config.Config.(type) {
		case *mgmtv1alpha1.TransformerConfig_PassthroughConfig:
			tconfig.Passthrough = &TransformerEmpty{}
		default:
			return nil, fmt.Errorf("this job mapping transformer is not currently supported by this provider: %w", errors.ErrUnsupported)
		}
		mapping := &JobMapping{
			Schema: types.StringValue(dtoMapping.Schema),
			Table:  types.StringValue(dtoMapping.Table),
			Column: types.StringValue(dtoMapping.Column),
			Transformer: &Transformer{
				Source: types.StringValue(dtoMapping.Transformer.Source),
				Config: tconfig,
			},
		}
		model.Mappings = append(model.Mappings, mapping)
	}

	if dto.SyncOptions != nil && dto.SyncOptions.ScheduleToCloseTimeout != nil && dto.SyncOptions.StartToCloseTimeout != nil && dto.SyncOptions.RetryPolicy != nil {
		model.SyncOptions = &ActivityOptions{
			ScheduleToCloseTimeout: types.Int64PointerValue(dto.SyncOptions.ScheduleToCloseTimeout),
			StartToCloseTimeout:    types.Int64PointerValue(dto.SyncOptions.StartToCloseTimeout),
		}
		if dto.SyncOptions.RetryPolicy != nil {
			model.SyncOptions.RetryPolicy = &RetryPolicy{
				MaximumAttempts: types.Int64PointerValue(i32Toi64(dto.SyncOptions.RetryPolicy.MaximumAttempts)),
			}
		}
	}
	if dto.WorkflowOptions != nil && dto.WorkflowOptions.RunTimeout != nil {
		model.WorkflowOptions = &WorkflowOptions{
			RunTimeout: types.Int64PointerValue(dto.WorkflowOptions.RunTimeout),
		}
	}
	return model, nil
}

func i32Toi64(input *int32) *int64 {
	if input == nil {
		return nil
	}
	output := int64(*input)
	return &output
}

// if input is unsafe, returns nil.
func i64Toi32(input *int64) *int32 {
	if input == nil {
		return nil
	}

	if *input < math.MinInt32 || *input > math.MaxInt32 {
		return nil
	}
	output := int32(*input)
	return &output
}

type JobSource struct {
	Postgres *JobSourcePostgresOptions `tfsdk:"postgres"`
	Mysql    *JobSourceMysqlOptions    `tfsdk:"mysql"`
	Generate *JobSourceGenerateOptions `tfsdk:"generate"`
	AwsS3    *JobSourceAwsS3Options    `tfsdk:"aws_s3"`
}
type JobSourcePostgresOptions struct {
	HaltOnNewColumnAddition types.Bool                             `tfsdk:"halt_on_new_column_addition"`
	ConnectionId            types.String                           `tfsdk:"connection_id"`
	SchemaOptions           []*JobSourcePostgresSourceSchemaOption `tfsdk:"schema_options"`
}
type JobSourcePostgresSourceSchemaOption struct {
	Schema types.String                          `tfsdk:"schema"`
	Tables []*JobSourcePostgresSourceTableOption `tfsdk:"tables"`
}
type JobSourcePostgresSourceTableOption struct {
	Table       types.String `tfsdk:"table"`
	WhereClause types.String `tfsdk:"where_clause"`
}

type JobSourceMysqlOptions struct {
	HaltOnNewColumnAddition types.Bool                          `tfsdk:"halt_on_new_column_addition"`
	ConnectionId            types.String                        `tfsdk:"connection_id"`
	SchemaOptions           []*JobSourceMysqlSourceSchemaOption `tfsdk:"schema_options"`
}
type JobSourceMysqlSourceSchemaOption struct {
	Schema types.String                       `tfsdk:"schema"`
	Tables []*JobSourceMysqlSourceTableOption `tfsdk:"tables"`
}
type JobSourceMysqlSourceTableOption struct {
	Table       types.String `tfsdk:"table"`
	WhereClause types.String `tfsdk:"where_clause"`
}

type JobSourceGenerateOptions struct {
	Schemas              []*JobSourceGenerateSchemaOption `tfsdk:"schemas"`
	FkSourceConnectionId types.String                     `tfsdk:"fk_source_connection_id"`
}
type JobSourceGenerateSchemaOption struct {
	Schema types.String                    `tfsdk:"schema"`
	Tables []*JobSourceGenerateTableOption `tfsdk:"tables"`
}
type JobSourceGenerateTableOption struct {
	Table    types.String `tfsdk:"table"`
	RowCount types.Int64  `tfsdk:"row_count"`
}
type JobSourceAwsS3Options struct {
	ConnectionId types.String `tfsdk:"connection_id"`
}

type JobDestination struct {
	Id           types.String `tfsdk:"id"`
	ConnectionId types.String `tfsdk:"connection_id"`

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
type JobDestinationMysqlOptions struct {
	TruncateTable   *MysqlDestinationTruncateTable `tfsdk:"truncate_table"`
	InitTableSchema types.Bool                     `tfsdk:"init_table_schema"`
}
type MysqlDestinationTruncateTable struct {
	TruncateBeforeInsert types.Bool `tfsdk:"truncate_before_insert"`
}
type JobDestinationAwsS3Options struct{}

type JobMapping struct {
	Schema      types.String `tfsdk:"schema"`
	Table       types.String `tfsdk:"table"`
	Column      types.String `tfsdk:"column"`
	Transformer *Transformer `tfsdk:"transformer"`
}
type Transformer struct {
	Source types.String       `tfsdk:"source"`
	Config *TransformerConfig `tfsdk:"config"`
}

// todo: maybe flatten this config to just live on the transformer
type TransformerConfig struct {
	Passthrough *TransformerEmpty `tfsdk:"passthrough"`
}
type TransformerEmpty struct{}
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
								Required:    true,
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
						Attributes: map[string]schema.Attribute{
							"halt_on_new_column_addition": schema.BoolAttribute{
								Description: "",
								Required:    true,
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
					"aws_s3": schema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"connection_id": schema.StringAttribute{
								Description: "",
								Required:    true,
							},
						},
					},
					"generate": schema.SingleNestedAttribute{
						Description: "",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"connection_id": schema.StringAttribute{
								Description: "",
								Required:    true,
							},
							"schema_options": schema.ListNestedAttribute{
								Description: "",
								Required:    true,
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
													"row_count": schema.Int64Attribute{
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
				},
			},
			"destinations": schema.ListNestedAttribute{
				Description: "",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
						},
						"connection_id": schema.StringAttribute{
							Description: "",
							Optional:    true, // required if id is not set
						},
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
									Required:    true,
								},
							},
						},
						"mysql": schema.SingleNestedAttribute{
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
									},
								},
								"init_table_schema": schema.BoolAttribute{
									Description: "",
									Required:    true,
								},
							},
						},
						"aws_s3": schema.SingleNestedAttribute{
							Description: "",
							Optional:    true,
							Attributes:  map[string]schema.Attribute{},
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
									Attributes: map[string]schema.Attribute{
										"passthrough": schema.SingleNestedAttribute{
											Description: "",
											Optional:    true,
										},
									}, // todo
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
	data.AccountId = types.StringValue(accountId)

	jobRequest, err := toCreateJobDto(&data)
	if err != nil {
		resp.Diagnostics.AddError("unable to create job request", err.Error())
		return
	}

	jobResp, err := r.client.CreateJob(ctx, connect.NewRequest(jobRequest))
	if err != nil {
		resp.Diagnostics.AddError("create job error", err.Error())
		return
	}

	job := jobResp.Msg.Job

	newModel, err := fromJobDto(job)
	if err != nil {
		resp.Diagnostics.AddError("job translate error", err.Error())
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created job resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, newModel)...)
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

	updatedModel, err := fromJobDto(job)
	if err != nil {
		resp.Diagnostics.AddError("unable to convert dto to state", err.Error())
		return
	}

	// data.Id = types.StringValue(job.Id)
	// data.Name = types.StringValue(job.Name)
	// data.AccountId = types.StringValue(job.AccountId)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *JobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel JobResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateModel JobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if planModel.CronSchedule.ValueString() != stateModel.CronSchedule.ValueString() {
		_, err := r.client.UpdateJobSchedule(ctx, connect.NewRequest(&mgmtv1alpha1.UpdateJobScheduleRequest{
			Id:           planModel.Id.ValueString(),
			CronSchedule: planModel.CronSchedule.ValueStringPointer(),
		}))
		if err != nil {
			resp.Diagnostics.AddError("unable to update cron schedule", err.Error())
			return
		}
	}
	_, err := r.client.UpdateJobSourceConnection(ctx, connect.NewRequest(&mgmtv1alpha1.UpdateJobSourceConnectionRequest{
		Id:       planModel.Id.ValueString(),
		Source:   nil,
		Mappings: nil,
	}))
	if err != nil {
		resp.Diagnostics.AddError("unable to update job source connection", err.Error())
		return
	}
	destinationsToCreate := []*JobDestination{}
	destinationsToUpdate := []*JobDestination{}
	destinationsToDelete := []*JobDestination{}

	stateDestinationsMap := map[string]*JobDestination{}
	for _, dst := range stateModel.Destinations {
		stateDestinationsMap[dst.Id.ValueString()] = dst
	}

	for _, dst := range planModel.Destinations {
		if dst.Id.IsUnknown() {
			destinationsToCreate = append(destinationsToCreate, dst)
			continue
		}
		if _, ok := stateDestinationsMap[dst.Id.ValueString()]; !ok {
			destinationsToDelete = append(destinationsToDelete, dst)
			continue
		}
		destinationsToUpdate = append(destinationsToUpdate, dst) // should do work here to see if it has actually changed at all
	}

	if len(destinationsToCreate) > 0 {
		dsts, err := fromModelDestinationsToCreate(destinationsToCreate)
		if err != nil {
			resp.Diagnostics.AddError("unable to model new destinations to create", err.Error())
			return
		}
		_, err = r.client.CreateJobDestinationConnections(ctx, connect.NewRequest(&mgmtv1alpha1.CreateJobDestinationConnectionsRequest{
			JobId:        planModel.Id.ValueString(),
			Destinations: dsts,
		}))
		if err != nil {
			resp.Diagnostics.AddError("unable to create job destination connections", err.Error())
			return
		}
	}
	if len(destinationsToDelete) > 0 {
		for _, jd := range destinationsToDelete {
			_, err = r.client.DeleteJobDestinationConnection(ctx, connect.NewRequest(&mgmtv1alpha1.DeleteJobDestinationConnectionRequest{
				DestinationId: jd.Id.ValueString(),
			}))
			if err != nil {
				resp.Diagnostics.AddError("unable to delete job destination connection", err.Error())
				return
			}
		}
	}
	if len(destinationsToUpdate) > 0 {
		jds, err := fromModelDestinations(destinationsToUpdate)
		if err != nil {
			resp.Diagnostics.AddError("unable to model destinations to update", err.Error())
			return
		}
		for _, jd := range jds {
			_, err = r.client.UpdateJobDestinationConnection(ctx, connect.NewRequest(&mgmtv1alpha1.UpdateJobDestinationConnectionRequest{
				DestinationId: jd.Id,
				JobId:         planModel.Id.ValueString(),
				ConnectionId:  jd.ConnectionId,
				Options:       jd.Options,
			}))
			if err != nil {
				resp.Diagnostics.AddError("unable to update job destination connection", err.Error())
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

	updatedModel, err := fromJobDto(getResp.Msg.Job)
	if err != nil {
		resp.Diagnostics.AddError("unable to model latest job dto", err.Error())
		return
	}

	tflog.Trace(ctx, "updated job")
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &updatedModel)...)
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
