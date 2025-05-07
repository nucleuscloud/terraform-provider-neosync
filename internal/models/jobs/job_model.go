package job_model

import (
	"errors"
	"math"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgmtv1alpha1 "github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1"
	transformer_model "github.com/nucleuscloud/terraform-provider-neosync/internal/models/transformers"
)

type JobResourceModel struct {
	Id                 types.String                   `tfsdk:"id"`
	Name               types.String                   `tfsdk:"name"`
	AccountId          types.String                   `tfsdk:"account_id"`
	JobSource          *JobSource                     `tfsdk:"source"`
	Destinations       []*JobDestination              `tfsdk:"destinations"`
	Mappings           []*JobMapping                  `tfsdk:"mappings"`
	CronSchedule       types.String                   `tfsdk:"cron_schedule"`
	SyncOptions        *ActivityOptions               `tfsdk:"sync_options"`
	WorkflowOptions    *WorkflowOptions               `tfsdk:"workflow_options"`
	VirtualForeignKeys []*VirtualForeignKeyConstraint `tfsdk:"virtual_foreign_keys"`
}

type VirtualForeignKeyConstraint struct {
	Schema     types.String       `tfsdk:"schema"`
	Table      types.String       `tfsdk:"table"`
	Columns    []types.String     `tfsdk:"columns"`
	ForeignKey *VirtualForeignKey `tfsdk:"foreign_key"`
}

type VirtualForeignKey struct {
	Schema  types.String   `tfsdk:"schema"`
	Table   types.String   `tfsdk:"table"`
	Columns []types.String `tfsdk:"columns"`
}

type JobSource struct {
	Postgres *JobSourcePostgresOptions `tfsdk:"postgres"`
	Mysql    *JobSourceMysqlOptions    `tfsdk:"mysql"`
	Generate *JobSourceGenerateOptions `tfsdk:"generate"`
	AwsS3    *JobSourceAwsS3Options    `tfsdk:"aws_s3"`
	// mssql, ai generate
}
type JobSourcePostgresOptions struct {
	NewColumnAdditionStrategy     *PostgresNewColumnAdditionStrategy     `tfsdk:"new_column_addition_strategy"`
	ColumnRemovalStrategy         *PostgresColumnRemovalStrategy         `tfsdk:"column_removal_strategy"`
	ConnectionId                  types.String                           `tfsdk:"connection_id"`
	SubsetByForeignKeyConstraints types.Bool                             `tfsdk:"subset_by_foreign_key_constraints"`
	SchemaOptions                 []*JobSourcePostgresSourceSchemaOption `tfsdk:"schemas"`
}

type PostgresNewColumnAdditionStrategy struct {
	HaltJob *PostgresNewColumnAdditionStrategyHaltJob `tfsdk:"halt_job"`
	AutoMap *PostgresNewColumnAdditionStrategyAutoMap `tfsdk:"auto_map"`
}

type PostgresNewColumnAdditionStrategyHaltJob struct{}
type PostgresNewColumnAdditionStrategyAutoMap struct{}

type PostgresColumnRemovalStrategy struct {
	HaltJob     *PostgresHaltJobColumnRemovalStrategy     `tfsdk:"halt_job"`
	ContinueJob *PostgresContinueJobColumnRemovalStrategy `tfsdk:"continue_job"`
}
type PostgresHaltJobColumnRemovalStrategy struct{}
type PostgresContinueJobColumnRemovalStrategy struct{}

type JobSourcePostgresSourceSchemaOption struct {
	Schema types.String                          `tfsdk:"schema"`
	Tables []*JobSourcePostgresSourceTableOption `tfsdk:"tables"`
}
type JobSourcePostgresSourceTableOption struct {
	Table       types.String `tfsdk:"table"`
	WhereClause types.String `tfsdk:"where_clause"`
}

type JobSourceMysqlOptions struct {
	ConnectionId                  types.String                        `tfsdk:"connection_id"`
	SchemaOptions                 []*JobSourceMysqlSourceSchemaOption `tfsdk:"schemas"`
	SubsetByForeignKeyConstraints types.Bool                          `tfsdk:"subset_by_foreign_key_constraints"`
	ColumnRemovalStrategy         *MssqlColumnRemovalStrategy         `tfsdk:"column_removal_strategy"`
}
type MssqlColumnRemovalStrategy struct {
	HaltJob     *MssqlHaltJobColumnRemovalStrategy     `json:"haltJob,omitempty"`
	ContinueJob *MssqlContinueJobColumnRemovalStrategy `json:"continueJob,omitempty"`
}
type MssqlHaltJobColumnRemovalStrategy struct{}
type MssqlContinueJobColumnRemovalStrategy struct{}

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
	// todo: fill out remaining destinations
}
type JobDestinationPostgresOptions struct {
	TruncateTable            *PostgresDestinationTruncateTable `tfsdk:"truncate_table"`
	InitTableSchema          types.Bool                        `tfsdk:"init_table_schema"`
	OnConflictConfig         *PostgresOnConflictConfig         `tfsdk:"on_conflict_config"`
	SkipForeignKeyViolations types.Bool                        `tfsdk:"skip_foreign_key_violations"`
	MaxInFlight              types.Int64                       `tfsdk:"max_in_flight"`
	Batch                    *BatchConfig                      `tfsdk:"batch"`
}

type BatchConfig struct {
	Count  types.Int64  `tfsdk:"count"`
	Period types.String `tfsdk:"period"`
}

type PostgresOnConflictConfig struct {
	OnConflictStrategy *PostgresOnConflictStrategy `tfsdk:"on_conflict_strategy"`
}

type PostgresOnConflictStrategy struct {
	Nothing *PostgresOnConflictNothing `tfsdk:"nothing"`
	Update  *PostgresOnConflictUpdate  `tfsdk:"update"`
}

type PostgresOnConflictNothing struct{}

