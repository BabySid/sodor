package task_runner

import (
	"encoding/json"
	"github.com/BabySid/gobase"
	"io/ioutil"
	"path/filepath"
	"sodor/thomas/config"
	"sync"
)

type metaDB struct {
	mutex sync.Mutex
	// taskWorkDir->bool
	metaOfTasks map[string]interface{}
}

func newMetaDB() *metaDB {
	return &metaDB{
		metaOfTasks: make(map[string]interface{}),
	}
}

const (
	taskMetaDB = "task_meta.db"
)

func (m *metaDB) inertTaskMeta(key string, value interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.metaOfTasks[key] = value

	bs, err := json.Marshal(m.metaOfTasks)
	if err != nil {
		return err
	}

	err = gobase.WriteFile(filepath.Join(config.GetInstance().DataPath, taskMetaDB), bs, defaultPerm)
	if err != nil {
		return err
	}

	return nil
}

func (m *metaDB) remove(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.metaOfTasks, key)

	bs, err := json.Marshal(m.metaOfTasks)
	if err != nil {
		return err
	}

	err = gobase.WriteFile(filepath.Join(config.GetInstance().DataPath, taskMetaDB), bs, defaultPerm)
	if err != nil {
		return err
	}

	return nil
}

func (m *metaDB) load() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ok, err := gobase.PathExists(filepath.Join(config.GetInstance().DataPath, taskMetaDB))
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	bs, err := ioutil.ReadFile(filepath.Join(config.GetInstance().DataPath, taskMetaDB))
	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, &m.metaOfTasks)
	return err
}

func (m *metaDB) Traversal(call func(k string, v interface{}) error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for k, v := range m.metaOfTasks {
		if err := call(k, v); err != nil {
			break
		}
	}
}
