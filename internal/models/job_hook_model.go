package models

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgmtv1alpha1 "github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1"
)

type JobHookResourceModel struct {
	Id          types.String  `tfsdk:"id"`
	Name        types.String  `tfsdk:"name"`
	Description types.String  `tfsdk:"description"`
	JobId       types.String  `tfsdk:"job_id"`
	Enabled     types.Bool    `tfsdk:"enabled"`
	Priority    types.Int32   `tfsdk:"priority"`
	Config      JobHookConfig `tfsdk:"config"`
}

type JobHookConfig struct {
	// oneof
	Sql *JobHookSql `tfsdk:"sql"`
}

type JobHookSql struct {
	Query        types.String     `tfsdk:"query"`
	ConnectionId types.String     `tfsdk:"connection_id"`
	Timing       JobHookSqlTiming `tfsdk:"timing"`
}

type JobHookSqlTiming struct {
	// oneof
	PreSync  *JobHookTimingPreSync  `tfsdk:"pre_sync"`
	PostSync *JobHookTimingPostSync `tfsdk:"post_sync"`
}

type JobHookTimingPreSync struct{}
type JobHookTimingPostSync struct{}

func (j *JobHookResourceModel) ToCreateJobHookDto() (*mgmtv1alpha1.CreateJobHookRequest, error) {
	if j == nil {
		return nil, errors.New("JobHookResourceModel was nil")
	}

	config, err := j.Config.ToDto()
	if err != nil {
		return nil, err
	}

	priority := j.Priority.ValueInt32()
	if priority < 0 {
		priority = 0
	}

	return &mgmtv1alpha1.CreateJobHookRequest{
		JobId: j.JobId.ValueString(),
		Hook: &mgmtv1alpha1.NewJobHook{
			Name:        j.Name.ValueString(),
			Description: j.Description.ValueString(),
			Enabled:     j.Enabled.ValueBool(),
			Priority:    uint32(priority),
			Config:      config,
		},
	}, nil
}

func (j *JobHookResourceModel) ToUpdateJobHookDto() (*mgmtv1alpha1.UpdateJobHookRequest, error) {
	if j == nil {
		return nil, errors.New("JobHookResourceModel was nil")
	}

	config, err := j.Config.ToDto()
	if err != nil {
		return nil, err
	}

	priority := j.Priority.ValueInt32()
	if priority < 0 {
		priority = 0
	}

	return &mgmtv1alpha1.UpdateJobHookRequest{
		Id:          j.Id.ValueString(),
		Name:        j.Name.ValueString(),
		Description: j.Description.ValueString(),
		Enabled:     j.Enabled.ValueBool(),
		Priority:    uint32(priority),
		Config:      config,
	}, nil
}

func (j *JobHookResourceModel) FromDto(dto *mgmtv1alpha1.JobHook) error {
	if j == nil {
		return errors.New("JobHookResourceModel is nil, must call FromDto on non-nil model struct")
	}

	if dto == nil {
		return errors.New("job hook dto was nil")
	}

	j.Id = types.StringValue(dto.GetId())
	j.JobId = types.StringValue(dto.GetJobId())
	j.Name = types.StringValue(dto.GetName())
	j.Description = types.StringValue(dto.GetDescription())
	j.Enabled = types.BoolValue(dto.GetEnabled())
	j.Priority = types.Int32Value(int32(dto.GetPriority()))
	j.Config = JobHookConfig{}
	err := j.Config.FromDto(dto.GetConfig())
	if err != nil {
		return err
	}
	return nil
}

