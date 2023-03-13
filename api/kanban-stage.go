package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// CreateKanbanStage returns created kanban stage
func (us *Uspacy) CreateKanbanStage(entity string, body interface{}) (kanbanStage crm.KanbanStage, err error) {
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.KanbanStageUrl, entity)), body)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(responseBody, &kanbanStage)
}
