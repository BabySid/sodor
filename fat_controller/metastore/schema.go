package metastore

import "gorm.io/gorm"

var (
	totalTable []interface{}
)

func init() {
	totalTable = []interface{}{
		&Job{},
		&Task{},
		&TaskRelation{},
		&Alert{},
		&TaskInstance{},
		&Thomas{},
	}
}

type Job struct {
	gorm.Model
	Name string `gorm:"not null;size:64;unique"`
}

type Task struct {
	gorm.Model
	JobID           int64  `gorm:"not null;uniqueIndex:uniq_task"`
	Name            string `gorm:"not null;size:64;uniqueIndex:uniq_task"`
	Host            string `gorm:"not null;default:'';size:64"`
	SchedulerMode   int    `gorm:"not null;default:0"`
	SchedulerCTSpec string `gorm:"not null;default:'';size:32"`
	Script          string `grom:"not null;default:'';type:mediumtext"`
}

type TaskRelation struct {
	gorm.Model
}

type Alert struct {
	gorm.Model
}

type TaskInstance struct {
	gorm.Model
	JobID      int64  `gorm:"not null"`
	InstanceID int64  `gorm:"not null"`
	StartTS    int64  `gorm:"not null;default:0"`
	StopTS     int64  `gorm:"not null;default:0"`
	PID        int64  `gorm:"not null;default:0"`
	ExitCode   int64  `gorm:"not null;default:0"`
	ExitMsg    string `gorm:"not null;default:''"`
	Vars       string `gorm:"not null;default:'';type:mediumtext"`
}

type Thomas struct {
	gorm.Model
	Version           string `gorm:"size:64;not null"`
	Proto             string `gorm:"size:16;not null"`
	Host              string `gorm:"size:32;not null"`
	Port              int    `gorm:"not null"`
	PID               int    `gorm:"not null"`
	Tags              string `gorm:"not null;default:'';size:64"`
	LastStartTime     int64  `gorm:"not null"`
	LastHeartbeatTime int64  `gorm:"not null"`
}

func (t Thomas) TableName() string {
	return "thomas"
}
