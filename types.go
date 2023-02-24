package uspacy_go_sdk

type (

	// CRM entity field
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
)
