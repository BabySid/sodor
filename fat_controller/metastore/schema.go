package metastore

import (
	"github.com/BabySid/proto/sodor"
	"time"
)

// This is the data structure of the storage layer

var (
	totalTables []interface{}
)

func init() {
	totalTables = []interface{}{
		&Job{},
		&Task{},
		&TaskRelation{},
		&AlertGroup{},
		&AlertGroupInstance{},
		&ScheduleState{},
		&JobInstance{},
		&TaskInstance{},
		&Thomas{},
		&ThomasInstance{},
	}
}

type TableModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AlertGroup struct {
	TableModel
	Name string `gorm:"not null;size:64;uniqueIndex:uniq_alert_group"`
	// pluginName => properties
	PluginValues map[sodor.AlertPluginName]map[string]interface{} `gorm:"not null;serializer:json;type:text"`
}

func (t AlertGroup) UpdateFields() []string {
	return []string{
		"Name",
		"PluginValues",
	}
}

type AlertGroupInstance struct {
	TableModel
	InstanceId  int32                  `gorm:"not null"`
	GroupID     int32                  `gorm:"not null"`
	PluginName  string                 `gorm:"not null;size:64"`
	ParamsValue map[string]interface{} `gorm:"not null;serializer:json;type:text"`
}

type Job struct {
	TableModel
	Name          string `gorm:"not null;size:64;unique"`
	SchedulerMode string `gorm:"not null;default:''"`
	RoutineSpec   string `gorm:"not null;default:'';size:128"` // {"ct_spec":"* * *"}
	//AlertRule    string `gorm:"not null;default:'';type:text"` // json
	AlertGroupID int32 `gorm:"not null"`
}

func (t Job) UpdateFields() []string {
	return []string{
		"Name",
		"SchedulerMode",
		"RoutineSpec",
		"AlertGroupID",
	}
}

type Task struct {
	TableModel
	JobID        int32  `gorm:"not null;uniqueIndex:uniq_task"`
	Name         string `gorm:"not null;size:64;uniqueIndex:uniq_task"`
	RunningHosts string `gorm:"not null;default:'';size:256"` // [{"tag":["a","b"]},{"hosts":["1.1.1.1"]}]
	Type         string `gorm:"not null;default:'';size:16"`
	Script       string `gorm:"not null;default:'';type:mediumtext"`
}

func (t Task) UpdateFields() []string {
	return []string{
		"JobID",
		"Name",
		"RunningHosts",
		"Type",
		"Script",
	}
}

type TaskRelation struct {
	TableModel
	JobID      int32 `gorm:"not null"`
	FromTaskID int32 `gorm:"not null"`
	ToTaskID   int32 `gorm:"not null"`
}

// ScheduleState stores jobs with a crontab-scheduler
type ScheduleState struct {
	TableModel
	JobID int32  `gorm:"not null;uniqueIndex:uniq_job"`
	Host  string `gorm:"not null;size:64;uniqueIndex:uniq_job"`
}

type JobInstance struct {
	TableModel
	JobID      int32  `gorm:"not null"`
	ScheduleTS int32  `gorm:"not null:default:0"`
	StartTS    int32  `gorm:"not null;default:0"`
	StopTS     int32  `gorm:"not null;default:0"`
	ExitCode   int32  `gorm:"not null;default:0"`
	ExitMsg    string `gorm:"not null;default:''"`
}

func (t JobInstance) UpdateFields() []string {
	return []string{
		"JobID",
		"ScheduleTS",
		"StartTS",
		"StopTS",
		"ExitCode",
		"ExitMsg",
	}
}

type TaskInstance struct {
	TableModel
	JobID         int32  `gorm:"not null;uniqueIndex:uniq_task"`
	TaskID        int32  `gorm:"not null;uniqueIndex:uniq_task"`
	JobInstanceID int32  `gorm:"not null;uniqueIndex:uniq_task"`
	StartTS       int32  `gorm:"not null;default:0"`
	StopTS        int32  `gorm:"not null;default:0"`
	Host          string `gorm:"not null;default:''"`
	PID           int32  `gorm:"not null;default:0"`
	ExitCode      int32  `gorm:"not null;default:0"`
	ExitMsg       string `gorm:"not null;default:''"`
	//InputVars     string `gorm:"not null;default:'';type:mediumtext"` // json
	OutputVars map[string]interface{} `gorm:"not null;serializer:json;default:'';type:mediumtext"` // json
}

func (t TaskInstance) UpdateFields() []string {
	return []string{
		"JobID",
		"TaskID",
		"JobInstanceID",
		"StartTS",
		"StopTS",
		"Host",
		"PID",
		"ExitCode",
		"ExitMsg",
		"OutputVars",
	}
}

type Thomas struct {
	TableModel
	Name          string                 `gorm:"size:64;not null"`
	Version       string                 `gorm:"size:64;not null"`
	Proto         string                 `gorm:"size:16;not null"`
	Host          string                 `gorm:"size:32;not null;uniqueIndex:uniq_thomas"` // ip. e.g. 1.2.3.4
	Port          int                    `gorm:"not null;uniqueIndex:uniq_thomas"`
	PID           int                    `gorm:"not null;column:pid"`
	StartTime     int32                  `gorm:"not null"`
	HeartBeatTime int32                  `gorm:"not null;column:heart_beat_time"`
	ThomasType    string                 `gorm:"not null;default:'';size:32"`
	Status        string                 `gorm:"not null;default:'';size:256"`
	Metrics       map[string]interface{} `gorm:"not null;serializer:json;default:'';type:mediumtext"` // json
}

type ThomasInstance struct {
	TableModel
	ThomasID int32                  `gorm:"not null"`
	Metrics  map[string]interface{} `gorm:"not null;serializer:json;default:'';type:mediumtext"` // json
}

func (t Thomas) TableName() string {
	return "thomas"
}

func (t Thomas) UpdateFields() []string {
	return []string{
		"Name",
		"Version",
		"Proto",
		"PID",
		"StartTime",
		"HeartBeatTime",
		"Status",
		"Metrics",
	}
}
