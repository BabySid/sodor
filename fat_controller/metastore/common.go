package metastore

import (
	"errors"
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc/http/codec"
	"github.com/BabySid/proto/sodor"
	"google.golang.org/protobuf/types/known/structpb"
	"time"
)

var (
	ErrNotFound         = errors.New("not found")
	ThomasStatusExpired = "Expired"
	ThomasStatusOK      = "OK"
)

func toJob(in *sodor.Job, out *Job) error {
	out.Name = in.Name

	if in.Id > 0 {
		out.ID = uint(in.Id)
	}

	out.SchedulerMode = in.ScheduleMode.String()
	if in.RoutineSpec != nil {
		temp, err := codec.DefaultProtoMarshal.Marshal(in.RoutineSpec)
		if err != nil {
			return err
		}
		out.RoutineSpec = string(temp)
	}

	out.AlertGroupID = in.AlertGroupId

	return nil
}

func fromJob(in *Job, out *sodor.Job) error {
	out.Id = int32(in.ID)
	out.Name = in.Name
	out.CreateAt = int32(in.CreatedAt.Unix())
	out.UpdateAt = int32(in.UpdatedAt.Unix())

	out.ScheduleMode = sodor.ScheduleMode(sodor.ScheduleMode_value[in.SchedulerMode])

	if out.ScheduleMode == sodor.ScheduleMode_ScheduleMode_Crontab {
		var spec sodor.RoutineSpec
		err := codec.DefaultProtoMarshal.Unmarshal([]byte(in.RoutineSpec), &spec)
		if err != nil {
			return err
		}

		out.RoutineSpec = &spec
	}

	out.AlertGroupId = in.AlertGroupID

	return nil
}

func toTask(in *sodor.Task, jobID int32, out *Task) error {
	if in.Id > 0 {
		out.ID = uint(in.Id)
	}
	out.JobID = jobID
	out.Name = in.Name

	out.Type = in.Type.String()
	out.Content = in.Content
	if in.RunningHosts != nil {
		jsonBytes, err := codec.DefaultProtoMarshal.Marshal(in.RunningHosts)
		if err != nil {
			return err
		}
		out.RunningHosts = string(jsonBytes)
	}

	return nil
}

func fromTask(in *Task, out *sodor.Task) error {
	out.Id = int32(int64(in.ID))
	out.JobId = in.JobID
	out.Name = in.Name
	if in.RunningHosts != "" {
		err := codec.DefaultProtoMarshal.Unmarshal([]byte(in.RunningHosts), &out.RunningHosts)
		if err != nil {
			return err
		}
	}

	out.Type = sodor.TaskType(sodor.TaskType_value[in.Type])
	out.Content = in.Content

	out.CreateAt = int32(in.CreatedAt.Unix())
	out.UpdateAt = int32(in.UpdatedAt.Unix())

	return nil
}

func fromThomas(in *Thomas, out *sodor.ThomasInfo) error {
	out.Id = int32(in.ID)
	out.CreateAt = int32(in.CreatedAt.Unix())
	out.UpdateAt = int32(in.UpdatedAt.Unix())
	out.Name = in.Name
	out.Version = in.Version
	out.Tags = in.Tags
	out.Proto = in.Proto
	out.Host = in.Host
	out.Port = int32(in.Port)
	out.Pid = int32(in.PID)
	out.StartTime = in.StartTime
	out.HeartBeatTime = in.HeartBeatTime
	out.ThomasType = sodor.ThomasType(sodor.ThomasType_value[in.ThomasType])
	out.Status = in.Status
	if out.Status == "" {
		if int64(out.HeartBeatTime) >= time.Now().Unix()-maxThomasLife {
			out.Status = ThomasStatusOK
		} else {
			out.Status = ThomasStatusExpired
		}
	}
	metrics, err := structpb.NewStruct(in.Metrics)
	if err != nil {
		return err
	}
	out.LatestMetrics = metrics
	return nil
}

func fromThomasMetrics(in *ThomasInstance, out *sodor.ThomasMetrics) error {
	out.Id = int32(in.ID)
	out.CreateAt = int32(in.CreatedAt.Unix())
	out.UpdateAt = int32(in.UpdatedAt.Unix())
	metrics, err := structpb.NewStruct(in.Metrics)
	if err != nil {
		return err
	}
	out.Metrics = metrics
	return nil
}

func toThomas(in *sodor.ThomasInfo, out *Thomas) error {
	if in.Id > 0 {
		out.ID = uint(in.Id)
	}

	out.Name = in.Name
	out.Version = in.Version
	out.Tags = in.Tags
	out.Proto = in.Proto
	out.Host = in.Host
	out.Port = int(in.Port)
	out.PID = int(in.Pid)
	out.StartTime = in.StartTime
	out.HeartBeatTime = int32(time.Now().Unix())
	out.ThomasType = in.ThomasType.String()
	out.Status = in.Status
	out.Metrics = in.LatestMetrics.AsMap()

	return nil
}

func toSimpleThomas(in *sodor.ThomasInfo, out *Thomas) error {
	out.Host = in.Host
	out.Port = int(in.Port)
	out.Tags = in.Tags
	out.ThomasType = in.ThomasType.String()

	return nil
}

func toJobInstance(in *sodor.JobInstance, out *JobInstance) error {
	if in.Id > 0 {
		out.ID = uint(in.Id)
	}

	out.JobID = in.JobId
	out.ScheduleTS = in.ScheduleTs
	out.StartTS = in.StartTs
	out.StopTS = in.StopTs
	out.ExitCode = in.ExitCode
	out.ExitMsg = in.ExitMsg

	return nil
}

