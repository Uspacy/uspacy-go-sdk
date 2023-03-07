package crm

const (
	VersionUrl = "crm/v1/"
)

const (
	//GetFieldsUrl         = "entity/%s/fields/"
	GetFieldsUrl         = "entities/%s/fields/%s/"
	CreateFunnelUrl      = "entities/%s/funnel"
	CreateKanbanStageUrl = "entities/%s/kanban/stage"
)

type (

	// CRM entity Field
	Field struct {
		Name              string      `json:"name"`
		Code              string      `json:"code"`
		EntityReferenceId interface{} `json:"entity_reference_id"`
		Type              string      `json:"type"`
		Required          bool        `json:"required"`
		Editable          bool        `json:"editable"`
		Show              bool        `json:"show"`
		Hidden            bool        `json:"hidden"`
		Multiple          bool        `json:"multiple"`
		SystemField       bool        `json:"system_field"`
		BaseField         bool        `json:"base_field"`
		Sort              string      `json:"sort"`
		DefaultValue      string      `json:"default_value"`
		Tooltip           string      `json:"tooltip"`
		Values            []struct {
			Title    string `json:"title"`
			Value    string `json:"value"`
			Color    string `json:"color"`
			Sort     string `json:"sort"`
			Selected bool   `json:"selected"`
		} `json:"values,omitempty"`
	}

	// CRM Fields type.
	Fields struct {
		Data []Field `json:"data"`
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
