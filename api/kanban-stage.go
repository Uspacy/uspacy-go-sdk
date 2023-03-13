package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// CreateKanbanStage returns created kanban stage
func (us *Uspacy) CreateKanbanStage(entityType crm.Entity, body interface{}) (kanbanStage crm.KanbanStage, err error) {
	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionCRMUrl, fmt.Sprintf(crm.KanbanStageUrl, entityType.GetUrl(crm.EntityType))), body)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(responseBody, &kanbanStage)
}
