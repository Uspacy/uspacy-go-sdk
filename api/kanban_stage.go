package api

import (
	"encoding/json"
	"fmt"
	"uspacy-go-sdk/crm"
)

// CreateCanbanStage returns created kanban stage
func (us *Uspacy) CreateCanbanStage(entity string, body interface{}) (crm.KanbanStage, error) {
	var kanbanStage crm.KanbanStage

	responceBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateFunnelUrl, entity)), body)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(responceBody, &kanbanStage)
}
