package crm

const (
	VersionCRMUrl = "crm/v1/"
)

const (
	FieldsUrl      = "entities/%s/fields/%s"
	ListsUrl       = "entities/%s/lists/%s"
	FunnelUrl      = "entities/%s/funnel"
	KanbanStageUrl = "entities/%s/kanban/stage"
	EntityUrl      = "entities/%s"
)

type (
	Entity int64
	Lists  int64
)

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
)

type Companies struct {
	Data  []Company `json:"data"`
	Links Links     `json:"links"`
	Meta  Meta      `json:"meta"`
}

type Contacts struct {
	Data  []Contact `json:"data"`
	Links Links     `json:"links"`
	Meta  Meta      `json:"meta"`
}

type Deals struct {
	Data  []Deal `json:"data"`
	Links Links  `json:"links"`
	Meta  Meta   `json:"meta"`
}

type Leads struct {
	Data  []Lead `json:"data"`
	Links Links  `json:"links"`
	Meta  Meta   `json:"meta"`
}

type (
	Company struct {
		Id               int             `json:"id"`
		CreatedAt        int             `json:"created_at"`
		UpdatedAt        int             `json:"updated_at"`
		Title            string          `json:"title"`
		Owner            int             `json:"owner"`
		CreatedBy        int             `json:"created_by"`
		ChangedBy        string          `json:"changed_by"`
		Converted        bool            `json:"converted"`
		CompanyName      string          `json:"company_name"`
		UtmSource        string          `json:"utm_source"`
		UtmMedium        string          `json:"utm_medium"`
		UtmCampaign      string          `json:"utm_campaign"`
		UtmContent       string          `json:"utm_content"`
		UtmTerm          string          `json:"utm_term"`
		Messengers       []Messenger     `json:"messengers"`
		Phone            []Phone         `json:"phone"`
		Email            []Email         `json:"email"`
		Comments         string          `json:"comments"`
		Source           []GeneralSource `json:"source"`
		CompanyLabel     []Label         `json:"company_label"`
		Site             string          `json:"site"`
		Address          string          `json:"address"`
		Contacts         []interface{}   `json:"contacts"`
		Ownercopy        string          `json:"ownercopy"`
		Titlecopy        string          `json:"titlecopy"`
		UtmSourcecopy    string          `json:"utm_sourcecopy"`
		CompanyNamecopy  string          `json:"company_namecopy"`
		UtmMediumcopy    string          `json:"utm_mediumcopy"`
		UtmCampaigncopy  string          `json:"utm_campaigncopy"`
		UtmContentcopy   string          `json:"utm_contentcopy"`
		UtmTermcopy      string          `json:"utm_termcopy"`
		Messengerscopy   string          `json:"messengerscopy"`
		Phonecopy        string          `json:"phonecopy"`
		Emailcopy        string          `json:"emailcopy"`
		Commentscopy     string          `json:"commentscopy"`
		CompanyLabelcopy string          `json:"company_labelcopy"`
		Sourcecopy       string          `json:"sourcecopy"`
		Sitecopy         string          `json:"sitecopy"`
		Addresscopy      string          `json:"addresscopy"`
		Contactscopy     []interface{}   `json:"contactscopy"`
		KanbanStageId    string          `json:"kanban_stage_id"`
	}

	Contact struct {
		Id         int         `json:"id"`
		CreatedAt  int         `json:"created_at"`
		UpdatedAt  int         `json:"updated_at"`
		Title      string      `json:"title"`
		Owner      int         `json:"owner"`
		CreatedBy  int         `json:"created_by"`
		ChangedBy  interface{} `json:"changed_by"`
		Converted  bool        `json:"converted"`
		FirstName  string      `json:"first_name"`
		LastName   string      `json:"last_name"`
		Patronymic string      `json:"patronymic"`
		Companies  []struct {
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
			Phone            []Phone     `json:"phone"`
			Email            []Email     `json:"email"`
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
		} `json:"companies"`
		Position         string          `json:"position"`
		UtmSource        string          `json:"utm_source"`
		UtmMedium        string          `json:"utm_medium"`
		UtmCampaign      string          `json:"utm_campaign"`
		UtmContent       string          `json:"utm_content"`
		UtmTerm          string          `json:"utm_term"`
		Messengers       interface{}     `json:"messengers"`
		Phone            interface{}     `json:"phone"`
		Email            interface{}     `json:"email"`
		Comments         string          `json:"comments"`
		Source           []GeneralSource `json:"source"`
		ContactLabel     []Label         `json:"contact_label"`
		Ownercopy        string          `json:"ownercopy"`
		LastNamecopy     string          `json:"last_namecopy"`
		Companiescopy    []interface{}   `json:"companiescopy"`
		FirstNamecopy    string          `json:"first_namecopy"`
		Titlecopy        string          `json:"titlecopy"`
		Patronymiccopy   string          `json:"patronymiccopy"`
		Positioncopy     string          `json:"positioncopy"`
		UtmSourcecopy    string          `json:"utm_sourcecopy"`
		UtmMediumcopy    string          `json:"utm_mediumcopy"`
		UtmCampaigncopy  string          `json:"utm_campaigncopy"`
		UtmContentcopy   string          `json:"utm_contentcopy"`
		UtmTermcopy      string          `json:"utm_termcopy"`
		Emailcopy        string          `json:"emailcopy"`
		Messengerscopy   string          `json:"messengerscopy"`
		Commentscopy     string          `json:"commentscopy"`
		Phonecopy        string          `json:"phonecopy"`
		Sourcecopy       string          `json:"sourcecopy"`
		ContactLabelcopy string          `json:"contact_labelcopy"`
		KanbanStageId    string          `json:"kanban_stage_id"`
	}

	Deal struct {
		Id              int         `json:"id"`
		CreatedAt       int         `json:"created_at"`
		UpdatedAt       int         `json:"updated_at"`
		Title           string      `json:"title"`
		Owner           int         `json:"owner"`
		CreatedBy       int         `json:"created_by"`
		ChangedBy       interface{} `json:"changed_by"`
		Converted       bool        `json:"converted"`
		AmountOfTheDeal string      `json:"amount_of_the_deal"`
		Contacts        []struct {
			Id               int          `json:"id"`
			CreatedAt        int          `json:"created_at"`
			UpdatedAt        int          `json:"updated_at"`
			Title            string       `json:"title"`
			Owner            int          `json:"owner"`
			CreatedBy        int          `json:"created_by"`
			ChangedBy        interface{}  `json:"changed_by"`
			Converted        bool         `json:"converted"`
			FirstName        string       `json:"first_name"`
			LastName         string       `json:"last_name"`
			Patronymic       string       `json:"patronymic"`
			Companies        string       `json:"companies"`
			Position         string       `json:"position"`
			UtmSource        string       `json:"utm_source"`
			UtmMedium        string       `json:"utm_medium"`
			UtmCampaign      string       `json:"utm_campaign"`
			UtmContent       string       `json:"utm_content"`
			UtmTerm          string       `json:"utm_term"`
			Messengers       *[]Messenger `json:"messengers"`
			Phone            *[]Phone     `json:"phone"`
			Email            *[]Email     `json:"email"`
			Comments         string       `json:"comments"`
			Source           string       `json:"source"`
			ContactLabel     string       `json:"contact_label"`
			Ownercopy        string       `json:"ownercopy"`
			LastNamecopy     string       `json:"last_namecopy"`
			Companiescopy    string       `json:"companiescopy"`
			FirstNamecopy    string       `json:"first_namecopy"`
			Titlecopy        string       `json:"titlecopy"`
			Patronymiccopy   string       `json:"patronymiccopy"`
			Positioncopy     string       `json:"positioncopy"`
			UtmSourcecopy    string       `json:"utm_sourcecopy"`
			UtmMediumcopy    string       `json:"utm_mediumcopy"`
			UtmCampaigncopy  string       `json:"utm_campaigncopy"`
			UtmContentcopy   string       `json:"utm_contentcopy"`
			UtmTermcopy      string       `json:"utm_termcopy"`
			Emailcopy        string       `json:"emailcopy"`
			Messengerscopy   string       `json:"messengerscopy"`
			Commentscopy     string       `json:"commentscopy"`
			Phonecopy        string       `json:"phonecopy"`
			Sourcecopy       string       `json:"sourcecopy"`
			ContactLabelcopy string       `json:"contact_labelcopy"`
			RelatedEntity    bool         `json:"related_entity"`
			EntityType       string       `json:"entity_type"`
			KanbanStageId    string       `json:"kanban_stage_id"`
		} `json:"contacts"`
		Companies []struct {
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
			Phone            []Phone     `json:"phone"`
			Email            []Email     `json:"email"`
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
		} `json:"companies"`
		UtmSource      string          `json:"utm_source"`
		UtmMedium      string          `json:"utm_medium"`
		UtmCampaign    string          `json:"utm_campaign"`
		UtmContent     string          `json:"utm_content"`
		UtmTerm        string          `json:"utm_term"`
		Comments       string          `json:"comments"`
		Source         []GeneralSource `json:"source"`
		DealLabel      []Label         `json:"deal_label"`
		KanbanStatus   string          `json:"kanban_status"`
		KanbanReasonId string          `json:"kanban_reason_id"`
		Ownercopy      string          `json:"ownercopy"`
		KanbanStageId  int             `json:"kanban_stage_id"`
	}

	Lead struct {
		Id             int             `json:"id"`
		CreatedAt      int             `json:"created_at"`
		UpdatedAt      int             `json:"updated_at"`
		Title          string          `json:"title"`
		Owner          int             `json:"owner"`
		CreatedBy      int             `json:"created_by"`
		ChangedBy      string          `json:"changed_by"`
		Converted      bool            `json:"converted"`
		FirstName      string          `json:"first_name"`
		LastName       string          `json:"last_name"`
		Patronymic     string          `json:"patronymic"`
		CompanyName    string          `json:"company_name"`
		Position       string          `json:"position"`
		UtmSource      string          `json:"utm_source"`
		UtmMedium      string          `json:"utm_medium"`
		UtmCampaign    string          `json:"utm_campaign"`
		UtmContent     string          `json:"utm_content"`
		UtmTerm        string          `json:"utm_term"`
		Messengers     []Messenger     `json:"messengers"`
		Phone          []Phone         `json:"phone"`
		Email          *[]Email        `json:"email"`
		Comments       string          `json:"comments"`
		Source         []GeneralSource `json:"source"`
		LeadLabel      []Label         `json:"lead_label"`
		KanbanStatus   string          `json:"kanban_status"`
		KanbanReasonId interface{}     `json:"kanban_reason_id"`
		Ownercopy      string          `json:"ownercopy"`
		Sourcecopy     string          `json:"sourcecopy"`
		LeadLabelcopy  string          `json:"lead_labelcopy"`
		KanbanStageId  int             `json:"kanban_stage_id"`
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
	Messenger struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Link string `json:"link"`
		Sort string `json:"sort"`
	}

	Phone struct {
		Id    string `json:"id"`
		Type  string `json:"type"`
		Value string `json:"value"`
		Main  bool   `json:"main"`
		Sort  string `json:"sort"`
	}

	Email struct {
		Id    string `json:"id"`
		Type  string `json:"type"`
		Value string `json:"value"`
		Main  bool   `json:"main"`
		Sort  string `json:"sort"`
	}

	GeneralSource struct {
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

	Links struct {
		First string      `json:"first"`
		Last  string      `json:"last"`
		Prev  interface{} `json:"prev"`
		Next  interface{} `json:"next"`
	}

	Meta struct {
		CurrentPage int `json:"current_page"`
		From        int `json:"from"`
		LastPage    int `json:"last_page"`
		Links       []struct {
			Url    *string `json:"url"`
			Label  string  `json:"label"`
			Active bool    `json:"active"`
		} `json:"links"`
		Path    string `json:"path"`
		PerPage int    `json:"per_page"`
		To      int    `json:"to"`
		Total   int    `json:"total"`
	}
)
