package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/v2/task"
)

// CreateTask creates a new task
func (us *Uspacy) CreateTask(ctx context.Context, taskData url.Values) (task.Task, error) {
	resp, err := us.doPostEncodedForm(ctx, us.buildURL(task.VersionUrl, task.TaskUrl), taskData)
	return decodeJSON[task.Task](resp, err)
}

// CreateTaskThroughMap creates a new task through a map
func (us *Uspacy) CreateTaskThroughMap(ctx context.Context, taskData map[string]any, opts ...RequestOption) (task.Task, int, error) {
	resp, code, err := us.doPost(ctx, us.buildURL(task.VersionUrl, task.TaskUrl), taskData, opts...)
	v, err := decodeJSON[task.Task](resp, err)
	return v, code, err
}

// CreateTransferTask creates a new transfer task
func (us *Uspacy) CreateTransferTask(ctx context.Context, body any, opts ...RequestOption) (task.TransferTaskOutput, int, error) {
	resp, code, err := us.doPost(ctx, us.buildURL(task.VersionUrl, task.TransferUrl), body, opts...)
	v, err := decodeJSON[task.TransferTaskOutput](resp, err)
	return v, code, err
}

// PatchTask updates a task and returns it.
func (us *Uspacy) PatchTask(ctx context.Context, taskId int, taskData map[string]any) (task.Task, error) {
	resp, err := us.doPatchEmptyHeaders(ctx, us.buildURL(task.VersionUrl, fmt.Sprintf(task.TaskIdUrl, taskId)), taskData)
	return decodeJSON[task.Task](resp, err)
}

// GetTaskFields returns the task fields.
func (us *Uspacy) GetTaskFields(ctx context.Context) ([]task.Field, error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(task.VersionUrl, task.TaskUrl, task.FieldUrl))
	resp, err := decodeJSON[task.TaskFields](body, err)
	if err != nil {
		return nil, err
	}
	return resp.Fields, nil
}

// GetTasksList returns the tasks matching params.
func (us *Uspacy) GetTasksList(ctx context.Context, params url.Values) (task.TasksList, error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(task.VersionUrl, task.TaskUrl)+"?"+params.Encode())
	return decodeJSON[task.TasksList](body, err)
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
func (us *Uspacy) GetTaskById(ctx context.Context, taskId int, params ...url.Values) (map[string]any, error) {
	urlStr := us.buildURL(task.VersionUrl, fmt.Sprintf(task.TaskIdUrl, taskId))
	if len(params) > 0 && params[0] != nil {
		urlStr += "?" + params[0].Encode()
	}
	body, err := us.doGetEmptyHeaders(ctx, urlStr)
	return decodeJSON[map[string]any](body, err)
}

// GetTaskStagesByGroupId returns task stages by group id
func (us *Uspacy) GetTaskStagesByGroupId(ctx context.Context, groupId int) ([]task.TaskGroupStage, error) {
	params := url.Values{}
	params.Set("groupId", fmt.Sprintf("%d", groupId))
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(task.VersionUrl, task.KanbanStages)+"?"+params.Encode())
	resp, err := decodeJSON[task.TaskGroupStages](body, err)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetTemplateById returns template by id
func (us *Uspacy) GetTemplateById(ctx context.Context, templateId int) (task.Template, error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(task.VersionUrl, task.TemplateUrl, fmt.Sprintf("%d", templateId)))
	return decodeJSON[task.Template](body, err)
}

// CreateTaskStage creates a new task stage
func (us *Uspacy) CreateTaskStage(ctx context.Context, stageData task.TaskGroupStage) (task.TaskGroupStage, int, error) {
	body, code, err := us.doPost(ctx, us.buildURL(task.VersionUrl, task.KanbanStages), stageData)
	v, err := decodeJSON[task.TaskGroupStage](body, err)
	return v, code, err
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
func (us *Uspacy) CreateTaskField(ctx context.Context, fieldData task.Field) (task.Field, int, error) {
	body, code, err := us.doPost(ctx, us.buildURL(task.VersionUrl, task.TaskUrl, task.FieldUrl), fieldData)
	v, err := decodeJSON[task.Field](body, err)
	return v, code, err
}