func toTaskInstance(in *sodor.TaskInstance, out *TaskInstance) error {
	if in.Id > 0 {
		out.ID = uint(in.Id)
	}

	out.JobID = in.JobId
	out.TaskID = in.TaskId
	out.JobInstanceID = in.JobInstanceId
	out.StartTS = in.StartTs
	out.StopTS = in.StopTs
	out.ParsedContent = in.ParsedContent
	out.Host = in.Host
	out.PID = in.Pid
	out.ExitCode = in.ExitCode
	out.ExitMsg = in.ExitMsg
	out.OutputVars = in.OutputVars.AsMap()

	return nil
}

func fromJobInstance(in *JobInstance, out *sodor.JobInstance) error {
	out.Id = int32(in.ID)
	out.CreateAt = int32(in.CreatedAt.Unix())
	out.UpdateAt = int32(in.UpdatedAt.Unix())
	out.JobId = in.JobID
	out.ScheduleTs = in.ScheduleTS
	out.StartTs = in.StartTS
	out.StopTs = in.StopTS
	out.ExitCode = in.ExitCode
	out.ExitMsg = in.ExitMsg

	return nil
}

func fromTaskInstance(in *TaskInstance, out *sodor.TaskInstance) error {
	out.Id = int32(in.ID)
	out.CreateAt = int32(in.CreatedAt.Unix())
	out.UpdateAt = int32(in.UpdatedAt.Unix())
	out.JobId = in.JobID
	out.TaskId = in.TaskID
	out.JobInstanceId = in.JobInstanceID
	out.StartTs = in.StartTS
	out.StopTs = in.StopTS
	out.Host = in.Host
	out.Pid = in.PID
	out.ExitCode = in.ExitCode
	out.ExitMsg = in.ExitMsg
	out.ParsedContent = in.ParsedContent
	outVars, err := structpb.NewStruct(in.OutputVars)
	if err != nil {
		return err
	}
	out.OutputVars = outVars

	return nil
}

func findTaskName(ts []*Task, id uint) string {
	for _, t := range ts {
		if t.ID == id {
			return t.Name
		}
	}
	gobase.AssertHere()
	return ""
}

func findTaskID(ts []*Task, name string) uint {
	for _, t := range ts {
		if t.Name == name {
			return t.ID
		}
	}
	gobase.AssertHere()
	return 0
}

func toAlertPluginInstance(in *sodor.AlertPluginInstance, out *AlertPluginInstance) error {
	if in.Id > 0 {
		out.ID = uint(in.Id)
	}

	out.Name = in.Name
	out.PluginName = in.PluginName
	bs, err := codec.DefaultProtoMarshal.Marshal(in.Dingding)
	if err != nil {
		return err
	}
	out.PluginValue = string(bs)
	return nil
}

func toAlertGroup(in *sodor.AlertGroup, out *AlertGroup) error {
	if in.Id > 0 {
		out.ID = uint(in.Id)
	}

	out.Name = in.Name
	out.PluginInstance = make([]uint, len(in.PluginInstance))
	for i, v := range in.PluginInstance {
		out.PluginInstance[i] = uint(v)
	}

	return nil
}

func fromAlertPluginInstance(in *AlertPluginInstance, out *sodor.AlertPluginInstance) error {
	out.Id = int32(in.ID)
	out.CreateAt = int32(in.CreatedAt.Unix())
	out.UpdateAt = int32(in.UpdatedAt.Unix())
	out.Name = in.Name
	out.PluginName = in.PluginName

	if out.PluginName == sodor.AlertPluginName_APN_DingDing.String() {
		var v sodor.AlertPluginDingDing
		err := codec.DefaultProtoMarshal.Unmarshal([]byte(in.PluginValue), &v)
		if err != nil {
			return err
		}
		out.Dingding = &v
	}

	return nil
}

func fromAlertGroup(in *AlertGroup, out *sodor.AlertGroup) error {
	out.Id = int32(in.ID)
	out.CreateAt = int32(in.CreatedAt.Unix())
	out.UpdateAt = int32(in.UpdatedAt.Unix())
	out.Name = in.Name

	out.PluginInstance = make([]int32, len(in.PluginInstance))
	for i, v := range in.PluginInstance {
		out.PluginInstance[i] = int32(v)
	}

	return nil
}

func fromAlertPluginInstanceHistory(in *AlertPluginInstanceHistory, out *sodor.AlertPluginInstanceHistory) error {
	out.Id = int32(in.ID)
	out.CreateAt = int32(in.CreatedAt.Unix())
	out.UpdateAt = int32(in.UpdatedAt.Unix())
	out.InstanceId = in.InstanceId
	out.GroupId = in.GroupID
	out.AlertMsg = in.AlertMsg
	out.StatusMsg = in.StatusMsg

	return nil
}

func toAlertPluginInstanceHistory(in *sodor.AlertPluginInstanceHistory, out *AlertPluginInstanceHistory) {
	if in.Id > 0 {
		out.ID = uint(in.Id)
	}

	out.InstanceId = in.InstanceId
	out.GroupID = in.GroupId
	out.InstanceId = in.InstanceId
	out.AlertMsg = in.AlertMsg
	out.StatusMsg = in.StatusMsg
}
