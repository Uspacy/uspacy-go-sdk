package group

import (
	"github.com/Uspacy/uspacy-go-sdk/common"
)

const (
	VersionUrl  = "groups/v1"
	GroupUrl    = "groups"
	TransferUrl = "transfer"
)

type (
	Groups struct {
		Data  []Group      `json:"data"`
		Links common.Links `json:"links"`
		Meta  Meta         `json:"meta"`
	}
	Group struct {
		ID            int         `json:"id"`
		Name          string      `json:"name"`
		GroupType     string      `json:"groupType"`
		OwnerID       string      `json:"ownerId"`
		Archived      bool        `json:"archived,omitempty"`
		Description   string      `json:"description"`
		GroupTheme    string      `json:"groupTheme,omitempty"`
		ModeratorsIds []string    `json:"moderatorsIds"`
		UsersIds      []string    `json:"usersIds"`
		Logo          interface{} `json:"logo,omitempty"`
	}

	TransferGroupOutput struct {
		Status         bool `json:"status"`
		CountGroups    int  `json:"countGroups"`
		CountModerator int  `json:"countModerator"`
	}
)

type Meta struct {
	CurrentPage int    `json:"currentPage"`
	From        int    `json:"from"`
	LastPage    int    `json:"lastPage"`
	Path        string `json:"path"`
	PerPage     int    `json:"perPage"`
	To          int    `json:"to"`
	Total       int    `json:"total"`
}
