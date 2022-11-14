package task_runner

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc/http/codec"
	"github.com/BabySid/proto/sodor"
	"github.com/go-cmd/cmd"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"sodor/thomas/config"
	"sodor/thomas/fat_ctrl"
	"strconv"
	"sync"
	"time"
)

type taskEnv struct {
	db *metaDB
}

var (
	once      sync.Once
	singleton *taskEnv
)

func GetTaskEnv() *taskEnv {
	once.Do(func() {
		singleton = &taskEnv{
			db: newMetaDB(),
		}
	})
	return singleton
}

type CmdContext struct {
	*cmd.Cmd
	ID string
}

func (e *taskEnv) SetUp(req *sodor.RunTaskRequest) (*CmdContext, error) {
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

	err = e.db.inertTaskMeta(path, true)
	if err != nil {
		return nil, err
	}
	c := createCmd(path)
	return &CmdContext{Cmd: c, ID: path}, nil
}

func (e *taskEnv) LoadTasksStatus() error {
	err := e.db.load()
	if err != nil {
		return err
	}

	taskInstances := make(map[string]*sodor.TaskInstance)
	e.db.Traversal(func(k string, v interface{}) error {
		resp, _ := e.GetTaskResponse(k)
		if resp == nil {
			log.Warnf("GetTaskResponse(%s) return nil resp", k)
			return nil
		}

		taskInstances[k] = resp
		return nil
	})

	go e.updateTaskInstances(taskInstances)

	return nil
}

func (e *taskEnv) updateTaskInstances(taskInstances map[string]*sodor.TaskInstance) {
	for path, ins := range taskInstances {
		if !gobase.IsProcessAlive(int(ins.Pid)) {
			if err := fat_ctrl.GetInstance().UpdateTaskInstance(ins); err != nil {
				log.Warnf("UpdateTaskInstance(%s) failed. retry after %v", path, config.GetInstance().RetryInterval)
				time.Sleep(config.GetInstance().RetryInterval)
				continue
			}
			delete(taskInstances, path)
			e.Remove(path)
		}
	}
}

func (e *taskEnv) Remove(task string) {
	_ = e.db.remove(task)
}

func (e *taskEnv) GetTaskResponse(taskPath string) (*sodor.TaskInstance, error) {
	fName := filepath.Join(taskPath, responseFile)
	ok, err := gobase.PathExists(fName)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	bs, err := ioutil.ReadFile(fName)
	if err != nil {
		return nil, err
	}

	var ins sodor.TaskInstance
	err = codec.DefaultProtoMarshal.Unmarshal(bs, &ins)
	if err != nil {
		return nil, err
	}

	return &ins, nil
}