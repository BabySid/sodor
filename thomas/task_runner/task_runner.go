package task_runner

import (
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc/http/codec"
	"github.com/BabySid/proto/sodor"
	"io/ioutil"
	"os"
	"sodor/base"
	"time"
)

type Task interface {
	Run()
}

func GetRunner(typ sodor.TaskType) Task {
	if typ == sodor.TaskType_TaskType_Shell {
		return &ShellRunner{}
	}

	gobase.AssertHere()
	return nil
}

type TaskRunner struct {
	request  *sodor.RunTaskRequest
	response *sodor.TaskInstance
}

const (
	requestFile  = "./task_request.json"
	responseFile = "./task_response.json"
	defaultPerm  = 666
)

func NewTaskRunner() *TaskRunner {
	return &TaskRunner{}
}

func (r *TaskRunner) SetUp() error {
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
	resp.Id = req.TaskInstanceId
	resp.JobId = req.JobId
	resp.TaskId = req.TaskId
	resp.JobInstanceId = req.JobInstanceId
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

func (r *TaskRunner) run() {

}