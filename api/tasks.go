package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/task"
)

// CreateTask creates a new task
func (us *Uspacy) CreateTask(taskData url.Values) (newTask task.Task, err error) {
	resp, err := us.doPostEncodedForm(us.buildURL(task.VersionUrl, task.TaskUrl), taskData)
	if err != nil {
		return newTask, err
	}
	return newTask, json.Unmarshal(resp, &newTask)
}

// CreateTaskThroughMap creates a new task through a map
func (us *Uspacy) CreateTaskThroughMap(taskData map[string]any, headers ...map[string]string) (newTask task.Task, statusCode int, err error) {
	resp, code, err := us.doPost(us.buildURL(task.VersionUrl, task.TaskUrl), taskData, headers...)
	if err != nil {
		return newTask, code, err
	}
	return newTask, code, json.Unmarshal(resp, &newTask)
}

// CreateTransferTask creates a new transfer task
func (us *Uspacy) CreateTransferTask(body any, headers ...map[string]string) (tasks task.TransferTaskOutput, statusCode int, err error) {
	resp, code, err := us.doPost(us.buildURL(task.VersionUrl, task.TransferUrl), body, headers...)
	if err != nil {
		return tasks, code, err
	}
	return tasks, code, json.Unmarshal(resp, &tasks)
}

// PatchTask patch task by Id
func (us *Uspacy) PatchTask(taskId int, taskData map[string]any) (updatedTask task.Task, err error) {
	resp, err := us.doPatchEmptyHeaders(us.buildURL(task.VersionUrl, fmt.Sprintf(task.TaskIdUrl, taskId)), taskData)
	if err != nil {
		return updatedTask, err
	}
	return updatedTask, json.Unmarshal(resp, &updatedTask)
}

// GetFields returns Fields struct
func (us *Uspacy) GetTaskFields() (fields []task.Field, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(task.VersionUrl, task.TaskUrl, task.FieldUrl))
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

// CreateTaskStage creates a new task stage
func (us *Uspacy) CreateTaskStage(stageData task.TaskGroupStage) (kanbanStage task.TaskGroupStage, statusCode int, err error) {
	body, code, err := us.doPost(us.buildURL(task.VersionUrl, task.KanbanStages), stageData)
	if err != nil {
		return kanbanStage, code, err
	}
	var resp task.TaskGroupStage
	return resp, code, json.Unmarshal(body, &resp)
}

// DeleteTaskStage deletes a task stage
func (us *Uspacy) DeleteTaskStage(stageId int) (err error) {
	_, err = us.doDeleteEmptyHeaders(us.buildURL(task.VersionUrl, task.KanbanStages, fmt.Sprintf("%d", stageId)), nil)
	return err
}

// TaskStatusReady marks task as ready
func (us *Uspacy) TaskStatusReady(taskId int) (err error) {
	_, err = us.doPatchEmptyHeaders(us.buildURL(task.VersionUrl, fmt.Sprintf(task.TaskIdUrl, taskId), task.TaskStatusReady), nil)
	return err
}

// CreateTaskField creates a new task field
func (us *Uspacy) CreateTaskField(fieldData task.Field) (field task.Field, statusCode int, err error) {
	resp, code, err := us.doPost(us.buildURL(task.VersionUrl, task.TaskUrl, task.FieldUrl), fieldData)
	if err != nil {
		return field, code, err
	}
	var respField task.Field
	err = json.Unmarshal(resp, &respField)
	return respField, code, err
}
