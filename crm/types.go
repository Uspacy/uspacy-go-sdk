package crm

import (
	"github.com/Uspacy/uspacy-go-sdk/common"
)

const (
	VersionUrl = "crm/v1/"
)

const (
	EntityUrl          = "entities/%s/"
	FieldsUrl          = "entities/%s/fields/%s/"
	ListsUrl           = "entities/%s/lists/%s"
	FunnelUrl          = "entities/%s/funnel"
	KanbanStageUrl     = "entities/%s/kanban/stage"
	MoveKanbanStageUrl = "entities/%s/%d/move/stage/%s"
)

type Entity int64

const (
	ContactsNum Entity = iota + 1
	LeadsNum
	DealsNum
	CompaniesNum
)

func (e Entity) GetUrl() string {
	uris := map[Entity]string{
		1: "contacts",
		2: "leads",
		3: "deals",
		4: "companies",
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
		Name              string  `json:"name"`
		Code              string  `json:"code"`
		EntityReferenceId string  `json:"entity_reference_id"`
		Type              string  `json:"type"`
		Required          bool    `json:"required"`
		Editable          bool    `json:"editable"`
		Show              bool    `json:"show"`
		Hidden            bool    `json:"hidden"`
		Multiple          bool    `json:"multiple"`
		SystemField       bool    `json:"system_field"`
		BaseField         bool    `json:"base_field"`
		Sort              string  `json:"sort"`
		DefaultValue      string  `json:"default_value"`
		Tooltip           string  `json:"tooltip"`
		Values            []Value `json:"values,omitempty"`
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
)

type (
	EntityCRM struct {
		Id            int     `json:"id"`
		CreatedAt     int     `json:"created_at"`
		UpdatedAt     int     `json:"updated_at"`
		Owner         int     `json:"owner"`
		CreatedBy     int     `json:"created_by"`
		ChangedBy     int     `json:"changed_by,string"`
		Converted     bool    `json:"converted"`
		RelatedEntity bool    `json:"related_entity,omitempty"`
		Title         string  `json:"title"`
		Comments      string  `json:"comments"`
		EntityType    string  `json:"entity_type,omitempty"`
		KanbanStageId string  `json:"kanban_stage_id"`
		Source        []Value `json:"source"`

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
		Messengers []Messanger   `json:"messengers"`
		Phone      []ContactData `json:"phone"`
		Email      []ContactData `json:"email"`
	}

	PersonData struct {
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Patronymic string `json:"patronymic"`
		Position   string `json:"position"`
	}

	Messanger struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Link string `json:"link"`
		Sort string `json:"sort"`
	}

	ContactData struct {
		Id    string `json:"id"`
		Type  string `json:"type"`
		Value string `json:"value"`
		Main  bool   `json:"main"`
		Sort  string `json:"sort"`
	}

	Value struct {
		Title    string `json:"title"`
		Value    string `json:"value"`
		Color    string `json:"color"`
		Sort     string `json:"sort"`
		Selected bool   `json:"selected"`
	}
)
