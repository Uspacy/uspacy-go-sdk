package crm

import (
	"github.com/Uspacy/uspacy-go-sdk/common"
)

const (
	VersionUrl = "crm/v1/"
)

const (
	GetFieldsUrl         = "entities/%s/fields/%s/"
	CreateFunnelUrl      = "entities/%s/funnel"
	CreateKanbanStageUrl = "entities/%s/kanban/stage"
	CreateEntity         = "entities/%s/"
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
)

type (
	Company struct {
		Id            int    `json:"id"`
		CreatedAt     int    `json:"created_at"`
		UpdatedAt     int    `json:"updated_at"`
		Title         string `json:"title"`
		Owner         int    `json:"owner"`
		CreatedBy     int    `json:"created_by"`
		ChangedBy     int    `json:"changed_by,string"`
		Comments      string `json:"comments"`
		Converted     bool   `json:"converted"`
		KanbanStageId string `json:"kanban_stage_id"`
		RelatedEntity bool   `json:"related_entity,omitempty"`
		EntityType    string `json:"entity_type,omitempty"`

		Site    string `json:"site"`
		Address string `json:"address"`

		Messengers   []Messanger   `json:"messengers"`
		Phone        []ContactData `json:"phone"`
		Email        []ContactData `json:"email"`
		Source       []Value       `json:"source"`
		CompanyLabel []Value       `json:"company_label"`
		Contacts     []Contact     `json:"contacts"`

		UtmSource   string `json:"utm_source"`
		UtmMedium   string `json:"utm_medium"`
		UtmCampaign string `json:"utm_campaign"`
		UtmContent  string `json:"utm_content"`
		UtmTerm     string `json:"utm_term"`
	}

	Contact struct {
		Id            int    `json:"id"`
		CreatedAt     int    `json:"created_at"`
		UpdatedAt     int    `json:"updated_at"`
		Title         string `json:"title"`
		Owner         int    `json:"owner"`
		CreatedBy     int    `json:"created_by"`
		ChangedBy     int    `json:"changed_by,string"`
		Comments      string `json:"comments"`
		Converted     bool   `json:"converted"`
		KanbanStageId string `json:"kanban_stage_id"`
		RelatedEntity bool   `json:"related_entity,omitempty"`
		EntityType    string `json:"entity_type,omitempty"`

		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Patronymic string `json:"patronymic"`
		Position   string `json:"position"`

		Messengers   []Messanger   `json:"messengers"`
		Phone        []ContactData `json:"phone"`
		Email        []ContactData `json:"email"`
		Source       []Value       `json:"source"`
		ContactLabel []Value       `json:"contact_label"`
		Companies    []Company     `json:"companies"`

		UtmSource   string `json:"utm_source"`
		UtmMedium   string `json:"utm_medium"`
		UtmCampaign string `json:"utm_campaign"`
		UtmContent  string `json:"utm_content"`
		UtmTerm     string `json:"utm_term"`
	}

	Deal struct {
		Id              int    `json:"id"`
		CreatedAt       int    `json:"created_at"`
		UpdatedAt       int    `json:"updated_at"`
		Title           string `json:"title"`
		Owner           int    `json:"owner"`
		CreatedBy       int    `json:"created_by"`
		ChangedBy       int    `json:"changed_by,string"`
		Comments        string `json:"comments"`
		Converted       bool   `json:"converted"`
		KanbanStageId   int    `json:"kanban_stage_id"`
		KanbanStatus    string `json:"kanban_status"`
		KanbanReasonId  string `json:"kanban_reason_id"`
		AmountOfTheDeal string `json:"amount_of_the_deal"`

		Contacts  []Contact `json:"contacts"`
		Companies []Company `json:"companies"`
		Source    []Value   `json:"source"`
		DealLabel []Value   `json:"deal_label"`

		UtmSource   string `json:"utm_source"`
		UtmMedium   string `json:"utm_medium"`
		UtmCampaign string `json:"utm_campaign"`
		UtmContent  string `json:"utm_content"`
		UtmTerm     string `json:"utm_term"`
	}

	Lead struct {
		Id             int    `json:"id"`
		CreatedAt      int    `json:"created_at"`
		UpdatedAt      int    `json:"updated_at"`
		Title          string `json:"title"`
		Owner          int    `json:"owner"`
		CreatedBy      int    `json:"created_by"`
		ChangedBy      int    `json:"changed_by,string"`
		Comments       string `json:"comments"`
		Converted      bool   `json:"converted"`
		KanbanStageId  int    `json:"kanban_stage_id"`
		KanbanStatus   string `json:"kanban_status"`
		KanbanReasonId int    `json:"kanban_reason_id"`

		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Patronymic  string `json:"patronymic"`
		CompanyName string `json:"company_name"`
		Position    string `json:"position"`

		Messengers []Messanger   `json:"messengers"`
		Phone      []ContactData `json:"phone"`
		Email      []ContactData `json:"email"`
		Source     []Value       `json:"source"`
		LeadLabel  []Value       `json:"lead_label"`

		UtmSource   string `json:"utm_source"`
		UtmMedium   string `json:"utm_medium"`
		UtmCampaign string `json:"utm_campaign"`
		UtmContent  string `json:"utm_content"`
		UtmTerm     string `json:"utm_term"`
	}
)

type (
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
