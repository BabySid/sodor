package metastore

import (
	"gorm.io/gorm"
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
		&AlertPlugin{},
		&AlertHistory{},
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

type Job struct {
	gorm.Model
	Name          string `gorm:"not null;size:64;unique"`
	SchedulerMode string `gorm:"not null;default:''"`
	RoutineSpec   string `gorm:"not null;default:'';size:128"` // {"ct_spec":"* * *"}
	//AlertRule    string `gorm:"not null;default:'';type:text"` // json
	//AlertGroupID int32  `gorm:"not null"`
}

type Task struct {
	gorm.Model
	JobID        int32  `gorm:"not null;uniqueIndex:uniq_task"`
	Name         string `gorm:"not null;size:64;uniqueIndex:uniq_task"`
	RunningHosts string `gorm:"not null;default:'';size:256"` // [{"tag":["a","b"]},{"hosts":["1.1.1.1"]}]
	Type         string `gorm:"not null;default:'';size:16"`
	Script       string `gorm:"not null;default:'';type:mediumtext"`
}

type TaskRelation struct {
	gorm.Model
	JobID      int32 `gorm:"not null"`
	FromTaskID int32 `gorm:"not null"`
	ToTaskID   int32 `gorm:"not null"`
}

type AlertGroup struct {
	gorm.Model
	Name      string `gorm:"not null;size:64;uniqueIndex:uniq_alter_group"`
	PluginIDs string `gorm:"not null;size:64;column:plugin_ids"` // json [1,2]
}

type AlertPlugin struct {
	gorm.Model
	Catalog string `gorm:"not null;size:32;uniqueIndex:uniq_plugin"`
	Params  string `gorm:"not null;type:text"`
}

type AlertHistory struct {
	gorm.Model
	GroupID     int32  `gorm:"not null"`
	PluginID    int32  `gorm:"not null"`
	ParamsValue string `gorm:"not null;type:text"`
}

// ScheduleState stores jobs with a crontab-scheduler
type ScheduleState struct {
	TableModel
	JobID int32  `gorm:"not null;uniqueIndex:uniq_job"`
	Host  string `gorm:"not null;size:64;uniqueIndex:uniq_job"`
}

type JobInstance struct {
	gorm.Model
	JobID      int32  `gorm:"not null"`
	ScheduleTS int32  `gorm:"not null:default:0"`
	StartTS    int32  `gorm:"not null;default:0"`
	StopTS     int32  `gorm:"not null;default:0"`
	ExitCode   int32  `gorm:"not null;default:0"`
	ExitMsg    string `gorm:"not null;default:''"`
}

type TaskInstance struct {
	gorm.Model
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



