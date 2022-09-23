package metastore

import "gorm.io/gorm"

// This is the data structure of the storage layer

var (
	totalTable []interface{}
)

func init() {
	totalTable = []interface{}{
		&Job{},
		&Task{},
		&TaskRelation{},
		&AlertGroup{},
		&AlertPlugin{},
		&AlertHistory{},
		&JobInstance{},
		&TaskInstance{},
		&Thomas{},
	}
}

type Job struct {
	gorm.Model
	Name string `gorm:"not null;size:64;unique"`
	//AlertRule    string `gorm:"not null;default:'';type:text"` // json
	//AlertGroupID int32  `gorm:"not null"`
}

type Task struct {
	gorm.Model
	JobID         int32  `gorm:"not null;uniqueIndex:uniq_task"`
	Name          string `gorm:"not null;size:64;uniqueIndex:uniq_task"`
	RunningHosts  string `gorm:"not null;default:'';size:256"` // [{"tag":["a","b"]},{"hosts":["1.1.1.1"]}]
	SchedulerMode string `gorm:"not null;default:''"`
	RoutineSpec   string `gorm:"not null;default:'';size:128"` // {"ct_spec":"* * *"}
	Script        string `gorm:"not null;default:'';type:mediumtext"`
	RunTimeout    int    `gorm:"not null;default:0"` // seconds
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

type JobInstance struct {
	gorm.Model
	JobID    int32  `gorm:"not null"`
	StartTS  int32  `gorm:"not null;default:0"`
	StopTS   int32  `gorm:"not null;default:0"`
	Progress string `gorm:"not null;default:''"` // json
	ExitCode int32  `gorm:"not null;default:0"`
	ExitMsg  string `gorm:"not null;default:''"`
}

type TaskInstance struct {
	gorm.Model
	JobID         int32  `gorm:"not null;uniqueIndex:uniq_task"`
	TaskID        int32  `gorm:"not null;uniqueIndex:uniq_task"`
	JobInstanceID int32  `gorm:"not null;uniqueIndex:uniq_task"`
	StartTS       int32  `gorm:"not null;default:0"`
	StopTS        int32  `gorm:"not null;default:0"`
	PID           int32  `gorm:"not null;default:0"`
	ExitCode      int32  `gorm:"not null;default:0"`
	ExitMsg       string `gorm:"not null;default:''"`
	InputVars     string `gorm:"not null;default:'';type:mediumtext"` // json
	OutputVars    string `gorm:"not null;default:'';type:mediumtext"` // json
}

type Thomas struct {
	gorm.Model
	Version           string `gorm:"size:64;not null"`
	Proto             string `gorm:"size:16;not null"`
	Host              string `gorm:"size:32;not null"` // ip. e.g. 1.2.3.4
	Port              int    `gorm:"not null"`
	PID               int    `gorm:"not null;column:pid"`
	Tags              string `gorm:"not null;default:'';size:64"`
	LastStartTime     int32  `gorm:"not null"`
	LastHeartbeatTime int32  `gorm:"not null"`
}

func (t Thomas) TableName() string {
	return "thomas"
}
