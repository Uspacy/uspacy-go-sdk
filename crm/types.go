package crm

import (
	"github.com/Uspacy/uspacy-go-sdk/common"
)

const (
	VersionUrl = "crm/v1"
)

const (
	EntityUrl          = "entities/%s/"
	FieldsUrl          = "entities/%s/fields/%s/"
	ListsUrl           = "entities/%s/lists/%s"
	FunnelUrl          = "entities/%s/funnel"
	KanbanStageUrl     = "entities/%s/kanban/stage"
	MoveKanbanStageUrl = "entities/%s/%d/move/stage/%s"
	ReasonsUrl         = "reasons/%d"
)

type Entity int64

const (
	LeadsNum Entity = iota + 1
	DealsNum
	ContactsNum
	CompaniesNum
	TasksNum
)

func (e Entity) GetUrl() string {
	uris := map[Entity]string{
		ContactsNum:  "contacts",
		LeadsNum:     "leads",
		DealsNum:     "deals",
		CompaniesNum: "companies",
		TasksNum:     "tasks",
	}
	if entity, ok := uris[e]; !ok {
		return "unknown"
	} else {
		return entity
	}
}

type (

	// CRM Fields type.
	Fields struct {
		Data []Field `json:"data"`
	}

	// CRM entity Field
	Field struct {
		Name              string      `json:"name"`
		Code              string      `json:"code"`
		EntityReferenceId interface{} `json:"entity_reference_id"`

		Type         string      `json:"type"`
		Required     bool        `json:"required"`
		Editable     bool        `json:"editable"`
		Show         bool        `json:"show"`
		Hidden       bool        `json:"hidden"`
		Multiple     bool        `json:"multiple"`
		SystemField  bool        `json:"system_field"`
		BaseField    bool        `json:"base_field"`
		Sort         interface{} `json:"sort"` // int
		DefaultValue interface{} `json:"default_value"`
		Tooltip      interface{} `json:"tooltip"`
		Values       []Value     `json:"values,omitempty"`
	}

	// CRM Funnel
	Funnel struct {
		Title      string `json:"title"`
		FunnelCode string `json:"funnel_code"`
		Active     bool   `json:"active"`
		Default    bool   `json:"default"`
		ID         int    `json:"id"`
	}

	// CRM KanbanStage
	KanbanStage struct {
		Title     string `json:"title"`
		StageCode string `json:"stage_code"`
		Sort      string `json:"sort"`
		Color     string `json:"color"`
		ID        int    `json:"id"`
	}
)

type (
	Companies struct {
		Data  []Company    `json:"data"`
		Links common.Links `json:"links"`
		Meta  common.Meta  `json:"meta"`
	}

	Contacts struct {
		Data  []Contact    `json:"data"`
		Links common.Links `json:"links"`
		Meta  common.Meta  `json:"meta"`
	}

	Deals struct {
		Data  []Deal       `json:"data"`
		Links common.Links `json:"links"`
		Meta  common.Meta  `json:"meta"`
	}

	Leads struct {
		Data  []Lead       `json:"data"`
		Links common.Links `json:"links"`
		Meta  common.Meta  `json:"meta"`
	}

	List struct {
		Title    string `json:"title"`
		Value    string `json:"value"`
		Color    string `json:"color"`
		Sort     string `json:"sort"`
		Selected bool   `json:"selected"`
	}
)

type (
	Company struct {
		EntityCRM
		PersonContactData

		Site    string `json:"site"`
		Address string `json:"address"`

		CompanyLabel []Value   `json:"company_label"`
		Contacts     []Contact `json:"contacts"`
	}

	Contact struct {
		EntityCRM
		PersonData
		PersonContactData

		ContactLabel []Value   `json:"contact_label"`
		Companies    []Company `json:"companies"`
	}

	Deal struct {
		EntityCRM
		KanbanCRM

		AmountOfTheDeal string `json:"amount_of_the_deal"`

		Contacts  []Contact `json:"contacts"`
		Companies []Company `json:"companies"`
		DealLabel []Value   `json:"deal_label"`
	}

	Lead struct {
		EntityCRM
		KanbanCRM
		PersonData
		PersonContactData

		CompanyName string  `json:"company_name"`
		LeadLabel   []Value `json:"lead_label"`
	}

	Task struct {
		ID            int         `json:"id"`
		Title         string      `json:"title"`
		Description   string      `json:"description"`
		Type          string      `json:"type"`
		Status        string      `json:"status"`
		CreatedAt     interface{} `json:"created_at"`     // int
		UpdatedAt     interface{} `json:"updated_at"`     // int
		CreatedBy     interface{} `json:"created_by"`     // int
		ResponsibleID interface{} `json:"responsible_id"` // int
		StartTime     interface{} `json:"start_time"`     // int
		EndTime       interface{} `json:"end_time"`       // int
		Deals         []Deal      `json:"deals"`
	}
)

type (
	EntityCRM struct {
		Id            int         `json:"id"`
		CreatedAt     interface{} `json:"created_at"` // int
		UpdatedAt     interface{} `json:"updated_at"` // int
		Owner         interface{} `json:"owner"`      // int
		CreatedBy     interface{} `json:"created_by"` // int
		ChangedBy     interface{} `json:"changed_by"` //`json:"changed_by,string"` // int
		Converted     bool        `json:"converted"`
		RelatedEntity bool        `json:"related_entity,omitempty"`
		Title         string      `json:"title"`
		Comments      string      `json:"comments"`
		EntityType    string      `json:"entity_type,omitempty"`
		KanbanStageId interface{} `json:"kanban_stage_id"` // int
		Source        []Value     `json:"source"`

		UtmSource   string `json:"utm_source"`
		UtmMedium   string `json:"utm_medium"`
		UtmCampaign string `json:"utm_campaign"`
		UtmContent  string `json:"utm_content"`
		UtmTerm     string `json:"utm_term"`
	}

	KanbanCRM struct {
		KanbanStatus   string `json:"kanban_status"`
		KanbanReasonId string `json:"kanban_reason_id"`
	}

	PersonContactData struct {
		/*
			Messengers []Messenger          `json:"messengers"`
			Phone      []common.ContactData `json:"phone"`
			Email      []common.ContactData `json:"email"`
		*/
		Messengers interface{} `json:"messengers"`
		Phone      interface{} `json:"phone"`
		Email      interface{} `json:"email"`
	}

	PersonData struct {
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Patronymic string `json:"patronymic"`
		Position   string `json:"position"`
	}

	Messenger struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Link string `json:"link"`
		Sort string `json:"sort"`
	}

	Value struct {
		Title    string `json:"title"`
		Value    string `json:"value"`
		Color    string `json:"color"`
		Sort     string `json:"sort"`
		Selected bool   `json:"selected"`
	}
)

type (
	Reasons struct {
		Success []interface{} `json:"SUCCESS"`
		Fail    []Fail        `json:"FAIL"`
	}
	Fail struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Sort  int    `json:"sort"`
	}

	FailWrite struct {
		Title string `json:"title"`
		Type  string `json:"type"`
		Sort  int    `json:"sort"`
	}
)
