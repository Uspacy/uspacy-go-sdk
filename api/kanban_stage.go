package api

import (
	"encoding/json"
	"fmt"
	"io"
	"uspacy-go-sdk/crm"
)

// CreateCanbanStage returns created kanban stage
func (us *Uspacy) CreateCanbanStage(entity string, body io.Reader) (crm.KanbanStage, error) {
	var kanbanStage crm.KanbanStage

	requestBody, err := us.doPostEmptyHeaders(buildURL(mainHost, crm.VersionUrl, fmt.Sprintf(crm.CreateFunnelUrl, entity)), body)
	if err != nil {
		return kanbanStage, err
	}
	return kanbanStage, json.Unmarshal(requestBody, &kanbanStage)
}
