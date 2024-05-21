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
				Source: stateSourceToTransformerSource(inputMapping.Transformer.Source.ValueString()),
				Config: &mgmtv1alpha1.TransformerConfig{},
			},
		}
		config, err := fromModelTransformerConfig(inputMapping.Transformer.Config)
		if err != nil {
			return nil, err
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
		tconfig, err := toTransformerConfigFromDto(dtoMapping.Transformer.Config)
		if err != nil {
			return nil, err
		}
		mapping := &JobMapping{
			Schema: types.StringValue(dtoMapping.Schema),
			Table:  types.StringValue(dtoMapping.Table),
			Column: types.StringValue(dtoMapping.Column),
			Transformer: &Transformer{
				Source: types.StringValue(transformerSourceToStateSource(dtoMapping.Transformer.Source)),
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

func fromModelTransformerConfig(model *TransformerConfig) (*mgmtv1alpha1.TransformerConfig, error) {
	dto := &mgmtv1alpha1.TransformerConfig{}
	if model.GenerateEmail != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateEmailConfig{
			GenerateEmailConfig: &mgmtv1alpha1.GenerateEmail{},
		}
	} else if model.TransformEmail != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformEmailConfig{
			TransformEmailConfig: &mgmtv1alpha1.TransformEmail{
				PreserveDomain: model.TransformEmail.PreserveDomain.ValueBool(),
				PreserveLength: model.TransformEmail.PreserveLength.ValueBool(),
			},
		}
	} else if model.GenerateBool != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateBoolConfig{
			GenerateBoolConfig: &mgmtv1alpha1.GenerateBool{},
		}
	} else if model.GenerateCardNumber != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateCardNumberConfig{
			GenerateCardNumberConfig: &mgmtv1alpha1.GenerateCardNumber{
				ValidLuhn: model.GenerateCardNumber.ValidLuhn.ValueBool(),
			},
		}
	} else if model.GenerateCity != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateCityConfig{
			GenerateCityConfig: &mgmtv1alpha1.GenerateCity{},
		}
	} else if model.GenerateE164PhoneNumber != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateE164PhoneNumberConfig{
			GenerateE164PhoneNumberConfig: &mgmtv1alpha1.GenerateE164PhoneNumber{
				Min: model.GenerateE164PhoneNumber.Min.ValueInt64(),
				Max: model.GenerateE164PhoneNumber.Max.ValueInt64(),
			},
		}
	} else if model.GenerateFirstName != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateFirstNameConfig{
			GenerateFirstNameConfig: &mgmtv1alpha1.GenerateFirstName{},
		}
	} else if model.GenerateFloat64 != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateFloat64Config{
			GenerateFloat64Config: &mgmtv1alpha1.GenerateFloat64{
				RandomizeSign: model.GenerateFloat64.RandomizeSign.ValueBool(),
				Min:           model.GenerateFloat64.Min.ValueFloat64(),
				Max:           model.GenerateFloat64.Max.ValueFloat64(),
			},
		}
	} else if model.GenerateFullAddress != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateFullAddressConfig{
			GenerateFullAddressConfig: &mgmtv1alpha1.GenerateFullAddress{},
		}
	} else if model.GenerateFullName != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateFullNameConfig{
			GenerateFullNameConfig: &mgmtv1alpha1.GenerateFullName{},
		}
	} else if model.GenerateGender != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateGenderConfig{
			GenerateGenderConfig: &mgmtv1alpha1.GenerateGender{
				Abbreviate: model.GenerateGender.Abbreviate.ValueBool(),
			},
		}
	} else if model.GenerateInt64PhoneNumber != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateInt64PhoneNumberConfig{
			GenerateInt64PhoneNumberConfig: &mgmtv1alpha1.GenerateInt64PhoneNumber{},
		}
	} else if model.GenerateInt64 != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateInt64Config{
			GenerateInt64Config: &mgmtv1alpha1.GenerateInt64{
				RandomizeSign: model.GenerateInt64.RandomizeSign.ValueBool(),
				Min:           model.GenerateInt64.Min.ValueInt64(),
				Max:           model.GenerateInt64.Max.ValueInt64(),
			},
		}
	} else if model.GenerateLastName != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateLastNameConfig{
			GenerateLastNameConfig: &mgmtv1alpha1.GenerateLastName{},
		}
	} else if model.GenerateSha256Hash != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateSha256HashConfig{
			GenerateSha256HashConfig: &mgmtv1alpha1.GenerateSha256Hash{},
		}
	} else if model.GenerateSsn != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateSsnConfig{
			GenerateSsnConfig: &mgmtv1alpha1.GenerateSSN{},
		}
	} else if model.GenerateState != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateStateConfig{
			GenerateStateConfig: &mgmtv1alpha1.GenerateState{},
		}
	} else if model.GenerateStreetAddress != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateStreetAddressConfig{
			GenerateStreetAddressConfig: &mgmtv1alpha1.GenerateStreetAddress{},
		}
	} else if model.GenerateStringPhoneNumber != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateStringPhoneNumberConfig{
			GenerateStringPhoneNumberConfig: &mgmtv1alpha1.GenerateStringPhoneNumber{
				Min: model.GenerateStringPhoneNumber.Min.ValueInt64(),
				Max: model.GenerateStringPhoneNumber.Max.ValueInt64(),
			},
		}
	} else if model.GenerateString != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateStringConfig{
			GenerateStringConfig: &mgmtv1alpha1.GenerateString{
				Min: model.GenerateString.Min.ValueInt64(),
				Max: model.GenerateString.Max.ValueInt64(),
			},
		}
	} else if model.GenerateUnixtimestamp != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateUnixtimestampConfig{
			GenerateUnixtimestampConfig: &mgmtv1alpha1.GenerateUnixTimestamp{},
		}
	} else if model.GenerateUsername != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateUsernameConfig{
			GenerateUsernameConfig: &mgmtv1alpha1.GenerateUsername{},
		}
	} else if model.GenerateUtctimestamp != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateUtctimestampConfig{
			GenerateUtctimestampConfig: &mgmtv1alpha1.GenerateUtcTimestamp{},
		}
	} else if model.GenerateUuid != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateUuidConfig{
			GenerateUuidConfig: &mgmtv1alpha1.GenerateUuid{
				IncludeHyphens: model.GenerateUuid.IncludeHyphens.ValueBool(),
			},
		}
	} else if model.GenerateZipcode != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateZipcodeConfig{
			GenerateZipcodeConfig: &mgmtv1alpha1.GenerateZipcode{},
		}
	} else if model.TransformE164PhoneNumber != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateE164PhoneNumberConfig{
			GenerateE164PhoneNumberConfig: &mgmtv1alpha1.GenerateE164PhoneNumber{
				Min: model.GenerateE164PhoneNumber.Min.ValueInt64(),
				Max: model.GenerateE164PhoneNumber.Max.ValueInt64(),
			},
		}
	} else if model.TransformFirstName != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateFirstNameConfig{
			GenerateFirstNameConfig: &mgmtv1alpha1.GenerateFirstName{},
		}
	} else if model.TransformFloat64 != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformFloat64Config{
			TransformFloat64Config: &mgmtv1alpha1.TransformFloat64{
				RandomizationRangeMin: model.TransformFloat64.RandomizationRangeMin.ValueFloat64(),
				RandomizationRangeMax: model.TransformFloat64.RandomizationRangeMax.ValueFloat64(),
			},
		}
	} else if model.TransformFullName != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformFullNameConfig{
			TransformFullNameConfig: &mgmtv1alpha1.TransformFullName{
				PreserveLength: model.TransformFullName.PreserveLength.ValueBool(),
			},
		}
	} else if model.TransformInt64PhoneNumber != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformInt64PhoneNumberConfig{
			TransformInt64PhoneNumberConfig: &mgmtv1alpha1.TransformInt64PhoneNumber{
				PreserveLength: model.TransformInt64PhoneNumber.PreserveLength.ValueBool(),
			},
		}
	} else if model.TransformInt64 != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateInt64Config{
			GenerateInt64Config: &mgmtv1alpha1.GenerateInt64{
				RandomizeSign: model.GenerateInt64.RandomizeSign.ValueBool(),
				Min:           model.GenerateInt64.Min.ValueInt64(),
				Max:           model.GenerateInt64.Max.ValueInt64(),
			},
		}
	} else if model.TransformLastName != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformLastNameConfig{
			TransformLastNameConfig: &mgmtv1alpha1.TransformLastName{
				PreserveLength: model.TransformLastName.PreserveLength.ValueBool(),
			},
		}
	} else if model.TransformPhoneNumber != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformPhoneNumberConfig{
			TransformPhoneNumberConfig: &mgmtv1alpha1.TransformPhoneNumber{
				PreserveLength: model.TransformPhoneNumber.PreserveLength.ValueBool(),
			},
		}
	} else if model.TransformString != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformStringConfig{
			TransformStringConfig: &mgmtv1alpha1.TransformString{
				PreserveLength: model.TransformString.PreserveLength.ValueBool(),
			},
		}
	} else if model.Passthrough != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_PassthroughConfig{
			PassthroughConfig: &mgmtv1alpha1.Passthrough{},
		}
	} else if model.Null != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_Nullconfig{
			Nullconfig: &mgmtv1alpha1.Null{},
		}
	} else if model.UserDefinedTransformer != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_UserDefinedTransformerConfig{
			UserDefinedTransformerConfig: &mgmtv1alpha1.UserDefinedTransformerConfig{
				Id: model.UserDefinedTransformer.Id.ValueString(),
			},
		}
	} else if model.GenerateDefault != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateDefaultConfig{
			GenerateDefaultConfig: &mgmtv1alpha1.GenerateDefault{},
		}
	} else if model.TransformJavascript != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformJavascriptConfig{
			TransformJavascriptConfig: &mgmtv1alpha1.TransformJavascript{
				Code: model.TransformJavascript.Code.ValueString(),
			},
		}
	} else if model.GenerateCategorical != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateCategoricalConfig{
			GenerateCategoricalConfig: &mgmtv1alpha1.GenerateCategorical{
				Categories: model.GenerateCategorical.Categories.ValueString(),
			},
		}
	} else if model.TransformCharacterScramble != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformCharacterScrambleConfig{
			TransformCharacterScrambleConfig: &mgmtv1alpha1.TransformCharacterScramble{},
		}
	} else {
		return nil, fmt.Errorf("transformer config is not currently supported by this provider: %w", errors.ErrUnsupported)
	}

	return dto, nil
}

