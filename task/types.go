package task

const (
	VersionUrl  = "tasks/v1"
	TaskUrl     = "tasks"
	TaskIdUrl   = "tasks/%d"
	TransferUrl = "transfer"
)

type (
	Task struct {
		ID              string   `json:"id"`
		ParentID        string   `json:"parentId"`
		Title           string   `json:"title"`
		Body            string   `json:"body"`
		CreatedBy       string   `json:"createdBy"`
		CreatedDate     int      `json:"createdDate"`
		ClosedDate      int      `json:"closedDate"`
		Deadline        int      `json:"deadline"`
		DepartmentID    string   `json:"departmentId"`
		GroupID         string   `json:"groupId"`
		ClosedBy        string   `json:"closedBy"`
		ResponsibleID   string   `json:"responsibleId"`
		Status          string   `json:"status"`
		Priority        string   `json:"priority"`
		ResultCommentID string   `json:"resultCommentId"`
		RequiredResult  bool     `json:"requiredResult"`
		AcceptResult    bool     `json:"acceptResult"`
		AccomplicesIds  []string `json:"accomplicesIds"`
		AuditorsIds     []string `json:"auditorsIds"`
		KanbanStageID   string   `json:"kanbanStageId"`
	}

	TransferTaskOutput struct {
		Status                bool `json:"status"`
		CountOwnerTasks       int  `json:"countOwnerTasks"`
		CountResponsibleTasks int  `json:"countResponsibleTasks"`
	}
)
