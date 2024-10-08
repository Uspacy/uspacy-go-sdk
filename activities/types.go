package activities

import (
	"encoding/json"

	"github.com/Uspacy/uspacy-go-sdk/common"
)

const (
	VersionUrl = "activities/v1"
)

const (
	ActivitiesUrl = "activities"
	ActivityUrl   = "entities/%s"
	MassDeletion  = "entities/mass_deletion"
)

type ActivitiesList struct {
	Activities []Activity   `json:"data"`
	Links      common.Links `json:"links"`
	Meta       common.Meta  `json:"meta"`
}

type Activity struct {
	Id            int             `json:"id"`
	CreatedBy     int             `json:"created_by"`
	Title         string          `json:"title"`
	Type          string          `json:"type"`
	Status        string          `json:"status"`
	Priority      string          `json:"priority"`
	Description   string          `json:"description"`
	StartTime     int             `json:"start_time,omitempty"`
	EndTime       int             `json:"end_time"`
	ResponsibleID int             `json:"responsible_id"`
	CompanyID     any             `json:"company_id"`
	CreatedAt     int             `json:"created_at,omitempty"`
	UpdatedAt     int             `json:"updated_at,omitempty"`
	ClosedBy      int             `json:"closed_by,omitempty"`
	FirstClosedBy any             `json:"first_closed_by,omitempty"`
	ClosedAt      any             `json:"closed_at,omitempty"`
	FirstClosedAt any             `json:"first_closed_at,omitempty"`
	CrmID         any             `json:"crm_id"`
	Participants  []int           `json:"participants,omitempty"`
	CrmEntities   json.RawMessage `json:"crm_entities,omitempty"`
}

type CrmEntities map[string]CrmEntity

type CrmEntity struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Owner     int    `json:"owner"`
	TableName string `json:"table_name"`
}

type MassDeletionBody struct {
	All       bool     `json:"all,omitempty"`
	EntityIds []string `json:"entity_ids,omitempty"`
	ExceptIds []any    `json:"except_ids,omitempty"`
}