func (j *JobHookConfig) ToDto() (*mgmtv1alpha1.JobHookConfig, error) {
	if j == nil {
		return nil, errors.New("JobHookConfig was nil")
	}

	if j.Sql != nil {
		sql, err := j.Sql.ToDto()
		if err != nil {
			return nil, err
		}
		return &mgmtv1alpha1.JobHookConfig{
			Config: sql,
		}, nil
	}

	return nil, errors.New("no valid job hook config was provided")
}
func (j *JobHookConfig) FromDto(dto *mgmtv1alpha1.JobHookConfig) error {
	if dto == nil {
		dto = &mgmtv1alpha1.JobHookConfig{}
	}
	switch cfg := dto.GetConfig().(type) {
	case *mgmtv1alpha1.JobHookConfig_Sql:
		if cfg.Sql == nil {
			return errors.New("config was JobHookConfig_Sql but inner config was nil")
		}
		j.Sql = &JobHookSql{}
		err := j.Sql.FromDto(cfg)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("invalid or not yet added job hook config was specified: %w", errors.ErrUnsupported)
}

func (j *JobHookSql) ToDto() (*mgmtv1alpha1.JobHookConfig_Sql, error) {
	if j == nil {
		return nil, errors.New("JobHookSql was nil")
	}

	timing, err := j.Timing.ToDto()
	if err != nil {
		return nil, err
	}
	return &mgmtv1alpha1.JobHookConfig_Sql{
		Sql: &mgmtv1alpha1.JobHookConfig_JobSqlHook{
			Query:        j.Query.ValueString(),
			ConnectionId: j.ConnectionId.ValueString(),
			Timing:       timing,
		},
	}, nil
}

func (j *JobHookSql) FromDto(dto *mgmtv1alpha1.JobHookConfig_Sql) error {
	if j == nil {
		return errors.New("JobHookSql was nil")
	}
	if dto == nil {
		dto = &mgmtv1alpha1.JobHookConfig_Sql{
			Sql: &mgmtv1alpha1.JobHookConfig_JobSqlHook{},
		}
	}
	if dto.Sql == nil {
		dto.Sql = &mgmtv1alpha1.JobHookConfig_JobSqlHook{}
	}

	j.Query = types.StringValue(dto.Sql.GetQuery())
	j.ConnectionId = types.StringValue(dto.Sql.GetConnectionId())
	j.Timing = JobHookSqlTiming{}
	err := j.Timing.FromDto(dto.Sql.GetTiming())
	if err != nil {
		return err
	}

	return nil
}

func (j *JobHookSqlTiming) ToDto() (*mgmtv1alpha1.JobHookConfig_JobSqlHook_Timing, error) {
	if j == nil {
		return nil, errors.New("JobHookSqlTiming was nil")
	}
	if j.PreSync != nil {
		return &mgmtv1alpha1.JobHookConfig_JobSqlHook_Timing{
			Timing: &mgmtv1alpha1.JobHookConfig_JobSqlHook_Timing_PreSync{
				PreSync: &mgmtv1alpha1.JobHookTimingPreSync{},
			},
		}, nil
	}
	if j.PostSync != nil {
		return &mgmtv1alpha1.JobHookConfig_JobSqlHook_Timing{
			Timing: &mgmtv1alpha1.JobHookConfig_JobSqlHook_Timing_PostSync{
				PostSync: &mgmtv1alpha1.JobHookTimingPostSync{},
			},
		}, nil
	}

	return nil, errors.New("no valid job hook sql timing config was provided")
}

func (j *JobHookSqlTiming) FromDto(dto *mgmtv1alpha1.JobHookConfig_JobSqlHook_Timing) error {
	if j == nil {
		return errors.New("JobHookSqlTiming was nil")
	}
	if dto == nil {
		dto = &mgmtv1alpha1.JobHookConfig_JobSqlHook_Timing{}
	}
	switch dto.GetTiming().(type) {
	case *mgmtv1alpha1.JobHookConfig_JobSqlHook_Timing_PreSync:
		j.PreSync = &JobHookTimingPreSync{}
		return nil
	case *mgmtv1alpha1.JobHookConfig_JobSqlHook_Timing_PostSync:
		j.PostSync = &JobHookTimingPostSync{}
		return nil
	}

	return fmt.Errorf("invalid or not yet added job sql hook timing was specified: %w", errors.ErrUnsupported)
}
