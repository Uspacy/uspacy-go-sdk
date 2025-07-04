package task

const (
	VersionUrl   = "tasks/v1"
	TaskUrl      = "tasks"
	KanbanStages = "stages"
	TaskIdUrl    = "tasks/%d"
	TransferUrl  = "transfer"
	TemplateUrl  = "templates/recurring"
)

type TasksList struct {
	Data []any `json:"data"`
	Meta Meta  `json:"meta"`
}

type Meta struct {
	CurrentPage int `json:"currentPage"`
	From        int `json:"from"`
	LastPage    int `json:"lastPage"`
	PerPage     int `json:"perPage"`
	To          int `json:"to"`
	Total       int `json:"total"`
}

type (
	Task struct {
		ID              string   `json:"id"`
		ParentID        any      `json:"parentId"`
		Title           string   `json:"title"`
		Body            string   `json:"body"`
		CreatedBy       string   `json:"createdBy"`
		CreatedDate     int      `json:"createdDate"`
		ClosedDate      any      `json:"closedDate"`
		Deadline        any      `json:"deadline"`
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
		KanbanStageID   any      `json:"kanbanStageId"`
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

type TaskGroupStage struct {
	Id      int    `json:"id,string"`
	Title   string `json:"title"`
	Color   string `json:"color"`
	Sort    int    `json:"sort"`
	GroupId int    `json:"groupId"`
}

type TaskGroupStages struct {
	Data []TaskGroupStage `json:"data"`
}

type Template struct {
	Id                 int       `json:"id,string"`
	ParentID           string    `json:"parentId"`
	Title              string    `json:"title"`
	SetterID           string    `json:"setterId"`
	CreatedBy          string    `json:"createdBy"`
	CreatedDate        int       `json:"createdDate"`
	DepartmentID       string    `json:"departmentId"`
	GroupID            string    `json:"groupId"`
	ClosedBy           string    `json:"closedBy"`
	ResponsibleID      string    `json:"responsibleId"`
	Deadline           int       `json:"deadline"`
	ClosedDate         string    `json:"closedDate"`
	Body               string    `json:"body"`
	Status             string    `json:"status"`
	Priority           string    `json:"priority"`
	ResultCommentID    string    `json:"resultCommentId"`
	RequiredResult     bool      `json:"requiredResult"`
	AcceptResult       bool      `json:"acceptResult"`
	TimeTracking       bool      `json:"timeTracking"`
	TimeEstimate       int       `json:"timeEstimate"`
	Template           bool      `json:"template"`
	TemplateID         int       `json:"templateId"`
	Delegation         bool      `json:"delegation"`
	BasicTask          bool      `json:"basicTask"`
	TaskType           string    `json:"taskType"`
	AccomplicesIds     []any     `json:"accomplicesIds"`
	AuditorsIds        []any     `json:"auditorsIds"`
	FileIds            string    `json:"fileIds"`
	GroupKanbanStageID any       `json:"groupKanbanStageId"`
	Scheduler          Scheduler `json:"scheduler"`
	TableName          string    `json:"tableName"`
	KanbanStageID      string    `json:"kanbanStageId"`
}
type Scheduler struct {
	TaskID           int    `json:"taskId"`
	Active           bool   `json:"active"`
	NextRun          int    `json:"nextRun"`
	DateStart        string `json:"dateStart"`
	HourStart        int    `json:"hourStart"`
	DateStop         int    `json:"dateStop"`
	TimezoneOffset   int    `json:"timezoneOffset"`
	DeadlineDay      any    `json:"deadlineDay"`
	DeadlineHour     any    `json:"deadlineHour"`
	CurrentIteration int    `json:"currentIteration"`
	Iteration        int    `json:"iteration"`
	Period           string `json:"period"`
	Every            int    `json:"every"`
	DayOfWeek        any    `json:"dayOfWeek"`
	DayOfMonth       any    `json:"dayOfMonth"`
	WeekOfMonth      any    `json:"weekOfMonth"`
	StoppedByForce   int    `json:"stoppedByForce"`
	CreatedAt        int    `json:"createdAt"`
	UpdatedAt        int    `json:"updatedAt"`
	ActivationLimit  bool   `json:"activationLimit"`
}
