package api

import (
	"encoding/json"
	"uspacy-go-sdk/task"
)

// CreateTransferTask creates a new transfer task
func (us *Uspacy) CreateTransferTask() (task.TransferTaskOutput, error) {
	var tasks task.TransferTaskOutput

	body, err := us.doGetEmptyHeaders(buildURL(mainHost, task.VersionTaskUrl, task.TasksUrl))
	if err != nil {
		return tasks, err
	}
	return tasks, json.Unmarshal(body, &tasks)
}
