package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/v2/task"
)

// CreateTask creates a new task
func (us *Uspacy) CreateTask(ctx context.Context, taskData url.Values, opts ...RequestOption) (task.Task, error) {
	resp, err := us.doPostEncodedForm(ctx, us.buildURL(task.VersionUrl, task.TaskUrl), taskData, opts...)
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
func (us *Uspacy) PatchTask(ctx context.Context, taskId int, taskData map[string]any, opts ...RequestOption) (task.Task, error) {
	resp, err := us.doPatch(ctx, us.buildURL(task.VersionUrl, fmt.Sprintf(task.TaskIdUrl, taskId)), taskData, opts...)
	return decodeJSON[task.Task](resp, err)
}

// GetTaskFields returns the task fields.
func (us *Uspacy) GetTaskFields(ctx context.Context, opts ...RequestOption) ([]task.Field, error) {
	body, err := us.doGet(ctx, us.buildURL(task.VersionUrl, task.TaskUrl, task.FieldUrl), opts...)
	resp, err := decodeJSON[task.TaskFields](body, err)
	if err != nil {
		return nil, err
	}
	return resp.Fields, nil
}

// GetTasksList returns the tasks matching params.
func (us *Uspacy) GetTasksList(ctx context.Context, params url.Values, opts ...RequestOption) (task.TasksList, error) {
	body, err := us.doGet(ctx, us.buildURL(task.VersionUrl, task.TaskUrl)+"?"+params.Encode(), opts...)
	return decodeJSON[task.TasksList](body, err)
}

// GetTasksWithFilters returns the tasks matching params, each as a map of its raw JSON fields.
func (us *Uspacy) GetTasksWithFilters(ctx context.Context, params url.Values, opts ...RequestOption) ([]map[string]any, error) {
	body, err := us.doGet(ctx, us.buildURL(task.VersionUrl, task.TaskUrl)+"?"+params.Encode(), opts...)
	resp, err := decodeJSON[struct {
		Data []map[string]any `json:"data"`
	}](body, err)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetTaskById returns a task by ID as a map of its raw JSON fields. params is optional;
// pass nil for none.
func (us *Uspacy) GetTaskById(ctx context.Context, taskId int, params url.Values, opts ...RequestOption) (map[string]any, error) {
	urlStr := us.buildURL(task.VersionUrl, fmt.Sprintf(task.TaskIdUrl, taskId))
	if len(params) != 0 {
		urlStr += "?" + params.Encode()
	}
	body, err := us.doGet(ctx, urlStr, opts...)
	return decodeJSON[map[string]any](body, err)
}

// GetTaskStagesByGroupId returns task stages by group id
func (us *Uspacy) GetTaskStagesByGroupId(ctx context.Context, groupId int, opts ...RequestOption) ([]task.TaskGroupStage, error) {
	params := url.Values{}
	params.Set("groupId", fmt.Sprintf("%d", groupId))
	body, err := us.doGet(ctx, us.buildURL(task.VersionUrl, task.KanbanStages)+"?"+params.Encode(), opts...)
	resp, err := decodeJSON[task.TaskGroupStages](body, err)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetTemplateById returns template by id
func (us *Uspacy) GetTemplateById(ctx context.Context, templateId int, opts ...RequestOption) (task.Template, error) {
	body, err := us.doGet(ctx, us.buildURL(task.VersionUrl, task.TemplateUrl, fmt.Sprintf("%d", templateId)), opts...)
	return decodeJSON[task.Template](body, err)
}

// CreateTaskStage creates a new task stage
func (us *Uspacy) CreateTaskStage(ctx context.Context, stageData task.TaskGroupStage, opts ...RequestOption) (task.TaskGroupStage, int, error) {
	body, code, err := us.doPost(ctx, us.buildURL(task.VersionUrl, task.KanbanStages), stageData, opts...)
	v, err := decodeJSON[task.TaskGroupStage](body, err)
	return v, code, err
}

// DeleteTaskStage deletes a task stage
func (us *Uspacy) DeleteTaskStage(ctx context.Context, stageId int, opts ...RequestOption) (err error) {
	_, err = us.doDelete(ctx, us.buildURL(task.VersionUrl, task.KanbanStages, fmt.Sprintf("%d", stageId)), nil, opts...)
	return err
}

// TaskStatusReady marks task as ready
func (us *Uspacy) TaskStatusReady(ctx context.Context, taskId int, opts ...RequestOption) (err error) {
	_, err = us.doPatch(ctx, us.buildURL(task.VersionUrl, fmt.Sprintf(task.TaskIdUrl, taskId), task.TaskStatusReady), nil, opts...)
	return err
}

// CreateTaskField creates a new task field
func (us *Uspacy) CreateTaskField(ctx context.Context, fieldData task.Field, opts ...RequestOption) (task.Field, int, error) {
	body, code, err := us.doPost(ctx, us.buildURL(task.VersionUrl, task.TaskUrl, task.FieldUrl), fieldData, opts...)
	v, err := decodeJSON[task.Field](body, err)
	return v, code, err
}
