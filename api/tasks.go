package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/task"
)

// CreateTask creates a new task
func (us *Uspacy) CreateTask(taskData url.Values) (_task task.Task, err error) {
	resp, err := us.doPostEncodedForm(us.buildURL(task.VersionUrl, task.TaskUrl), taskData)
	if err != nil {
		return _task, err
	}
	return _task, json.Unmarshal(resp, &_task)
}

// CreateTaskThroughMap creates a new task through a map
func (us *Uspacy) CreateTaskThroughMap(taskData map[string]interface{}) (_task task.Task, err error) {
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

// PatchTask patch task by Id
func (us *Uspacy) PatchTask(taskId int, taskData map[string]interface{}) (_task task.Task, err error) {
	resp, err := us.doPatchEmptyHeaders(us.buildURL(task.VersionUrl, fmt.Sprintf(task.TaskIdUrl, taskId)), taskData)
	if err != nil {
		return _task, err
	}
	return _task, json.Unmarshal(resp, &_task)
}
