package common

type (
	ContactData struct {
		Id    string `json:"id"`
		Type  string `json:"type"`
		Value string `json:"value"`
		Main  bool   `json:"main"`
		Sort  string `json:"sort"`
	}

	Meta struct {
		CurrentPage int `json:"current_page"`
		From        int `json:"from"`
		LastPage    int `json:"last_page"`
		PerPage     int `json:"per_page"`
		To          int `json:"to"`
		Total       int `json:"total"`
	}

	Links struct {
		First string      `json:"first"`
		Last  string      `json:"last"`
		Prev  interface{} `json:"prev"`
		Next  interface{} `json:"next"`
	}

	SocialMedia struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Link string `json:"link"`
	}
)
