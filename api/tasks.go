package api

import (
	"encoding/json"
	"uspacy-go-sdk/task"
)

// CreateTransferTask creates a new transfer task
func (us *Uspacy) CreateTransferTask(body interface{}) (task.TransferTaskOutput, error) {
	var tasks task.TransferTaskOutput

	resp, err := us.doPostEmptyHeaders(buildURL(mainHost, task.VersionTaskUrl, task.TransferTasksUrl), body)
	if err != nil {
		return tasks, err
	}
	return tasks, json.Unmarshal(resp, &tasks)
}
