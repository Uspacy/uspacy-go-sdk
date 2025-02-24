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
func (us *Uspacy) CreateTaskThroughMap(taskData map[string]interface{}, headers ...map[string]string) (_task task.Task, err error) {
	resp, _, err := us.doPost(us.buildURL(task.VersionUrl, task.TaskUrl), taskData, headers...)
	if err != nil {
		return _task, err
	}
	return _task, json.Unmarshal(resp, &_task)
}

// CreateTransferTask creates a new transfer task
func (us *Uspacy) CreateTransferTask(body interface{}, headers ...map[string]string) (tasks task.TransferTaskOutput, err error) {
	resp, _, err := us.doPost(us.buildURL(task.VersionUrl, task.TransferUrl), body, headers...)
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

// GetFields returns Fields struct
func (us *Uspacy) GetTaskFields() (fields []task.Field, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(task.VersionUrl, "custom_fields", task.TaskUrl, "fields"))
	if err != nil {
		return fields, err
	}
	var resp task.TaskFields
	return resp.Fields, json.Unmarshal(body, &resp)
}

// GetTasksList returns TasksList struct
func (us *Uspacy) GetTasksList(params url.Values) (tasks task.TasksList, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(task.VersionUrl, task.TaskUrl) + "?" + params.Encode())
	if err != nil {
		return tasks, err
	}
	var resp task.TasksList
	return resp, json.Unmarshal(body, &resp)
}

// GetTaskStagesByGroupId
func (us *Uspacy) GetTaskStagesByGroupId(groupId int) (kanbanStages []task.TaskGroupStage, err error) {
	params := url.Values{}
	params.Set("groupId", fmt.Sprintf("%d", groupId))
	body, err := us.doGetEmptyHeaders(us.buildURL(task.VersionUrl, task.KanbanStages) + "?" + params.Encode())
	if err != nil {
		return kanbanStages, err
	}
	var resp task.TaskGroupStages
	return resp.Data, json.Unmarshal(body, &resp)
}

// GetTempleateById
func (us *Uspacy) GetTemplateById(templateId int) (template task.Template, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(task.VersionUrl, task.TemplateUrl, fmt.Sprintf("%d", templateId)))
	if err != nil {
		return template, err
	}
	var resp task.Template
	return resp, json.Unmarshal(body, &resp)
}
