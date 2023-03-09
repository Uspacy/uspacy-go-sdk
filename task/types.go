package task

const (
	VersionTaskUrl   = "tasks/v1/"
	TransferTasksUrl = "transfer"
)

type TransferTaskOutput struct {
	Status                bool `json:"status"`
	CountOwnerTasks       int  `json:"countOwnerTasks"`
	CountResponsibleTasks int  `json:"countResponsibleTasks"`
}
