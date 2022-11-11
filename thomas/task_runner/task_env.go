package task_runner

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc/http/codec"
	"github.com/BabySid/proto/sodor"
	"github.com/go-cmd/cmd"
	"os"
	"path/filepath"
	"sodor/thomas/config"
	"strconv"
)

type TaskEnv struct {
}

func (e *TaskEnv) SetUp(req *sodor.RunTaskRequest) (*cmd.Cmd, error) {
	path := filepath.Join(
		config.GetInstance().DataPath,
		strconv.Itoa(int(req.JobId)),
		strconv.Itoa(int(req.TaskId)),
		fmt.Sprintf("%d_%d", req.JobInstanceId, req.TaskInstanceId))
	err := os.MkdirAll(path, defaultPerm)
	if err != nil {
		return nil, err
	}

	respByte, err := codec.DefaultProtoMarshal.Marshal(req)
	if err != nil {
		return nil, err
	}

	err = gobase.WriteFile(filepath.Join(path, requestFile), respByte, defaultPerm)
	if err != nil {
		return nil, err
	}

	c := createCmd(path)
	return c, nil
}