type PostgresOnConflictUpdate struct {
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
	Schema      types.String                   `tfsdk:"schema"`
	Table       types.String                   `tfsdk:"table"`
	Column      types.String                   `tfsdk:"column"`
	Transformer *transformer_model.Transformer `tfsdk:"transformer"`
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

func (j *JobResourceModel) ToCreateJobDto() (*mgmtv1alpha1.CreateJobRequest, error) {
	if j == nil {
		return nil, errors.New("job resource model is nil")
	}

	mappings, err := ToJobMappingsDto(j.Mappings)
	if err != nil {
		return nil, err
	}
	source, err := j.JobSource.ToDto()
	if err != nil {
		return nil, err
	}

	var workflowOpts *mgmtv1alpha1.WorkflowOptions
	if j.WorkflowOptions != nil {
		workflowOpts, err = j.WorkflowOptions.ToDto()
		if err != nil {
			return nil, err
		}
	}

	var syncOpts *mgmtv1alpha1.ActivityOptions
	if j.SyncOptions != nil {
		syncOpts, err = j.SyncOptions.ToDto()
		if err != nil {
			return nil, err
		}
	}

	var destinations []*mgmtv1alpha1.CreateJobDestination
	if len(j.Destinations) > 0 {
		destinations = make([]*mgmtv1alpha1.CreateJobDestination, 0, len(j.Destinations))
		for _, destination := range j.Destinations {
			destinationDto, err := destination.ToCreateJobDestinationDto()
			if err != nil {
				return nil, err
			}
			destinations = append(destinations, destinationDto)
		}
	}

	var virtualForeignKeys []*mgmtv1alpha1.VirtualForeignConstraint
	if len(j.VirtualForeignKeys) > 0 {
		virtualForeignKeys = make([]*mgmtv1alpha1.VirtualForeignConstraint, 0, len(j.VirtualForeignKeys))
		for _, vfk := range j.VirtualForeignKeys {
			vfkDto, err := vfk.ToDto()
			if err != nil {
				return nil, err
			}
			virtualForeignKeys = append(virtualForeignKeys, vfkDto)
		}
	}

	return &mgmtv1alpha1.CreateJobRequest{
		JobName:            j.Name.ValueString(),
		AccountId:          j.AccountId.ValueString(),
		CronSchedule:       j.CronSchedule.ValueStringPointer(),
		Mappings:           mappings,
		Source:             source,
		Destinations:       destinations,
		InitiateJobRun:     false,
		WorkflowOptions:    workflowOpts,
		SyncOptions:        syncOpts,
		VirtualForeignKeys: virtualForeignKeys,
	}, nil
}

type UpdateJobRequest struct {
	UpdateJobScheduleRequest         *mgmtv1alpha1.UpdateJobScheduleRequest
	UpdateJobSourceConnectionRequest *mgmtv1alpha1.UpdateJobSourceConnectionRequest

	CreateJobDestinationConnectionsRequest *mgmtv1alpha1.CreateJobDestinationConnectionsRequest
	UpdateJobDestinationConnectionRequests []*mgmtv1alpha1.UpdateJobDestinationConnectionRequest
	DeleteJobDestinationConnectionRequests []*mgmtv1alpha1.DeleteJobDestinationConnectionRequest
}

func (j *JobResourceModel) ToUpdateJobDto(planModel *JobResourceModel, jobId string) (*UpdateJobRequest, error) {
	if j == nil {
		return nil, errors.New("job resource model is nil")
	}
	if planModel == nil {
		return nil, errors.New("plan model is nil")
	}
	var updateJobScheduleRequest *mgmtv1alpha1.UpdateJobScheduleRequest
	if j.CronSchedule.ValueString() != planModel.CronSchedule.ValueString() {
		updateJobScheduleRequest = &mgmtv1alpha1.UpdateJobScheduleRequest{
			Id:           j.Id.ValueString(),
			CronSchedule: planModel.CronSchedule.ValueStringPointer(),
		}
	}

	newSource, err := planModel.JobSource.ToDto()
	if err != nil {
		return nil, err
	}

	newMappings, err := ToJobMappingsDto(planModel.Mappings)
	if err != nil {
		return nil, err
	}

	var virtualForeignKeys []*mgmtv1alpha1.VirtualForeignConstraint
	if len(planModel.VirtualForeignKeys) > 0 {
		virtualForeignKeys = make([]*mgmtv1alpha1.VirtualForeignConstraint, 0, len(planModel.VirtualForeignKeys))
		for _, vfk := range planModel.VirtualForeignKeys {
			vfkDto, err := vfk.ToDto()
			if err != nil {
				return nil, err
			}
			virtualForeignKeys = append(virtualForeignKeys, vfkDto)
		}
	}

	// todo: compare plan and state and only conditionally create request if there are actually changes
	updateJobSourceConnectionRequest := &mgmtv1alpha1.UpdateJobSourceConnectionRequest{
		Id:                 j.Id.ValueString(),
		Source:             newSource,
		Mappings:           newMappings,
		VirtualForeignKeys: virtualForeignKeys,
	}

	destinationsToCreate := []*JobDestination{}
	destinationsToUpdate := []*JobDestination{}
	destinationsToDelete := []*JobDestination{}

	stateDestinationsMap := map[string]*JobDestination{}
	for _, dst := range j.Destinations {
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
		destinationsToUpdate = append(destinationsToUpdate, dst) // todo: should do work here to see if it has actually changed at all
	}

	var createJobDestinationConnectionsRequest *mgmtv1alpha1.CreateJobDestinationConnectionsRequest
	if len(destinationsToCreate) > 0 {
		destinationDtos := make([]*mgmtv1alpha1.CreateJobDestination, 0, len(destinationsToCreate))
		for _, dst := range destinationsToCreate {
			destinationDto, err := dst.ToCreateJobDestinationDto()
			if err != nil {
				return nil, err
			}
			destinationDtos = append(destinationDtos, destinationDto)
		}
		createJobDestinationConnectionsRequest = &mgmtv1alpha1.CreateJobDestinationConnectionsRequest{
			JobId:        j.Id.ValueString(),
			Destinations: destinationDtos,
		}
	}

	updateJobDestinationConnectionRequests := make([]*mgmtv1alpha1.UpdateJobDestinationConnectionRequest, 0, len(destinationsToUpdate))
	if len(destinationsToUpdate) > 0 {
		for _, dst := range destinationsToUpdate {
			destinationDto, err := dst.ToDto()
			if err != nil {
				return nil, err
			}
			updateJobDestinationConnectionRequests = append(updateJobDestinationConnectionRequests, &mgmtv1alpha1.UpdateJobDestinationConnectionRequest{
				JobId:         jobId,
				DestinationId: destinationDto.GetId(),
				ConnectionId:  destinationDto.GetConnectionId(),
				Options:       destinationDto.GetOptions(),
			})
		}
	}

	deleteJobDestinationConnectionRequests := make([]*mgmtv1alpha1.DeleteJobDestinationConnectionRequest, 0, len(destinationsToDelete))
	if len(destinationsToDelete) > 0 {
		for _, dst := range destinationsToDelete {
			deleteJobDestinationConnectionRequests = append(deleteJobDestinationConnectionRequests, &mgmtv1alpha1.DeleteJobDestinationConnectionRequest{
				DestinationId: dst.Id.ValueString(),
			})
		}
	}
	return &UpdateJobRequest{
		UpdateJobScheduleRequest:               updateJobScheduleRequest,
		UpdateJobSourceConnectionRequest:       updateJobSourceConnectionRequest,
		CreateJobDestinationConnectionsRequest: createJobDestinationConnectionsRequest,
		UpdateJobDestinationConnectionRequests: updateJobDestinationConnectionRequests,
		DeleteJobDestinationConnectionRequests: deleteJobDestinationConnectionRequests,
	}, nil
}

func (j *JobResourceModel) FromDto(dto *mgmtv1alpha1.Job) error {
	if j == nil {
		return errors.New("job resource model is nil")
	}
	if dto == nil {
		return errors.New("job dto is nil")
	}

	j.Id = types.StringValue(dto.Id)
	j.Name = types.StringValue(dto.Name)
	j.AccountId = types.StringValue(dto.AccountId)

	if dto.CronSchedule != nil {
		j.CronSchedule = types.StringValue(*dto.CronSchedule)
	}

	if dto.Mappings != nil {
		mappings, err := FromJobMappingsDto(dto.Mappings)
		if err != nil {
			return err
		}
		j.Mappings = mappings
	}

	if dto.Source != nil {
		source := &JobSource{}
		err := source.FromDto(dto.Source)
		if err != nil {
			return err
		}
		j.JobSource = source
	}

	if dto.Destinations != nil {
		destinations := make([]*JobDestination, 0, len(dto.Destinations))
		for _, destDto := range dto.Destinations {
			destination := &JobDestination{}
			err := destination.FromDto(destDto)
			if err != nil {
				return err
			}
			destinations = append(destinations, destination)
		}
		j.Destinations = destinations
	}

	if dto.SyncOptions != nil && (dto.SyncOptions.ScheduleToCloseTimeout != nil || dto.SyncOptions.StartToCloseTimeout != nil || dto.SyncOptions.RetryPolicy != nil) {
		syncOpts := &ActivityOptions{}
		err := syncOpts.FromDto(dto.SyncOptions)
		if err != nil {
			return err
		}
		j.SyncOptions = syncOpts
	}

	if dto.WorkflowOptions != nil && (dto.WorkflowOptions.RunTimeout != nil) {
		workflowOpts := &WorkflowOptions{}
		err := workflowOpts.FromDto(dto.WorkflowOptions)
		if err != nil {
			return err
		}
		j.WorkflowOptions = workflowOpts
	}

	if len(dto.GetVirtualForeignKeys()) > 0 {
		vkeys := make([]*VirtualForeignKeyConstraint, 0, len(dto.GetVirtualForeignKeys()))
		for _, dtoVfkey := range dto.GetVirtualForeignKeys() {
			vfk := &VirtualForeignKeyConstraint{}
			err := vfk.FromDto(dtoVfkey)
			if err != nil {
				return err
			}
			vkeys = append(vkeys, vfk)
		}
		j.VirtualForeignKeys = vkeys
	}

	return nil
}
func (v *VirtualForeignKeyConstraint) ToDto() (*mgmtv1alpha1.VirtualForeignConstraint, error) {
	if v == nil {
		return nil, errors.New("virtual foreign key constraint is nil")
	}

	var columns []string
	if len(v.Columns) > 0 {
		columns = make([]string, 0, len(v.Columns))
		for _, column := range v.Columns {
			columns = append(columns, column.ValueString())
		}
	}

	var foreignKey *mgmtv1alpha1.VirtualForeignKey
	if v.ForeignKey != nil {
		foreignKeyDto, err := v.ForeignKey.ToDto()
		if err != nil {
			return nil, err
		}
		foreignKey = foreignKeyDto
	}

	return &mgmtv1alpha1.VirtualForeignConstraint{
		Schema:     v.Schema.ValueString(),
		Table:      v.Table.ValueString(),
		Columns:    columns,
		ForeignKey: foreignKey,
	}, nil
}

func (v *VirtualForeignKeyConstraint) FromDto(dto *mgmtv1alpha1.VirtualForeignConstraint) error {
	if v == nil {
		return errors.New("virtual foreign key constraint is nil")
	}
	if dto == nil {
		return errors.New("virtual foreign key constraint dto is nil")
	}

	v.Schema = types.StringValue(dto.Schema)
	v.Table = types.StringValue(dto.Table)

	var columns []types.String
	if len(dto.Columns) > 0 {
		columns = make([]types.String, 0, len(dto.Columns))
		for _, column := range dto.Columns {
			columns = append(columns, types.StringValue(column))
		}
	}
	v.Columns = columns

	var foreignKey *VirtualForeignKey
	if dto.ForeignKey != nil {
		foreignKey = &VirtualForeignKey{}
		err := foreignKey.FromDto(dto.ForeignKey)
		if err != nil {
			return err
		}
	}

	v.ForeignKey = foreignKey
	return nil
}

func (v *VirtualForeignKey) ToDto() (*mgmtv1alpha1.VirtualForeignKey, error) {
	if v == nil {
		return nil, errors.New("virtual foreign key is nil")
	}

	var columns []string
	if len(v.Columns) > 0 {
		columns = make([]string, 0, len(v.Columns))
		for _, column := range v.Columns {
			columns = append(columns, column.ValueString())
		}
	}

	return &mgmtv1alpha1.VirtualForeignKey{
		Schema:  v.Schema.ValueString(),
		Table:   v.Table.ValueString(),
		Columns: columns,
	}, nil
}

func (v *VirtualForeignKey) FromDto(dto *mgmtv1alpha1.VirtualForeignKey) error {
	if v == nil {
		return errors.New("virtual foreign key is nil")
	}
	if dto == nil {
		return errors.New("virtual foreign key dto is nil")
	}

	v.Schema = types.StringValue(dto.Schema)
	v.Table = types.StringValue(dto.Table)

	var columns []types.String
	if len(dto.Columns) > 0 {
		columns = make([]types.String, 0, len(dto.Columns))
		for _, column := range dto.Columns {
			columns = append(columns, types.StringValue(column))
		}
	}
	v.Columns = columns

	return nil
}

func (a *ActivityOptions) ToDto() (*mgmtv1alpha1.ActivityOptions, error) {
	if a == nil {
		return nil, errors.New("activity options is nil")
	}

	var retryPolicy *mgmtv1alpha1.RetryPolicy
	if a.RetryPolicy != nil {
		retryPolicy = &mgmtv1alpha1.RetryPolicy{
			MaximumAttempts: i64Toi32(a.RetryPolicy.MaximumAttempts.ValueInt64Pointer()),
		}
	}

	return &mgmtv1alpha1.ActivityOptions{
		ScheduleToCloseTimeout: a.ScheduleToCloseTimeout.ValueInt64Pointer(),
		StartToCloseTimeout:    a.StartToCloseTimeout.ValueInt64Pointer(),
		RetryPolicy:            retryPolicy,
	}, nil
}

func (a *ActivityOptions) FromDto(dto *mgmtv1alpha1.ActivityOptions) error {
	if a == nil {
		return errors.New("activity options is nil")
	}
	if dto == nil {
		return errors.New("activity options dto is nil")
	}

	a.ScheduleToCloseTimeout = types.Int64PointerValue(dto.ScheduleToCloseTimeout)
	a.StartToCloseTimeout = types.Int64PointerValue(dto.StartToCloseTimeout)

	if dto.RetryPolicy != nil {
		a.RetryPolicy = &RetryPolicy{
			MaximumAttempts: types.Int64Value(int64(*dto.RetryPolicy.MaximumAttempts)),
		}
	}
	return nil
}

func (w *WorkflowOptions) ToDto() (*mgmtv1alpha1.WorkflowOptions, error) {
	if w == nil {
		return nil, errors.New("workflow options is nil")
	}

	return &mgmtv1alpha1.WorkflowOptions{
		RunTimeout: w.RunTimeout.ValueInt64Pointer(),
	}, nil
}

func (w *WorkflowOptions) FromDto(dto *mgmtv1alpha1.WorkflowOptions) error {
	if w == nil {
		return errors.New("workflow options is nil")
	}
	if dto == nil {
		return errors.New("workflow options dto is nil")
	}

	w.RunTimeout = types.Int64PointerValue(dto.RunTimeout)
	return nil
}

func ToJobMappingsDto(mappings []*JobMapping) ([]*mgmtv1alpha1.JobMapping, error) {
	if mappings == nil {
		return nil, errors.New("mappings is nil")
	}

	dtos := make([]*mgmtv1alpha1.JobMapping, 0, len(mappings))
	for _, mapping := range mappings {
		dto, err := mapping.ToDto()
		if err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

func FromJobMappingsDto(dtos []*mgmtv1alpha1.JobMapping) ([]*JobMapping, error) {
	output := make([]*JobMapping, 0, len(dtos))
	for _, dto := range dtos {
		mapping := &JobMapping{}
		err := mapping.FromDto(dto)
		if err != nil {
			return nil, err
		}
		output = append(output, mapping)
	}
	return output, nil
}

func (j *JobMapping) ToDto() (*mgmtv1alpha1.JobMapping, error) {
	if j == nil {
		return nil, errors.New("job mapping is nil")
	}

	transformerDto, err := j.Transformer.ToDto()
	if err != nil {
		return nil, err
	}

	return &mgmtv1alpha1.JobMapping{
		Schema: j.Schema.ValueString(),
		Table:  j.Table.ValueString(),
		Column: j.Column.ValueString(),
		Transformer: &mgmtv1alpha1.JobMappingTransformer{
			Config: transformerDto,
		},
	}, nil
}

func (j *JobMapping) FromDto(dto *mgmtv1alpha1.JobMapping) error {
	if j == nil {
		return errors.New("job mapping is nil")
	}
	if dto == nil {
		return errors.New("job mapping dto is nil")
	}

	tconfig := &transformer_model.Transformer{}
	err := tconfig.FromDto(dto.Transformer.Config)
	if err != nil {
		return err
	}
	j.Transformer = tconfig
	j.Schema = types.StringValue(dto.Schema)
	j.Table = types.StringValue(dto.Table)
	j.Column = types.StringValue(dto.Column)

	return nil
}

func (j *JobDestination) ToCreateJobDestinationDto() (*mgmtv1alpha1.CreateJobDestination, error) {
	if j == nil {
		return nil, errors.New("job destination is nil")
	}

	jd := &mgmtv1alpha1.CreateJobDestination{
		ConnectionId: j.ConnectionId.ValueString(),
	}

	if j.Postgres != nil {
		pgDto, err := j.Postgres.ToDto()
		if err != nil {
			return nil, err
		}
		jd.Options = &mgmtv1alpha1.JobDestinationOptions{
			Config: pgDto,
		}
	}
	if j.Mysql != nil {
		mysqlDto, err := j.Mysql.ToDto()
		if err != nil {
			return nil, err
		}
		jd.Options = &mgmtv1alpha1.JobDestinationOptions{
			Config: mysqlDto,
		}
	}
	if j.AwsS3 != nil {
		awsS3Dto, err := j.AwsS3.ToDto()
		if err != nil {
			return nil, err
		}
		jd.Options = &mgmtv1alpha1.JobDestinationOptions{
			Config: awsS3Dto,
		}
	}
	return jd, nil
}

func (j *JobDestination) ToDto() (*mgmtv1alpha1.JobDestination, error) {
	if j == nil {
		return nil, errors.New("job destination is nil")
	}

	jd := &mgmtv1alpha1.JobDestination{
		Id:           j.Id.ValueString(),
		ConnectionId: j.ConnectionId.ValueString(),
	}

	if j.Postgres != nil {
		pgDto, err := j.Postgres.ToDto()
		if err != nil {
			return nil, err
		}
		jd.Options = &mgmtv1alpha1.JobDestinationOptions{
			Config: pgDto,
		}
	}
	if j.Mysql != nil {
		mysqlDto, err := j.Mysql.ToDto()
		if err != nil {
			return nil, err
		}
		jd.Options = &mgmtv1alpha1.JobDestinationOptions{
			Config: mysqlDto,
		}
	}
	if j.AwsS3 != nil {
		awsS3Dto, err := j.AwsS3.ToDto()
		if err != nil {
			return nil, err
		}
		jd.Options = &mgmtv1alpha1.JobDestinationOptions{
			Config: awsS3Dto,
		}
	}
	return jd, nil
}

func (j *JobDestination) FromDto(dto *mgmtv1alpha1.JobDestination) error {
	if j == nil {
		return errors.New("job destination is nil")
	}
	if dto == nil {
		return errors.New("job destination dto is nil")
	}

	j.Id = types.StringValue(dto.Id)
	j.ConnectionId = types.StringValue(dto.ConnectionId)

	switch config := dto.GetOptions().GetConfig().(type) {
	case *mgmtv1alpha1.JobDestinationOptions_PostgresOptions:
		j.Postgres = &JobDestinationPostgresOptions{}
		err := j.Postgres.FromDto(config.PostgresOptions)
		if err != nil {
			return err
		}
	case *mgmtv1alpha1.JobDestinationOptions_MysqlOptions:
		j.Mysql = &JobDestinationMysqlOptions{}
		err := j.Mysql.FromDto(config.MysqlOptions)
		if err != nil {
			return err
		}
	case *mgmtv1alpha1.JobDestinationOptions_AwsS3Options:
		j.AwsS3 = &JobDestinationAwsS3Options{}
		err := j.AwsS3.FromDto(config.AwsS3Options)
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *JobDestinationPostgresOptions) ToDto() (*mgmtv1alpha1.JobDestinationOptions_PostgresOptions, error) {
	if j == nil {
		return nil, errors.New("job destination postgres options is nil")
	}

	var truncateTable *mgmtv1alpha1.PostgresTruncateTableConfig
	if j.TruncateTable != nil {
		truncateTable = &mgmtv1alpha1.PostgresTruncateTableConfig{
			TruncateBeforeInsert: j.TruncateTable.TruncateBeforeInsert.ValueBool(),
			Cascade:              j.TruncateTable.Cascade.ValueBool(),
		}
	}

	dto := &mgmtv1alpha1.JobDestinationOptions_PostgresOptions{
		PostgresOptions: &mgmtv1alpha1.PostgresDestinationConnectionOptions{
			TruncateTable:            truncateTable,
			InitTableSchema:          j.InitTableSchema.ValueBool(),
			MaxInFlight:              nil,   // todo
			OnConflict:               nil,   // todo
			SkipForeignKeyViolations: false, // todo
			Batch:                    nil,   // todo
		},
	}

	return dto, nil
}

func (j *JobDestinationPostgresOptions) FromDto(dto *mgmtv1alpha1.PostgresDestinationConnectionOptions) error {
	if j == nil {
		return errors.New("job destination postgres options is nil")
	}
	if dto == nil {
		return errors.New("job destination postgres options dto is nil")
	}

	if dto.TruncateTable != nil {
		j.TruncateTable = &PostgresDestinationTruncateTable{
			TruncateBeforeInsert: types.BoolValue(dto.TruncateTable.TruncateBeforeInsert),
			Cascade:              types.BoolValue(dto.TruncateTable.Cascade),
		}
	}

	j.InitTableSchema = types.BoolValue(dto.InitTableSchema)

	return nil
}

func (j *JobDestinationMysqlOptions) ToDto() (*mgmtv1alpha1.JobDestinationOptions_MysqlOptions, error) {
	if j == nil {
		return nil, errors.New("job destination mysql options is nil")
	}

	var truncateTable *mgmtv1alpha1.MysqlTruncateTableConfig
	if j.TruncateTable != nil {
		truncateTable = &mgmtv1alpha1.MysqlTruncateTableConfig{
			TruncateBeforeInsert: j.TruncateTable.TruncateBeforeInsert.ValueBool(),
		}
	}

	dto := &mgmtv1alpha1.JobDestinationOptions_MysqlOptions{
		MysqlOptions: &mgmtv1alpha1.MysqlDestinationConnectionOptions{
			TruncateTable:            truncateTable,
			InitTableSchema:          j.InitTableSchema.ValueBool(),
			SkipForeignKeyViolations: false, // todo
			MaxInFlight:              nil,   // todo
			OnConflict:               nil,   // todo
			Batch:                    nil,   // todo
		},
	}

	return dto, nil
}

func (j *JobDestinationMysqlOptions) FromDto(dto *mgmtv1alpha1.MysqlDestinationConnectionOptions) error {
	if j == nil {
		return errors.New("job destination mysql options is nil")
	}
	if dto == nil {
		return errors.New("job destination mysql options dto is nil")
	}

	if dto.TruncateTable != nil {
		j.TruncateTable = &MysqlDestinationTruncateTable{
			TruncateBeforeInsert: types.BoolValue(dto.TruncateTable.TruncateBeforeInsert),
		}
	}

	j.InitTableSchema = types.BoolValue(dto.InitTableSchema)
	return nil
}

func (j *JobDestinationAwsS3Options) ToDto() (*mgmtv1alpha1.JobDestinationOptions_AwsS3Options, error) {
	if j == nil {
		return nil, errors.New("job destination aws s3 options is nil")
	}

	dto := &mgmtv1alpha1.JobDestinationOptions_AwsS3Options{
		AwsS3Options: &mgmtv1alpha1.AwsS3DestinationConnectionOptions{
			// todo
		},
	}

	return dto, nil
}

func (j *JobDestinationAwsS3Options) FromDto(dto *mgmtv1alpha1.AwsS3DestinationConnectionOptions) error {
	if j == nil {
		return errors.New("job destination aws s3 options is nil")
	}
	if dto == nil {
		return errors.New("job destination aws s3 options dto is nil")
	}

	// todo

	return nil
}

func (j *JobSource) ToDto() (*mgmtv1alpha1.JobSource, error) {
	if j == nil {
		return nil, errors.New("job source is nil")
	}

	if j.Postgres != nil {
		pgDto, err := j.Postgres.ToDto()
		if err != nil {
			return nil, err
		}
		return &mgmtv1alpha1.JobSource{
			Options: &mgmtv1alpha1.JobSourceOptions{
				Config: pgDto,
			},
		}, nil
	}
	if j.Mysql != nil {
		mysqlDto, err := j.Mysql.ToDto()
		if err != nil {
			return nil, err
		}
		return &mgmtv1alpha1.JobSource{
			Options: &mgmtv1alpha1.JobSourceOptions{
				Config: mysqlDto,
			},
		}, nil
	}
	if j.Generate != nil {
		generateDto, err := j.Generate.ToDto()
		if err != nil {
			return nil, err
		}
		return &mgmtv1alpha1.JobSource{
			Options: &mgmtv1alpha1.JobSourceOptions{
				Config: generateDto,
			},
		}, nil
	}
	if j.AwsS3 != nil {
		awsS3Dto, err := j.AwsS3.ToDto()
		if err != nil {
			return nil, err
		}
		return &mgmtv1alpha1.JobSource{
			Options: &mgmtv1alpha1.JobSourceOptions{
				Config: awsS3Dto,
			},
		}, nil
	}

	return nil, nil
}

func (j *JobSourcePostgresOptions) ToDto() (*mgmtv1alpha1.JobSourceOptions_Postgres, error) {
	if j == nil {
		return nil, errors.New("job source postgres options is nil")
	}

	schemas := make([]*mgmtv1alpha1.PostgresSourceSchemaOption, 0, len(j.SchemaOptions))
	for _, schema := range j.SchemaOptions {
		schemaDto, err := schema.ToDto()
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, schemaDto)
	}

	dto := &mgmtv1alpha1.JobSourceOptions_Postgres{
		Postgres: &mgmtv1alpha1.PostgresSourceConnectionOptions{
			ConnectionId:                  j.ConnectionId.ValueString(),
			SubsetByForeignKeyConstraints: j.SubsetByForeignKeyConstraints.ValueBool(),
			NewColumnAdditionStrategy:     nil, // set below
			ColumnRemovalStrategy:         nil, // set below
			Schemas:                       schemas,
		},
	}

	if j.NewColumnAdditionStrategy != nil {
		strategyDto, err := j.NewColumnAdditionStrategy.ToDto()
		if err != nil {
			return nil, err
		}
		dto.Postgres.NewColumnAdditionStrategy = strategyDto
	} else if !j.SubsetByForeignKeyConstraints.IsNull() {
		// Handle deprecated field
		if j.SubsetByForeignKeyConstraints.ValueBool() {
			dto.Postgres.NewColumnAdditionStrategy = &mgmtv1alpha1.PostgresSourceConnectionOptions_NewColumnAdditionStrategy{
				Strategy: &mgmtv1alpha1.PostgresSourceConnectionOptions_NewColumnAdditionStrategy_HaltJob_{
					HaltJob: &mgmtv1alpha1.PostgresSourceConnectionOptions_NewColumnAdditionStrategy_HaltJob{},
				},
			}
		}
	}

	if j.ColumnRemovalStrategy != nil {
		strategyDto, err := j.ColumnRemovalStrategy.ToDto()
		if err != nil {
			return nil, err
		}
		dto.Postgres.ColumnRemovalStrategy = strategyDto
	}

	return dto, nil
}

func (j *JobSourcePostgresOptions) FromDto(dto *mgmtv1alpha1.PostgresSourceConnectionOptions) error {
	if j == nil {
		return errors.New("job source postgres options is nil")
	}
	if dto == nil {
		return errors.New("job source postgres options dto is nil")
	}

	j.ConnectionId = types.StringValue(dto.ConnectionId)
	j.SubsetByForeignKeyConstraints = types.BoolValue(dto.SubsetByForeignKeyConstraints)
	if len(dto.Schemas) > 0 {
		j.SchemaOptions = make([]*JobSourcePostgresSourceSchemaOption, 0, len(dto.Schemas))
		for _, schemaDto := range dto.Schemas {
			schema := &JobSourcePostgresSourceSchemaOption{}
			err := schema.FromDto(schemaDto)
			if err != nil {
				return err
			}
			j.SchemaOptions = append(j.SchemaOptions, schema)
		}
	}
	j.SubsetByForeignKeyConstraints = types.BoolValue(dto.SubsetByForeignKeyConstraints)

	if dto.NewColumnAdditionStrategy != nil {
		j.NewColumnAdditionStrategy = &PostgresNewColumnAdditionStrategy{}
		err := j.NewColumnAdditionStrategy.FromDto(dto.NewColumnAdditionStrategy)
		if err != nil {
			return err
		}
	}

	if dto.ColumnRemovalStrategy != nil {
		j.ColumnRemovalStrategy = &PostgresColumnRemovalStrategy{}
		err := j.ColumnRemovalStrategy.FromDto(dto.ColumnRemovalStrategy)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *PostgresNewColumnAdditionStrategy) ToDto() (*mgmtv1alpha1.PostgresSourceConnectionOptions_NewColumnAdditionStrategy, error) {
	if j == nil {
		return nil, errors.New("postgres new column addition strategy is nil")
	}
	if j.HaltJob != nil {
		return &mgmtv1alpha1.PostgresSourceConnectionOptions_NewColumnAdditionStrategy{
			Strategy: &mgmtv1alpha1.PostgresSourceConnectionOptions_NewColumnAdditionStrategy_HaltJob_{
				HaltJob: &mgmtv1alpha1.PostgresSourceConnectionOptions_NewColumnAdditionStrategy_HaltJob{},
			},
		}, nil
	}
	if j.AutoMap != nil {
		return &mgmtv1alpha1.PostgresSourceConnectionOptions_NewColumnAdditionStrategy{
			Strategy: &mgmtv1alpha1.PostgresSourceConnectionOptions_NewColumnAdditionStrategy_AutoMap_{
				AutoMap: &mgmtv1alpha1.PostgresSourceConnectionOptions_NewColumnAdditionStrategy_AutoMap{},
			},
		}, nil
	}
	return nil, nil
}

func (j *PostgresNewColumnAdditionStrategy) FromDto(dto *mgmtv1alpha1.PostgresSourceConnectionOptions_NewColumnAdditionStrategy) error {
	if j == nil {
		return errors.New("postgres new column addition strategy is nil")
	}
	if dto == nil {
		return errors.New("postgres new column addition strategy dto is nil")
	}

	switch dto.GetStrategy().(type) {
	case *mgmtv1alpha1.PostgresSourceConnectionOptions_NewColumnAdditionStrategy_HaltJob_:
		j.HaltJob = &PostgresNewColumnAdditionStrategyHaltJob{}
	case *mgmtv1alpha1.PostgresSourceConnectionOptions_NewColumnAdditionStrategy_AutoMap_:
		j.AutoMap = &PostgresNewColumnAdditionStrategyAutoMap{}
	}

	return nil
}

func (j *PostgresColumnRemovalStrategy) ToDto() (*mgmtv1alpha1.PostgresSourceConnectionOptions_ColumnRemovalStrategy, error) {
	if j == nil {
		return nil, errors.New("postgres column removal strategy is nil")
	}
	if j.HaltJob != nil {
		return &mgmtv1alpha1.PostgresSourceConnectionOptions_ColumnRemovalStrategy{
			Strategy: &mgmtv1alpha1.PostgresSourceConnectionOptions_ColumnRemovalStrategy_HaltJob_{
				HaltJob: &mgmtv1alpha1.PostgresSourceConnectionOptions_ColumnRemovalStrategy_HaltJob{},
			},
		}, nil
	}
	if j.ContinueJob != nil {
		return &mgmtv1alpha1.PostgresSourceConnectionOptions_ColumnRemovalStrategy{
			Strategy: &mgmtv1alpha1.PostgresSourceConnectionOptions_ColumnRemovalStrategy_ContinueJob_{
				ContinueJob: &mgmtv1alpha1.PostgresSourceConnectionOptions_ColumnRemovalStrategy_ContinueJob{},
			},
		}, nil
	}
	return nil, nil
}

func (j *PostgresColumnRemovalStrategy) FromDto(dto *mgmtv1alpha1.PostgresSourceConnectionOptions_ColumnRemovalStrategy) error {
	if j == nil {
		return errors.New("postgres column removal strategy is nil")
	}
	if dto == nil {
		return errors.New("postgres column removal strategy dto is nil")
	}

	switch dto.GetStrategy().(type) {
	case *mgmtv1alpha1.PostgresSourceConnectionOptions_ColumnRemovalStrategy_HaltJob_:
		j.HaltJob = &PostgresHaltJobColumnRemovalStrategy{}
	case *mgmtv1alpha1.PostgresSourceConnectionOptions_ColumnRemovalStrategy_ContinueJob_:
		j.ContinueJob = &PostgresContinueJobColumnRemovalStrategy{}
	}

	return nil
}

func (j *JobSourcePostgresSourceSchemaOption) ToDto() (*mgmtv1alpha1.PostgresSourceSchemaOption, error) {
	if j == nil {
		return nil, errors.New("job source postgres source schema option is nil")
	}

	dto := &mgmtv1alpha1.PostgresSourceSchemaOption{
		Schema: j.Schema.ValueString(),
		Tables: make([]*mgmtv1alpha1.PostgresSourceTableOption, 0, len(j.Tables)),
	}

	for _, table := range j.Tables {
		tableDto, err := table.ToDto()
		if err != nil {
			return nil, err
		}
		dto.Tables = append(dto.Tables, tableDto)
	}
	return dto, nil
}

func (j *JobSourcePostgresSourceSchemaOption) FromDto(dto *mgmtv1alpha1.PostgresSourceSchemaOption) error {
	if j == nil {
		return errors.New("job source postgres source schema option is nil")
	}
	if dto == nil {
		return errors.New("job source postgres source schema option dto is nil")
	}

	j.Schema = types.StringValue(dto.Schema)
	j.Tables = make([]*JobSourcePostgresSourceTableOption, 0, len(dto.Tables))
	for _, dtoTable := range dto.Tables {
		table := &JobSourcePostgresSourceTableOption{}
		err := table.FromDto(dtoTable)
		if err != nil {
			return err
		}
		j.Tables = append(j.Tables, table)
	}

	return nil
}

func (j *JobSourcePostgresSourceTableOption) ToDto() (*mgmtv1alpha1.PostgresSourceTableOption, error) {
	if j == nil {
		return nil, errors.New("job source postgres source table option is nil")
	}

	dto := &mgmtv1alpha1.PostgresSourceTableOption{
		Table:       j.Table.ValueString(),
		WhereClause: j.WhereClause.ValueStringPointer(),
	}

	return dto, nil
}

func (j *JobSourcePostgresSourceTableOption) FromDto(dto *mgmtv1alpha1.PostgresSourceTableOption) error {
	if j == nil {
		return errors.New("job source postgres source table option is nil")
	}
	if dto == nil {
		return errors.New("job source postgres source table option dto is nil")
	}

	j.Table = types.StringValue(dto.Table)
	j.WhereClause = types.StringPointerValue(dto.WhereClause)

	return nil
}

func (j *JobSource) FromDto(dto *mgmtv1alpha1.JobSource) error {
	if j == nil {
		return errors.New("job source is nil")
	}
	if dto == nil {
		return errors.New("job source dto is nil")
	}

	switch source := dto.Options.Config.(type) {
	case *mgmtv1alpha1.JobSourceOptions_Postgres:
		j.Postgres = &JobSourcePostgresOptions{}
		err := j.Postgres.FromDto(source.Postgres)
		if err != nil {
			return err
		}
	case *mgmtv1alpha1.JobSourceOptions_Mysql:
		j.Mysql = &JobSourceMysqlOptions{}
		err := j.Mysql.FromDto(source.Mysql)
		if err != nil {
			return err
		}
	case *mgmtv1alpha1.JobSourceOptions_Generate:
		j.Generate = &JobSourceGenerateOptions{}
		err := j.Generate.FromDto(source.Generate)
		if err != nil {
			return err
		}
	case *mgmtv1alpha1.JobSourceOptions_AwsS3:
		j.AwsS3 = &JobSourceAwsS3Options{}
		err := j.AwsS3.FromDto(source.AwsS3)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *JobSourceMysqlSourceSchemaOption) ToDto() (*mgmtv1alpha1.MysqlSourceSchemaOption, error) {
	if j == nil {
		return nil, errors.New("job source mysql source schema option is nil")
	}

	dto := &mgmtv1alpha1.MysqlSourceSchemaOption{
		Schema: j.Schema.ValueString(),
		Tables: make([]*mgmtv1alpha1.MysqlSourceTableOption, 0, len(j.Tables)),
	}

	for _, table := range j.Tables {
		tableDto, err := table.ToDto()
		if err != nil {
			return nil, err
		}
		dto.Tables = append(dto.Tables, tableDto)
	}
	return dto, nil
}

func (j *JobSourceMysqlSourceSchemaOption) FromDto(dto *mgmtv1alpha1.MysqlSourceSchemaOption) error {
	if j == nil {
		return errors.New("job source mysql source schema option is nil")
	}
	if dto == nil {
		return errors.New("job source mysql source schema option dto is nil")
	}

	j.Schema = types.StringValue(dto.Schema)
	j.Tables = make([]*JobSourceMysqlSourceTableOption, 0, len(dto.Tables))
	for _, dtoTable := range dto.Tables {
		table := &JobSourceMysqlSourceTableOption{}
		err := table.FromDto(dtoTable)
		if err != nil {
			return err
		}
		j.Tables = append(j.Tables, table)
	}

	return nil
}

func (j *JobSourceMysqlSourceTableOption) ToDto() (*mgmtv1alpha1.MysqlSourceTableOption, error) {
	if j == nil {
		return nil, errors.New("job source mysql source table option is nil")
	}

	dto := &mgmtv1alpha1.MysqlSourceTableOption{
		Table:       j.Table.ValueString(),
		WhereClause: j.WhereClause.ValueStringPointer(),
	}

	return dto, nil
}

func (j *JobSourceMysqlSourceTableOption) FromDto(dto *mgmtv1alpha1.MysqlSourceTableOption) error {
	if j == nil {
		return errors.New("job source mysql source table option is nil")
	}
	if dto == nil {
		return errors.New("job source mysql source table option dto is nil")
	}

	j.Table = types.StringValue(dto.Table)
	j.WhereClause = types.StringPointerValue(dto.WhereClause)

	return nil
}
func (j *JobSourceGenerateOptions) ToDto() (*mgmtv1alpha1.JobSourceOptions_Generate, error) {
	if j == nil {
		return nil, errors.New("job source generate options is nil")
	}

	var schemas []*mgmtv1alpha1.GenerateSourceSchemaOption
	if len(j.Schemas) > 0 {
		schemas = make([]*mgmtv1alpha1.GenerateSourceSchemaOption, 0, len(j.Schemas))
		for _, schema := range j.Schemas {
			schemaDto, err := schema.ToDto()
			if err != nil {
				return nil, err
			}
			schemas = append(schemas, schemaDto)
		}
	}

	return &mgmtv1alpha1.JobSourceOptions_Generate{
		Generate: &mgmtv1alpha1.GenerateSourceOptions{
			Schemas:              schemas,
			FkSourceConnectionId: j.FkSourceConnectionId.ValueStringPointer(),
		},
	}, nil
}

func (j *JobSourceGenerateOptions) FromDto(dto *mgmtv1alpha1.GenerateSourceOptions) error {
	if j == nil {
		return errors.New("job source generate options is nil")
	}
	if dto == nil {
		return errors.New("job source generate options dto is nil")
	}

	j.FkSourceConnectionId = types.StringPointerValue(dto.FkSourceConnectionId)

	if len(dto.Schemas) > 0 {
		j.Schemas = make([]*JobSourceGenerateSchemaOption, 0, len(dto.Schemas))
		for _, schemaDto := range dto.Schemas {
			schema := &JobSourceGenerateSchemaOption{}
			err := schema.FromDto(schemaDto)
			if err != nil {
				return err
			}
			j.Schemas = append(j.Schemas, schema)
		}
	}

	return nil
}

func (j *JobSourceGenerateSchemaOption) ToDto() (*mgmtv1alpha1.GenerateSourceSchemaOption, error) {
	if j == nil {
		return nil, errors.New("job source generate schema option is nil")
	}

	var tables []*mgmtv1alpha1.GenerateSourceTableOption
	if len(j.Tables) > 0 {
		tables = make([]*mgmtv1alpha1.GenerateSourceTableOption, 0, len(j.Tables))
		for _, table := range j.Tables {
			tableDto, err := table.ToDto()
			if err != nil {
				return nil, err
			}
			tables = append(tables, tableDto)
		}
	}

	return &mgmtv1alpha1.GenerateSourceSchemaOption{
		Schema: j.Schema.ValueString(),
		Tables: tables,
	}, nil
}

func (j *JobSourceGenerateSchemaOption) FromDto(dto *mgmtv1alpha1.GenerateSourceSchemaOption) error {
	if j == nil {
		return errors.New("job source generate schema option is nil")
	}
	if dto == nil {
		return errors.New("job source generate schema option dto is nil")
	}

	j.Schema = types.StringValue(dto.Schema)
	if len(dto.Tables) > 0 {
		j.Tables = make([]*JobSourceGenerateTableOption, 0, len(dto.Tables))
		for _, dtoTable := range dto.Tables {
			table := &JobSourceGenerateTableOption{}
			err := table.FromDto(dtoTable)
			if err != nil {
				return err
			}
			j.Tables = append(j.Tables, table)
		}
	}
	return nil
}

func (j *JobSourceGenerateTableOption) ToDto() (*mgmtv1alpha1.GenerateSourceTableOption, error) {
	if j == nil {
		return nil, errors.New("job source generate table option is nil")
	}

	return &mgmtv1alpha1.GenerateSourceTableOption{
		Table:    j.Table.ValueString(),
		RowCount: j.RowCount.ValueInt64(),
	}, nil
}

func (j *JobSourceGenerateTableOption) FromDto(dto *mgmtv1alpha1.GenerateSourceTableOption) error {
	if j == nil {
		return errors.New("job source generate table option is nil")
	}
	if dto == nil {
		return errors.New("job source generate table option dto is nil")
	}

	j.Table = types.StringValue(dto.Table)
	j.RowCount = types.Int64Value(dto.RowCount)

	return nil
}

func (j *JobSourceAwsS3Options) ToDto() (*mgmtv1alpha1.JobSourceOptions_AwsS3, error) {
	if j == nil {
		return nil, errors.New("job source aws s3 options is nil")
	}

	return &mgmtv1alpha1.JobSourceOptions_AwsS3{
		AwsS3: &mgmtv1alpha1.AwsS3SourceConnectionOptions{
			ConnectionId: j.ConnectionId.ValueString(),
		},
	}, nil
}

func (j *JobSourceAwsS3Options) FromDto(dto *mgmtv1alpha1.AwsS3SourceConnectionOptions) error {
	if j == nil {
		return errors.New("job source aws s3 options is nil")
	}
	if dto == nil {
		return errors.New("job source aws s3 options dto is nil")
	}

	j.ConnectionId = types.StringValue(dto.ConnectionId)

	return nil
}

func (j *JobSourceMysqlOptions) ToDto() (*mgmtv1alpha1.JobSourceOptions_Mysql, error) {
	if j == nil {
		return nil, errors.New("job source mysql options is nil")
	}

	var schemas []*mgmtv1alpha1.MysqlSourceSchemaOption
	if len(j.SchemaOptions) > 0 {
		schemas = make([]*mgmtv1alpha1.MysqlSourceSchemaOption, 0, len(j.SchemaOptions))
		for _, schema := range j.SchemaOptions {
			schemaDto, err := schema.ToDto()
			if err != nil {
				return nil, err
			}
			schemas = append(schemas, schemaDto)
		}
	}

	dto := &mgmtv1alpha1.JobSourceOptions_Mysql{
		Mysql: &mgmtv1alpha1.MysqlSourceConnectionOptions{
			ConnectionId:                  j.ConnectionId.ValueString(),
			Schemas:                       schemas,
			SubsetByForeignKeyConstraints: j.SubsetByForeignKeyConstraints.ValueBool(),
			ColumnRemovalStrategy:         nil, // set below
		},
	}

	if j.ColumnRemovalStrategy != nil {
		strategyDto, err := j.ColumnRemovalStrategy.ToDto()
		if err != nil {
			return nil, err
		}
		dto.Mysql.ColumnRemovalStrategy = strategyDto
	}

	return dto, nil
}

func (j *JobSourceMysqlOptions) FromDto(dto *mgmtv1alpha1.MysqlSourceConnectionOptions) error {
	if j == nil {
		return errors.New("job source mysql options is nil")
	}
	if dto == nil {
		return errors.New("job source mysql options dto is nil")
	}

	j.ConnectionId = types.StringValue(dto.ConnectionId)
	if len(dto.Schemas) > 0 {
		j.SchemaOptions = make([]*JobSourceMysqlSourceSchemaOption, 0, len(dto.Schemas))
		for _, schemaDto := range dto.Schemas {
			schema := &JobSourceMysqlSourceSchemaOption{}
			err := schema.FromDto(schemaDto)
			if err != nil {
				return err
			}
			j.SchemaOptions = append(j.SchemaOptions, schema)
		}
	}
	j.SubsetByForeignKeyConstraints = types.BoolValue(dto.SubsetByForeignKeyConstraints)

	if dto.ColumnRemovalStrategy != nil {
		strategy := &MssqlColumnRemovalStrategy{}
		err := strategy.FromDto(dto.ColumnRemovalStrategy)
		if err != nil {
			return err
		}
		j.ColumnRemovalStrategy = strategy
	}

	return nil
}

func (j *MssqlColumnRemovalStrategy) ToDto() (*mgmtv1alpha1.MysqlSourceConnectionOptions_ColumnRemovalStrategy, error) {
	if j == nil {
		return nil, errors.New("mysql column removal strategy is nil")
	}
	if j.HaltJob != nil {
		return &mgmtv1alpha1.MysqlSourceConnectionOptions_ColumnRemovalStrategy{
			Strategy: &mgmtv1alpha1.MysqlSourceConnectionOptions_ColumnRemovalStrategy_HaltJob_{
				HaltJob: &mgmtv1alpha1.MysqlSourceConnectionOptions_ColumnRemovalStrategy_HaltJob{},
			},
		}, nil
	}
	if j.ContinueJob != nil {
		return &mgmtv1alpha1.MysqlSourceConnectionOptions_ColumnRemovalStrategy{
			Strategy: &mgmtv1alpha1.MysqlSourceConnectionOptions_ColumnRemovalStrategy_ContinueJob_{
				ContinueJob: &mgmtv1alpha1.MysqlSourceConnectionOptions_ColumnRemovalStrategy_ContinueJob{},
			},
		}, nil
	}
	return nil, nil
}

func (j *MssqlColumnRemovalStrategy) FromDto(dto *mgmtv1alpha1.MysqlSourceConnectionOptions_ColumnRemovalStrategy) error {
	if j == nil {
		return errors.New("mysql column removal strategy is nil")
	}
	if dto == nil {
		return errors.New("mysql column removal strategy dto is nil")
	}

	switch dto.GetStrategy().(type) {
	case *mgmtv1alpha1.MysqlSourceConnectionOptions_ColumnRemovalStrategy_HaltJob_:
		j.HaltJob = &MssqlHaltJobColumnRemovalStrategy{}
	case *mgmtv1alpha1.MysqlSourceConnectionOptions_ColumnRemovalStrategy_ContinueJob_:
		j.ContinueJob = &MssqlContinueJobColumnRemovalStrategy{}
	}

	return nil
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
