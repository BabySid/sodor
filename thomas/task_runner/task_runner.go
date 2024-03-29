package task_runner

import (
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc/codec"
	"github.com/BabySid/proto/sodor"
	"io/ioutil"
	"os"
	"sodor/base"
	"sodor/thomas/config"
	"time"
)

type Task interface {
	Run() error
}

func GetRunner() Task {
	return &ShellRunner{}
}

type TaskRunner struct {
	request  *sodor.RunTaskRequest
	response *sodor.TaskInstance
}

const (
	requestFile  = "./task_request.json"
	responseFile = "./task_response.json"
	defaultPerm  = os.ModePerm

	OKMsg       = "OK"
	SystemError = 999
)

func NewTaskRunner() *TaskRunner {
	return &TaskRunner{}
}

func (r *TaskRunner) SetUp() error {
	prefix := "[" + config.GetInstance().TaskIdentity + "] "
	Info.SetPrefix(prefix)
	Warn.SetPrefix(prefix)

	reqBytes, err := ioutil.ReadFile(requestFile)
	if err != nil {
		return err
	}

	var req sodor.RunTaskRequest
	err = codec.DefaultProtoMarshal.Unmarshal(reqBytes, &req)
	if err != nil {
		return err
	}

	var resp sodor.TaskInstance
	resp.Id = req.TaskInstance.Id
	resp.JobId = req.Task.JobId
	resp.TaskId = req.Task.Id
	resp.JobInstanceId = req.TaskInstance.JobInstanceId
	resp.ParsedContent = req.TaskInstance.ParsedContent
	resp.StartTs = int32(time.Now().Unix())
	resp.Host = base.LocalHost
	resp.Pid = int32(os.Getpid())

	respByte, err := codec.DefaultProtoMarshal.Marshal(&resp)
	if err != nil {
		return err
	}

	err = gobase.WriteFile(responseFile, respByte, defaultPerm)
	if err != nil {
		return err
	}

	r.request = &req
	r.response = &resp
	return nil
}

func (r *TaskRunner) TearDown() error {
	r.response.StopTs = int32(time.Now().Unix())

	respByte, err := codec.DefaultProtoMarshal.Marshal(r.response)
	if err != nil {
		return err
	}

	err = gobase.WriteFile(responseFile, respByte, defaultPerm)
	if err != nil {
		return err
	}
	return nil
}
