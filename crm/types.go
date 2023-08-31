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
	CreateFieldUrl     = "entities/%s/fields"
	ListsUrl           = "entities/%s/lists/%s"
	FunnelUrl          = "entities/%s/funnel"
	KanbanStageUrl     = "entities/%s/kanban/stage/%v"
	StageByFunnelIdUrl = "?funnel_id=%d"
	MoveKanbanStageUrl = "entities/%s/%d/move/stage/%s"
	ReasonsUrl         = "reasons/%d"
	TaskUrl            = "static/tasks/%s"
	CallUrl            = "events/call"
	ProductIdUrl       = "static/products/%s"
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

type FunnelsById []struct {
	ID         int           `json:"id"`
	Title      string        `json:"title"`
	FunnelCode string        `json:"funnel_code"`
	Default    bool          `json:"default"`
	Active     bool          `json:"active"`
	Stages     []KanbanStage `json:"stages"`
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

	// CRM KanbanStages List

	KanbanStages struct {
		Data []KanbanStage `json:"data"`
	}

	// CRM KanbanStage
	KanbanStage struct {
		FunnelStageBase
		ID          int      `json:"id"`
		SystemStage bool     `json:"system_stage,omitempty"`
		StageCode   string   `json:"stage_code"`
		Reasons     []Reason `json:"reasons,omitempty"`
	}

	FunnelStage struct {
		FunnelStageBase
		FunnelId int `json:"funnel_id"`
	}

	FunnelStageBase struct {
		Title string `json:"title"`
		Color string `json:"color"`
		Sort  string `json:"sort"`
	}
)

type (
	CRMEntity struct {
		Data []struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
		} `json:"data"`
		Links common.Links `json:"links"`
		Meta  common.Meta  `json:"meta"`
	}

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
		Title    string      `json:"title"`
		Value    string      `json:"value"`
		Color    string      `json:"color"`
		Sort     interface{} `json:"sort"`
		Selected bool        `json:"selected"`
	}
)

type (
	Company struct {
		EntityCRM
		PersonContactData

		Site    string `json:"site"`
		Address string `json:"address"`

		//CompanyLabel []Value   `json:"company_label"`
		//Contacts     []Contact `json:"contacts"`
	}

	Contact struct {
		EntityCRM
		PersonData
		PersonContactData

		//ContactLabel []Value   `json:"contact_label"`
		//Companies    []Company `json:"companies"`
	}

	Deal struct {
		EntityCRM
		KanbanCRM

		AmountOfTheDeal interface{} `json:"amount_of_the_deal"`

		//Contacts  []Contact `json:"contacts"`
		//Companies []Company `json:"companies"`
		//DealLabel []Value   `json:"deal_label"`
	}

	AmountOfTheDealData struct {
		Currency string `json:"currency"`
		Value    string `json:"value"`
	}

	Lead struct {
		EntityCRM
		KanbanCRM
		PersonData
		PersonContactData

		CompanyName string `json:"company_name"`
		//LeadLabel   []Value `json:"lead_label"`
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
		//Deals         []Deal      `json:"deals"`
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
		Converted     interface{} `json:"converted"`  // bool
		RelatedEntity bool        `json:"related_entity,omitempty"`
		Title         string      `json:"title"`
		Comments      string      `json:"comments"`
		EntityType    string      `json:"entity_type,omitempty"`
		KanbanStageId interface{} `json:"kanban_stage_id"` // int
		Source        interface{} `json:"source"`          // []Value

		UtmSource   string `json:"utm_source"`
		UtmMedium   string `json:"utm_medium"`
		UtmCampaign string `json:"utm_campaign"`
		UtmContent  string `json:"utm_content"`
		UtmTerm     string `json:"utm_term"`
	}

	KanbanFailReasonCRM struct {
		ReasonId string `json:"reason_id"`
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
		Id   string      `json:"id"`
		Name string      `json:"name"`
		Link string      `json:"link"`
		Sort interface{} `json:"sort"`
	}

	Value struct {
		Title    string      `json:"title"`
		Value    string      `json:"value"`
		Color    string      `json:"color"`
		Sort     interface{} `json:"sort"`
		Selected bool        `json:"selected"`
	}
)

type (
	Reasons struct {
		Success []interface{} `json:"SUCCESS"`
		Fail    []Reason      `json:"FAIL"`
	}
	Reason struct {
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

type (
	Call struct {
		IntegrationCode string `json:"integration_code"`
		ExternalID      int    `json:"external_id"`
		TmpID           int    `json:"tmp_id"`
		TaskID          any    `json:"task_id"`
		ContactID       any    `json:"contact_id"`
		CompanyID       any    `json:"company_id"`
		LeadID          any    `json:"lead_id"`
		DealID          any    `json:"deal_id"`
		EntityTable     string `json:"entity_table"`
		Employees       []int  `json:"employees"`
		Contacts        []int  `json:"contacts"`
		Companies       []int  `json:"companies"`
		Subject         string `json:"subject"`
		CallType        string `json:"call_type"`
		EndedCallStatus any    `json:"ended_call_status"`
		From            int64  `json:"from"`
		To              any    `json:"to"`
		BeginTime       int    `json:"begin_time"`
		EndTime         int    `json:"end_time"`
		Duration        int    `json:"duration"`
		CallRecordLink  any    `json:"call_record_link"`
		Note            string `json:"note"`
	}
)

type (
	Products struct {
		Data  []Product    `json:"data"`
		Links common.Links `json:"links"`
	}

	Product struct {
		ID                int             `json:"id"`
		ProductCategoryID any             `json:"product_category_id"`
		MeasurementUnitID int             `json:"measurement_unit_id"`
		Title             string          `json:"title"`
		Article           any             `json:"article"`
		Type              string          `json:"type"`
		IsActive          int             `json:"is_active"`
		Availability      string          `json:"availability"`
		Quantity          any             `json:"quantity"`
		ReservedQuantity  any             `json:"reserved_quantity"`
		Description       any             `json:"description"`
		Comment           any             `json:"comment"`
		Link              any             `json:"link"`
		CreatedAt         int             `json:"created_at"`
		UpdatedAt         int             `json:"updated_at"`
		ProductCategory   ProductCategory `json:"product_category"`
		MeasurementUnit   MeasurementUnit `json:"measurement_unit"`
		Prices            []Prices        `json:"prices"`
		Files             []File          `json:"files"`
	}

	ProductCategory struct {
		ID       int    `json:"id"`
		ParentID int    `json:"parent_id"`
		Name     string `json:"name"`
		IsActive int    `json:"is_active"`
	}

	MeasurementUnit struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Abbr      string `json:"abbr"`
		IsDefault int    `json:"is_default"`
	}

	Prices struct {
		Price         int    `json:"price"`
		Currency      string `json:"currency"`
		IsDefault     int    `json:"is_default"`
		IsTaxIncluded int    `json:"is_tax_included"`
		TaxID         int    `json:"tax_id"`
		Tax           Tax    `json:"tax"`
	}

	File struct {
		ID               int    `json:"id"`
		EntityType       string `json:"entityType"`
		EntityID         string `json:"entityId"`
		UploadID         string `json:"uploadId"`
		OriginalFilename string `json:"originalFilename"`
		LastModified     int    `json:"lastModified"`
		Size             int    `json:"size"`
		URL              string `json:"url"`
	}

	Tax struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Rate      int    `json:"rate"`
		IsActive  int    `json:"is_active"`
		CreatedAt int    `json:"created_at"`
		UpdatedAt int    `json:"updated_at"`
	}
)
