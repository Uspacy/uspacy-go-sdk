package api

import (
	"encoding/json"
	"fmt"
	"github.com/Uspacy/uspacy-go-sdk/crm"
)

// CreateCanbanStage returns created kanban stage
func (us *Uspacy) CreateCanbanStage(entity string, body interface{}) (crm.KanbanStage, error) {
	var kanbanStage crm.KanbanStage

	responseBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateFunnelUrl, entity)), body)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(responseBody, &kanbanStage)
}