func toTransformerConfigFromDto(dto *mgmtv1alpha1.TransformerConfig) (*TransformerConfig, error) {
	tconfig := &TransformerConfig{}
	switch config := dto.Config.(type) {
	case *mgmtv1alpha1.TransformerConfig_GenerateEmailConfig:
		tconfig.GenerateEmail = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_TransformEmailConfig:
		tconfig.TransformEmail = &TransformEmail{}
	case *mgmtv1alpha1.TransformerConfig_GenerateBoolConfig:
		tconfig.GenerateBool = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateCardNumberConfig:
		tconfig.GenerateCardNumber = &GenerateCardNumber{
			ValidLuhn: types.BoolValue(config.GenerateCardNumberConfig.ValidLuhn),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateCityConfig:
		tconfig.GenerateCity = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateE164PhoneNumberConfig:
		tconfig.GenerateE164PhoneNumber = &GenerateE164PhoneNumber{
			Min: types.Int64Value(config.GenerateE164PhoneNumberConfig.Min),
			Max: types.Int64Value(config.GenerateE164PhoneNumberConfig.Max),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateFirstNameConfig:
		tconfig.GenerateFirstName = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateFloat64Config:
		tconfig.GenerateFloat64 = &GenerateFloat64{
			RandomizeSign: types.BoolValue(config.GenerateFloat64Config.RandomizeSign),
			Min:           types.Float64Value(config.GenerateFloat64Config.Min),
			Max:           types.Float64Value(config.GenerateFloat64Config.Max),
			Precision:     types.Int64Value(config.GenerateFloat64Config.Precision),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateFullAddressConfig:
		tconfig.GenerateFullAddress = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateFullNameConfig:
		tconfig.GenerateFullName = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateGenderConfig:
		tconfig.GenerateGender = &GenerateGender{
			Abbreviate: types.BoolValue(config.GenerateGenderConfig.Abbreviate),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateInt64PhoneNumberConfig:
		tconfig.GenerateInt64PhoneNumber = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateInt64Config:
		tconfig.GenerateInt64 = &GenerateInt64{
			RandomizeSign: types.BoolValue(config.GenerateInt64Config.RandomizeSign),
			Min:           types.Int64Value(config.GenerateInt64Config.Min),
			Max:           types.Int64Value(config.GenerateInt64Config.Max),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateLastNameConfig:
		tconfig.GenerateLastName = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateSha256HashConfig:
		tconfig.GenerateSha256Hash = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateSsnConfig:
		tconfig.GenerateSsn = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateStateConfig:
		tconfig.GenerateState = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateStreetAddressConfig:
		tconfig.GenerateStreetAddress = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateStringPhoneNumberConfig:
		tconfig.GenerateStringPhoneNumber = &GenerateStringPhoneNumber{
			Min: types.Int64Value(config.GenerateStringPhoneNumberConfig.Min),
			Max: types.Int64Value(config.GenerateStringPhoneNumberConfig.Max),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateStringConfig:
		tconfig.GenerateString = &GenerateString{
			Min: types.Int64Value(config.GenerateStringConfig.Min),
			Max: types.Int64Value(config.GenerateStringConfig.Max),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateUnixtimestampConfig:
		tconfig.GenerateUnixtimestamp = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateUsernameConfig:
		tconfig.GenerateUsername = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateUtctimestampConfig:
		tconfig.GenerateUtctimestamp = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateUuidConfig:
		tconfig.GenerateUuid = &GenerateUuid{
			IncludeHyphens: types.BoolValue(config.GenerateUuidConfig.IncludeHyphens),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateZipcodeConfig:
		tconfig.GenerateZipcode = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_TransformE164PhoneNumberConfig:
		tconfig.TransformE164PhoneNumber = &TransformE164PhoneNumber{
			PreserveLength: types.BoolValue(config.TransformE164PhoneNumberConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformFirstNameConfig:
		tconfig.TransformFirstName = &TransformFirstName{
			PreserveLength: types.BoolValue(config.TransformFirstNameConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformFloat64Config:
		tconfig.TransformFloat64 = &TransformFloat64{
			RandomizationRangeMin: types.Float64Value(config.TransformFloat64Config.RandomizationRangeMin),
			RandomizationRangeMax: types.Float64Value(config.TransformFloat64Config.RandomizationRangeMax),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformFullNameConfig:
		tconfig.TransformFullName = &TransformFullName{
			PreserveLength: types.BoolValue(config.TransformFullNameConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformInt64PhoneNumberConfig:
		tconfig.TransformInt64PhoneNumber = &TransformInt64PhoneNumber{
			PreserveLength: types.BoolValue(config.TransformInt64PhoneNumberConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformInt64Config:
		tconfig.TransformInt64 = &TransformInt64{
			RandomizationRangeMin: types.Int64Value(config.TransformInt64Config.RandomizationRangeMin),
			RandomizationRangeMax: types.Int64Value(config.TransformInt64Config.RandomizationRangeMax),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformLastNameConfig:
		tconfig.TransformLastName = &TransformLastName{
			PreserveLength: types.BoolValue(config.TransformLastNameConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformPhoneNumberConfig:
		tconfig.TransformPhoneNumber = &TransformPhoneNumber{
			PreserveLength: types.BoolValue(config.TransformPhoneNumberConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformStringConfig:
		tconfig.TransformString = &TransformString{
			PreserveLength: types.BoolValue(config.TransformStringConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_PassthroughConfig:
		tconfig.Passthrough = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_Nullconfig:
		tconfig.Null = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_UserDefinedTransformerConfig:
		tconfig.UserDefinedTransformer = &UserDefinedTransformer{
			Id: types.StringValue(config.UserDefinedTransformerConfig.Id),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateDefaultConfig:
		tconfig.GenerateDefault = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_TransformJavascriptConfig:
		tconfig.TransformJavascript = &TransformJavascript{
			Code: types.StringValue(config.TransformJavascriptConfig.Code),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateCategoricalConfig:
		tconfig.GenerateCategorical = &GenerateCategorical{
			Categories: types.StringValue(config.GenerateCategoricalConfig.Categories),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformCharacterScrambleConfig:
		tconfig.TransformCharacterScramble = &TransformerEmpty{}
	default:
		return nil, fmt.Errorf("this job mapping transformer is not currently supported by this provider: %w", errors.ErrUnsupported)
	}
	return tconfig, nil
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
	SchemaOptions           []*JobSourcePostgresSourceSchemaOption `tfsdk:"schemas"`
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
	SchemaOptions           []*JobSourceMysqlSourceSchemaOption `tfsdk:"schemas"`
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

type TransformerConfig struct {
	GenerateEmail              *TransformerEmpty          `tfsdk:"generate_email"`
	TransformEmail             *TransformEmail            `tfsdk:"transform_email"`
	GenerateBool               *TransformerEmpty          `tfsdk:"generate_bool"`
	GenerateCardNumber         *GenerateCardNumber        `tfsdk:"generate_card_number"`
	GenerateCity               *TransformerEmpty          `tfsdk:"generate_city"`
	GenerateE164PhoneNumber    *GenerateE164PhoneNumber   `tfsdk:"generate_e164_phone_number"`
	GenerateFirstName          *TransformerEmpty          `tfsdk:"generate_firstname"`
	GenerateFloat64            *GenerateFloat64           `tfsdk:"generate_float64"`
	GenerateFullAddress        *TransformerEmpty          `tfsdk:"generate_full_address"`
	GenerateFullName           *TransformerEmpty          `tfsdk:"generate_fullname"`
	GenerateGender             *GenerateGender            `tfsdk:"generate_gender"`
	GenerateInt64PhoneNumber   *TransformerEmpty          `tfsdk:"generate_int64_phone_number"`
	GenerateInt64              *GenerateInt64             `tfsdk:"generate_int64"`
	GenerateLastName           *TransformerEmpty          `tfsdk:"generate_lastname"`
	GenerateSha256Hash         *TransformerEmpty          `tfsdk:"generate_sha256"`
	GenerateSsn                *TransformerEmpty          `tfsdk:"generate_ssn"`
	GenerateState              *TransformerEmpty          `tfsdk:"generate_state"`
	GenerateStreetAddress      *TransformerEmpty          `tfsdk:"generate_street_address"`
	GenerateStringPhoneNumber  *GenerateStringPhoneNumber `tfsdk:"generate_string_phone_number"`
	GenerateString             *GenerateString            `tfsdk:"generate_string"`
	GenerateUnixtimestamp      *TransformerEmpty          `tfsdk:"generate_unix_timestamp"`
	GenerateUsername           *TransformerEmpty          `tfsdk:"generate_username"`
	GenerateUtctimestamp       *TransformerEmpty          `tfsdk:"generate_utc_timestamp"`
	GenerateUuid               *GenerateUuid              `tfsdk:"generate_uuid"`
	GenerateZipcode            *TransformerEmpty          `tfsdk:"generate_zipcode"`
	TransformE164PhoneNumber   *TransformE164PhoneNumber  `tfsdk:"transform_e164_phone_number"`
	TransformFirstName         *TransformFirstName        `tfsdk:"transform_firstname"`
	TransformFloat64           *TransformFloat64          `tfsdk:"transform_float64"`
	TransformFullName          *TransformFullName         `tfsdk:"transform_fullname"`
	TransformInt64PhoneNumber  *TransformInt64PhoneNumber `tfsdk:"transform_int64_phone_number"`
	TransformInt64             *TransformInt64            `tfsdk:"transform_int64"`
	TransformLastName          *TransformLastName         `tfsdk:"transform_lastname"`
	TransformPhoneNumber       *TransformPhoneNumber      `tfsdk:"transform_phone_number"`
	TransformString            *TransformString           `tfsdk:"transform_string"`
	Passthrough                *TransformerEmpty          `tfsdk:"passthrough"`
	Null                       *TransformerEmpty          `tfsdk:"null"`
	UserDefinedTransformer     *UserDefinedTransformer    `tfsdk:"user_defined_transformer"`
	GenerateDefault            *TransformerEmpty          `tfsdk:"generate_default"`
	TransformJavascript        *TransformJavascript       `tfsdk:"transform_javascript"`
	GenerateCategorical        *GenerateCategorical       `tfsdk:"generate_categorical"`
	TransformCharacterScramble *TransformerEmpty          `tfsdk:"transform_character_scramble"`
}
type TransformerEmpty struct{}
type TransformEmail struct {
	PreserveDomain types.Bool `tfsdk:"preserve_domain"`
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type GenerateCardNumber struct {
	ValidLuhn types.Bool `tfsdk:"valid_luhn"`
}
type GenerateE164PhoneNumber struct {
	Min types.Int64 `tfsdk:"min"`
	Max types.Int64 `tfsdk:"max"`
}
type GenerateFloat64 struct {
	RandomizeSign types.Bool    `tfsdk:"randomize_sign"`
	Min           types.Float64 `tfsdk:"min"`
	Max           types.Float64 `tfsdk:"max"`
	Precision     types.Int64   `tfsdk:"precision"`
}
type GenerateGender struct {
	Abbreviate types.Bool `tfsdk:"abbreviate"`
}
type GenerateInt64 struct {
	RandomizeSign types.Bool  `tfsdk:"randomize_sign"`
	Min           types.Int64 `tfsdk:"min"`
	Max           types.Int64 `tfsdk:"max"`
}
type GenerateStringPhoneNumber struct {
	Min types.Int64 `tfsdk:"min"`
	Max types.Int64 `tfsdk:"max"`
}
type GenerateString struct {
	Min types.Int64 `tfsdk:"min"`
	Max types.Int64 `tfsdk:"max"`
}
type GenerateUuid struct {
	IncludeHyphens types.Bool `tfsdk:"include_hyphens"`
}
type TransformE164PhoneNumber struct {
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type TransformFirstName struct {
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type TransformFloat64 struct {
	RandomizationRangeMin types.Float64 `tfsdk:"randomization_range_min"`
	RandomizationRangeMax types.Float64 `tfsdk:"randomization_range_max"`
}
type TransformFullName struct {
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type TransformInt64PhoneNumber struct {
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type TransformInt64 struct {
	RandomizationRangeMin types.Int64 `tfsdk:"randomization_range_min"`
	RandomizationRangeMax types.Int64 `tfsdk:"randomization_range_max"`
}
type TransformLastName struct {
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type TransformPhoneNumber struct {
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type TransformString struct {
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type TransformJavascript struct {
	Code types.String `tfsdk:"code"`
}
type UserDefinedTransformer struct {
	Id types.String `tfsdk:"id"`
}
type GenerateCategorical struct {
	Categories types.String `tfsdk:"categories"`
}

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
				Description: "Configuration details about the source data connection",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"postgres": schema.SingleNestedAttribute{
						Description: "Postgres specific connection configurations",
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
							Description: "The unique identifier of the destination resource. This is set after creation",
							Optional:    true,
							Computed:    true,
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
								"source": schema.StringAttribute{
									Description: "The source of the transformer that will be used",
									Required:    true,
								},
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

	newSource, err := fromModelJobSource(planModel.JobSource)
	if err != nil {
		resp.Diagnostics.AddError("unable to map new job source", err.Error())
		return
	}

	newMappings, err := fromModelJobMappings(planModel.Mappings)
	if err != nil {
		resp.Diagnostics.AddError("unable to map new job mappings", err.Error())
		return
	}

	_, err = r.client.UpdateJobSourceConnection(ctx, connect.NewRequest(&mgmtv1alpha1.UpdateJobSourceConnectionRequest{
		Id:       planModel.Id.ValueString(),
		Source:   newSource,
		Mappings: newMappings,
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
