package variables

import "sodor/fat_controller/metastore"

type VariableHandle interface {
	Handle() (interface{}, error)
}

type LastSuccessVars struct {
	taskId int32
}

func (v *LastSuccessVars) SetTaskID(id int32) {
	v.taskId = id
}

func (v *LastSuccessVars) Handle() (interface{}, error) {
	ins, err := metastore.GetInstance().SelectLastTaskInstance(v.taskId)
	if err != nil {
		return nil, err
	}

	if ins == nil {
		return "", nil
	}

	return ins.OutputVars.AsMap(), nil
}
