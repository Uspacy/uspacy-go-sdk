package smartobjects

const (
	FieldsUrl = "entities/%s/fields"
)

type SmartObjectCreateRequest struct {
	Title           string `json:"title"`
	Type            string `json:"type"`
	Kanban          bool   `json:"kanban"`
	ForAllUsers     bool   `json:"for_all_users"`
	DisplayInMenu   bool   `json:"display_in_menu"`
	ActivitySupport bool   `json:"activity_support"`
	TaskSupport     bool   `json:"task_support"`
}

type CrmSmartObject struct {
	Id              int    `json:"id"`
	Title           string `json:"title"`
	TitleSingular   string `json:"title_singular"`
	TableName       string `json:"table_name"`
	Type            string `json:"type"`
	Sort            string `json:"sort"`
	Kanban          bool   `json:"kanban"`
	ForAllUsers     bool   `json:"for_all_users"`
	DisplayInMenu   bool   `json:"display_in_menu"`
	ActivitySupport bool   `json:"activity_support"`
	TaskSupport     bool   `json:"task_support"`
	Avatar          any    `json:"avatar"`
	CreatedBy       int    `json:"created_by"`
	UpdatedAt       int    `json:"updated_at"`
	CreatedAt       int    `json:"created_at"`
}

type Field struct {
	Name              string  `json:"name"`
	Type              string  `json:"type"`
	Required          bool    `json:"required"`
	Editable          bool    `json:"editable"`
	Show              bool    `json:"show"`
	Hidden            bool    `json:"hidden"`
	Multiple          bool    `json:"multiple"`
	FieldSectionID    string  `json:"field_section_id"`
	SystemField       bool    `json:"system_field"`
	Code              string  `json:"code"`
	Values            []Value `json:"values"`
	EntityReferenceID int     `json:"entity_reference_id"`
}
type Value struct {
	Title    string `json:"title"`
	Value    string `json:"value"`
	Color    string `json:"color"`
	Sort     int    `json:"sort"`
	Selected bool   `json:"selected"`
}
