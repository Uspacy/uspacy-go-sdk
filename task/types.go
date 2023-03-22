package task

const (
	VersionUrl  = "tasks/v1/"
	TransferUrl = "transfer"
)

type TransferTaskOutput struct {
	Status                bool `json:"status"`
	CountOwnerTasks       int  `json:"countOwnerTasks"`
	CountResponsibleTasks int  `json:"countResponsibleTasks"`
}
