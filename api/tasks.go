package api

import (
	"encoding/json"

	"github.com/Uspacy/uspacy-go-sdk/task"
)

// CreateTransferTask creates a new transfer task
func (us *Uspacy) CreateTransferTask(body interface{}) (tasks task.TransferTaskOutput, err error) {
	resp, err := us.doPostEmptyHeaders(buildURL(mainHost, task.VersionUrl, task.TransferUrl), body)
	if err != nil {
		return tasks, err
	}
	return tasks, json.Unmarshal(resp, &tasks)
}
