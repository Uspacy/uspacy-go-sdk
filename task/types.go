package task

import "github.com/Uspacy/uspacy-go-sdk/common"

const (
	VersionUrl  = "tasks/v1"
	TaskUrl     = "tasks"
	TaskIdUrl   = "tasks/%d"
	TransferUrl = "transfer"
)

type TasksList struct {
	Data []any       `json:"data"`
	Meta common.Meta `json:"meta"`
}

type (
	Task struct {
		ID              string      `json:"id"`
		ParentID        string      `json:"parentId"`
		Title           string      `json:"title"`
		Body            string      `json:"body"`
		CreatedBy       string      `json:"createdBy"`
		CreatedDate     int         `json:"createdDate"`
		ClosedDate      int         `json:"closedDate"`
		Deadline        int         `json:"deadline"`
		DepartmentID    string      `json:"departmentId"`
		GroupID         string      `json:"groupId"`
		ClosedBy        string      `json:"closedBy"`
		ResponsibleID   string      `json:"responsibleId"`
		Status          string      `json:"status"`
		Priority        string      `json:"priority"`
		ResultCommentID string      `json:"resultCommentId"`
		RequiredResult  bool        `json:"requiredResult"`
		AcceptResult    bool        `json:"acceptResult"`
		AccomplicesIds  []string    `json:"accomplicesIds"`
		AuditorsIds     []string    `json:"auditorsIds"`
		KanbanStageID   interface{} `json:"kanbanStageId"`
	}

	TransferTaskOutput struct {
		Status                bool `json:"status"`
		CountOwnerTasks       int  `json:"countOwnerTasks"`
		CountResponsibleTasks int  `json:"countResponsibleTasks"`
	}
)
type TaskFields struct {
	Fields []Field `json:"data"`
}
type Field struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	Type           string `json:"type"`
	Required       bool   `json:"required"`
	Editable       bool   `json:"editable"`
	Show           bool   `json:"show"`
	Hidden         bool   `json:"hidden"`
	Multiple       bool   `json:"multiple"`
	SystemField    bool   `json:"systemField"`
	Sort           string `json:"sort"`
	DefaultValue   string `json:"defaultValue"`
	ListUUID       string `json:"listUuid"`
	FieldSectionID string `json:"fieldSectionId"`
	Tooltip        string `json:"tooltip"`
}
