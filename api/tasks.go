package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/v2/task"
)

// CreateTask creates a new task
func (us *Uspacy) CreateTask(ctx context.Context, taskData url.Values) (newTask task.Task, err error) {
	resp, err := us.doPostEncodedForm(ctx, us.buildURL(task.VersionUrl, task.TaskUrl), taskData)
	if err != nil {
		return newTask, err
	}
	return newTask, json.Unmarshal(resp, &newTask)
}

// CreateTaskThroughMap creates a new task through a map
func (us *Uspacy) CreateTaskThroughMap(ctx context.Context, taskData map[string]any, opts ...RequestOption) (newTask task.Task, statusCode int, err error) {
	resp, code, err := us.doPost(ctx, us.buildURL(task.VersionUrl, task.TaskUrl), taskData, opts...)
	if err != nil {
		return newTask, code, err
	}
	return newTask, code, json.Unmarshal(resp, &newTask)
}

// CreateTransferTask creates a new transfer task
func (us *Uspacy) CreateTransferTask(ctx context.Context, body any, opts ...RequestOption) (tasks task.TransferTaskOutput, statusCode int, err error) {
	resp, code, err := us.doPost(ctx, us.buildURL(task.VersionUrl, task.TransferUrl), body, opts...)
	if err != nil {
		return tasks, code, err
	}
	return tasks, code, json.Unmarshal(resp, &tasks)
}

// PatchTask patch task by Id
func (us *Uspacy) PatchTask(ctx context.Context, taskId int, taskData map[string]any) (updatedTask task.Task, err error) {
	resp, err := us.doPatchEmptyHeaders(ctx, us.buildURL(task.VersionUrl, fmt.Sprintf(task.TaskIdUrl, taskId)), taskData)
	if err != nil {
		return updatedTask, err
	}
	return updatedTask, json.Unmarshal(resp, &updatedTask)
}

// GetTaskFields returns Fields struct
func (us *Uspacy) GetTaskFields(ctx context.Context) (fields []task.Field, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(task.VersionUrl, task.TaskUrl, task.FieldUrl))
	if err != nil {
		return fields, err
	}
	var resp task.TaskFields
	if err := json.Unmarshal(body, &resp); err != nil {
		return fields, err
	}
	return resp.Fields, nil
}

// GetTasksList returns TasksList struct
func (us *Uspacy) GetTasksList(ctx context.Context, params url.Values) (tasks task.TasksList, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(task.VersionUrl, task.TaskUrl)+"?"+params.Encode())
	if err != nil {
		return tasks, err
	}
	var resp task.TasksList
	return resp, json.Unmarshal(body, &resp)
}

// GetTasksWithFilters returns tasks with filters as a map
func (us *Uspacy) GetTasksWithFilters(ctx context.Context, params url.Values) (tasks []map[string]any, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(task.VersionUrl, task.TaskUrl)+"?"+params.Encode())
	if err != nil {
		return tasks, err
	}
	var result map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return tasks, err
	}
	// Extract data array from response
	if data, ok := result["data"].([]any); ok {
		tasks = make([]map[string]any, len(data))
		for i, item := range data {
			if taskMap, ok := item.(map[string]any); ok {
				tasks[i] = taskMap
			}
		}
	}
	return tasks, nil
}

// GetTaskById returns task by ID as a map
func (us *Uspacy) GetTaskById(ctx context.Context, taskId int, params ...url.Values) (taskData map[string]any, err error) {
	var urlStr string
	if len(params) > 0 && params[0] != nil {
		urlStr = us.buildURL(task.VersionUrl, fmt.Sprintf(task.TaskIdUrl, taskId)) + "?" + params[0].Encode()
	} else {
		urlStr = us.buildURL(task.VersionUrl, fmt.Sprintf(task.TaskIdUrl, taskId))
	}
	body, err := us.doGetEmptyHeaders(ctx, urlStr)
	if err != nil {
		return taskData, err
	}
	var result map[string]any
	return result, json.Unmarshal(body, &result)
}

// GetTaskStagesByGroupId returns task stages by group id
func (us *Uspacy) GetTaskStagesByGroupId(ctx context.Context, groupId int) (kanbanStages []task.TaskGroupStage, err error) {
	params := url.Values{}
	params.Set("groupId", fmt.Sprintf("%d", groupId))
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(task.VersionUrl, task.KanbanStages)+"?"+params.Encode())
	if err != nil {
		return kanbanStages, err
	}
	var resp task.TaskGroupStages
	if err := json.Unmarshal(body, &resp); err != nil {
		return kanbanStages, err
	}
	return resp.Data, nil
}

// GetTemplateById returns template by id
func (us *Uspacy) GetTemplateById(ctx context.Context, templateId int) (template task.Template, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(task.VersionUrl, task.TemplateUrl, fmt.Sprintf("%d", templateId)))
	if err != nil {
		return template, err
	}
	var resp task.Template
	return resp, json.Unmarshal(body, &resp)
}

// CreateTaskStage creates a new task stage
func (us *Uspacy) CreateTaskStage(ctx context.Context, stageData task.TaskGroupStage) (kanbanStage task.TaskGroupStage, statusCode int, err error) {
	body, code, err := us.doPost(ctx, us.buildURL(task.VersionUrl, task.KanbanStages), stageData)
	if err != nil {
		return kanbanStage, code, err
	}
	var resp task.TaskGroupStage
	return resp, code, json.Unmarshal(body, &resp)
}

// DeleteTaskStage deletes a task stage
func (us *Uspacy) DeleteTaskStage(ctx context.Context, stageId int) (err error) {
	_, err = us.doDeleteEmptyHeaders(ctx, us.buildURL(task.VersionUrl, task.KanbanStages, fmt.Sprintf("%d", stageId)), nil)
	return err
}

// TaskStatusReady marks task as ready
func (us *Uspacy) TaskStatusReady(ctx context.Context, taskId int) (err error) {
	_, err = us.doPatchEmptyHeaders(ctx, us.buildURL(task.VersionUrl, fmt.Sprintf(task.TaskIdUrl, taskId), task.TaskStatusReady), nil)
	return err
}

// CreateTaskField creates a new task field
func (us *Uspacy) CreateTaskField(ctx context.Context, fieldData task.Field) (field task.Field, statusCode int, err error) {
	resp, code, err := us.doPost(ctx, us.buildURL(task.VersionUrl, task.TaskUrl, task.FieldUrl), fieldData)
	if err != nil {
		return field, code, err
	}
	var respField task.Field
	err = json.Unmarshal(resp, &respField)
	return respField, code, err
}
