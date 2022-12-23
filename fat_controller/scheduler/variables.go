package scheduler

import (
	"encoding/json"
	"github.com/BabySid/proto/sodor"
	"sodor/fat_controller/variables"
	"strings"
)

func parseTaskContent(task *sodor.Task, ins *sodor.TaskInstance) error {
	// todo switch task.Type
	content := task.Content
	err := variables.ReplaceVariables(task.Content, func(term string, handle variables.VariableHandle) error {
		if handle == nil {
			return nil
		}

		if term == variables.SystemLastSuccessVars {
			h := handle.(*variables.LastSuccessVars)
			h.SetTaskID(task.Id)
			data, err := h.Handle()
			if err != nil {
				return err
			}

			bs, err := json.Marshal(data)
			if err != nil {
				return err
			}

			content = strings.ReplaceAll(content, variables.SystemLastSuccessVars, string(bs))
		}

		return nil
	})

	ins.ParsedContent = content
	return err
}
