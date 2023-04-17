package api

import (
	"encoding/json"

	"github.com/Uspacy/uspacy-go-sdk/task"
)

// CreateTask creates a new task
func (us *Uspacy) CreateTask(taskData task.Task) (_task task.Task, err error) {
	resp, err := us.doPostEmptyHeaders(us.buildURL(task.VersionUrl, task.TaskUrl), taskData)
	if err != nil {
		return _task, err
	}
	return _task, json.Unmarshal(resp, &_task)
}

// CreateTransferTask creates a new transfer task
func (us *Uspacy) CreateTransferTask(body interface{}) (tasks task.TransferTaskOutput, err error) {
	resp, err := us.doPostEmptyHeaders(us.buildURL(task.VersionUrl, task.TransferUrl), body)
	if err != nil {
		return tasks, err
	}
	return tasks, json.Unmarshal(resp, &tasks)
}
