package crm

const (
	VersionUrl = "crm/v1/"
)

const (
	FieldsUrl      = "entities/%s/fields/%s/"
	ListsUrl       = "entities/%s/lists/%s/"
	FunnelUrl      = "entities/%s/funnel"
	KanbanStageUrl = "entities/%s/kanban/stage"
	EntityUrl      = "entities/%s/"
)

type (
	Entity int64
)

const (
	Contacts Entity = iota + 1
	Leads
	Deals
	Companies
)

var (
	EntityType = map[Entity]string{
		1: "contacts",
		2: "leads",
		3: "deals",
		4: "companies",
	}

	StatusId = map[Entity]string{
		1: "SOURCE",
		2: "DEAL_TYPE",
		3: "DEAL_STAGE",
		4: "CONTACT_TYPE",
		5: "COMPANY_TYPE",
		6: "INDUSTRY",
		7: "EMPLOYEES",
	}

	ListValue = map[Entity]string{
		1: "source",
		2: "deal_label",
		3: "DEAL_STAGE", ///!!!!!!!!!!!!!!!!!!
		4: "contact_label",
		5: "company_label",
		6: "industry_label",
		7: "employees_label",
	}
)

func (e Entity) GetUrl(list map[Entity]string) string {
	uris := list
	if entity, ok := uris[e]; !ok {
		return "unknown"
	} else {
		return entity
	}
}

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
		Values            []List      `json:"values,omitempty"`
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

	// CRM EntityResponse, returns when entity object was created successfully
	EntityResponse struct {
		Id               int           `json:"id"`
		CreatedAt        int           `json:"created_at"`
		UpdatedAt        int           `json:"updated_at"`
		Title            string        `json:"title"`
		Owner            int           `json:"owner"`
		CreatedBy        int           `json:"created_by"`
		ChangedBy        string        `json:"changed_by"`
		Converted        bool          `json:"converted"`
		FirstName        string        `json:"first_name"`
		LastName         string        `json:"last_name"`
		Patronymic       string        `json:"patronymic"`
		Companies        []Company     `json:"companies"`
		Position         string        `json:"position"`
		UtmSource        string        `json:"utm_source"`
		UtmMedium        string        `json:"utm_medium"`
		UtmCampaign      string        `json:"utm_campaign"`
		UtmContent       string        `json:"utm_content"`
		UtmTerm          string        `json:"utm_term"`
		Messengers       string        `json:"messengers"`
		Phone            string        `json:"phone"`
		Email            string        `json:"email"`
		Comments         string        `json:"comments"`
		Source           []Source      `json:"source"`
		ContactLabel     []Label       `json:"contact_label"`
		Ownercopy        string        `json:"ownercopy"`
		LastNamecopy     string        `json:"last_namecopy"`
		Companiescopy    []interface{} `json:"companiescopy"`
		FirstNamecopy    string        `json:"first_namecopy"`
		Titlecopy        string        `json:"titlecopy"`
		Patronymiccopy   string        `json:"patronymiccopy"`
		Positioncopy     string        `json:"positioncopy"`
		UtmSourcecopy    string        `json:"utm_sourcecopy"`
		UtmMediumcopy    string        `json:"utm_mediumcopy"`
		UtmCampaigncopy  string        `json:"utm_campaigncopy"`
		UtmContentcopy   string        `json:"utm_contentcopy"`
		UtmTermcopy      string        `json:"utm_termcopy"`
		Emailcopy        string        `json:"emailcopy"`
		Messengerscopy   string        `json:"messengerscopy"`
		Commentscopy     string        `json:"commentscopy"`
		Phonecopy        string        `json:"phonecopy"`
		Sourcecopy       string        `json:"sourcecopy"`
		ContactLabelcopy string        `json:"contact_labelcopy"`
		KanbanStageId    string        `json:"kanban_stage_id"`
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
		Id               int         `json:"id"`
		CreatedAt        int         `json:"created_at"`
		UpdatedAt        int         `json:"updated_at"`
		Title            string      `json:"title"`
		Owner            int         `json:"owner"`
		CreatedBy        int         `json:"created_by"`
		ChangedBy        string      `json:"changed_by"`
		Converted        bool        `json:"converted"`
		CompanyName      string      `json:"company_name"`
		UtmSource        string      `json:"utm_source"`
		UtmMedium        string      `json:"utm_medium"`
		UtmCampaign      string      `json:"utm_campaign"`
		UtmContent       string      `json:"utm_content"`
		UtmTerm          string      `json:"utm_term"`
		Messengers       []Messenger `json:"messengers"`
		Phone            []Contact   `json:"phone"`
		Email            []Contact   `json:"email"`
		Comments         string      `json:"comments"`
		Source           string      `json:"source"`
		CompanyLabel     string      `json:"company_label"`
		Site             string      `json:"site"`
		Address          string      `json:"address"`
		Contacts         string      `json:"contacts"`
		Ownercopy        string      `json:"ownercopy"`
		Titlecopy        string      `json:"titlecopy"`
		UtmSourcecopy    string      `json:"utm_sourcecopy"`
		CompanyNamecopy  string      `json:"company_namecopy"`
		UtmMediumcopy    string      `json:"utm_mediumcopy"`
		UtmCampaigncopy  string      `json:"utm_campaigncopy"`
		UtmContentcopy   string      `json:"utm_contentcopy"`
		UtmTermcopy      string      `json:"utm_termcopy"`
		Messengerscopy   string      `json:"messengerscopy"`
		Phonecopy        string      `json:"phonecopy"`
		Emailcopy        string      `json:"emailcopy"`
		Commentscopy     string      `json:"commentscopy"`
		CompanyLabelcopy string      `json:"company_labelcopy"`
		Sourcecopy       string      `json:"sourcecopy"`
		Sitecopy         string      `json:"sitecopy"`
		Addresscopy      string      `json:"addresscopy"`
		Contactscopy     string      `json:"contactscopy"`
		RelatedEntity    bool        `json:"related_entity"`
		EntityType       string      `json:"entity_type"`
		KanbanStageId    string      `json:"kanban_stage_id"`
	}

	Source struct {
		Title    string `json:"title"`
		Value    string `json:"value"`
		Color    string `json:"color"`
		Sort     string `json:"sort"`
		Selected bool   `json:"selected"`
	}

	Label struct {
		Title    string `json:"title"`
		Value    string `json:"value"`
		Color    string `json:"color"`
		Sort     string `json:"sort"`
		Selected bool   `json:"selected"`
	}

	Contact struct {
		Id    string `json:"id"`
		Type  string `json:"type"`
		Value string `json:"value"`
		Main  bool   `json:"main"`
		Sort  string `json:"sort"`
	}

	Messenger struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Link string `json:"link"`
		Sort string `json:"sort"`
	}
)
