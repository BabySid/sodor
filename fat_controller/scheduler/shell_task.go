package scheduler

import "github.com/BabySid/proto/sodor"

func parseTaskContent(task *sodor.Task, ins *sodor.TaskInstance) {
	ins.ParsedContent = task.Content
}
