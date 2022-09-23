package metastore

import (
	"errors"
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc/http/codec"
	"github.com/BabySid/proto/sodor"
)

var (
	ErrNotFound = errors.New("not Found")
)

func toJob(in *sodor.Job, out *Job) error {
	out.Name = in.Name

	if in.Id > 0 {
		out.ID = uint(in.Id)
	}

	return nil
}

func fromJob(in *Job, out *sodor.Job) error {
	out.Id = int64(in.ID)
	out.Name = in.Name
	out.CreateAt = in.CreatedAt.Unix()
	out.UpdateAt = in.UpdatedAt.Unix()

	return nil
}

func toTask(in *sodor.Task, jobID int64, out *Task) error {
	if in.Id > 0 {
		out.ID = uint(in.Id)
	}
	out.JobID = jobID
	out.Name = in.Name

	out.SchedulerMode = in.SchedulerMode.String()
	if in.RoutineSpec != nil {
		temp, err := codec.DefaultProtoMarshal.Marshal(in.RoutineSpec)
		if err != nil {
			return err
		}
		out.RoutineSpec = string(temp)
	}

	out.Script = in.Script
	if in.RunningHosts != nil {
		jsonBytes, err := codec.DefaultProtoMarshal.Marshal(in.RunningHosts)
		if err != nil {
			return err
		}
		out.RunningHosts = string(jsonBytes)
	}

	out.RunTimeout = int(in.RunningTimeout)

	return nil
}

func fromTask(in *Task, out *sodor.Task) error {
	out.Id = int64(in.ID)
	out.JobId = in.JobID
	out.Name = in.Name
	if in.RunningHosts != "" {
		err := codec.DefaultProtoMarshal.Unmarshal([]byte(in.RunningHosts), &out.RunningHosts)
		if err != nil {
			return err
		}
	}

	out.Script = in.Script
	out.SchedulerMode = sodor.SchedulerMode(sodor.SchedulerMode_value[in.SchedulerMode])

	if out.SchedulerMode == sodor.SchedulerMode_SM_Crontab {
		var spec sodor.RoutineSpec
		err := codec.DefaultProtoMarshal.Unmarshal([]byte(in.RoutineSpec), &spec)
		if err != nil {
			return err
		}

		out.RoutineSpec = &spec
	}

	out.RunningTimeout = int32(in.RunTimeout)

	out.CreateAt = in.CreatedAt.Unix()
	out.UpdateAt = in.UpdatedAt.Unix()

	return nil
}

func findTaskName(ts []Task, id uint) string {
	for _, t := range ts {
		if t.ID == id {
			return t.Name
		}
	}
	gobase.AssertHere()
	return ""
}

func findTaskID(ts []Task, name string) uint {
	for _, t := range ts {
		if t.Name == name {
			return t.ID
		}
	}
	gobase.AssertHere()
	return 0
}
