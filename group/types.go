package group

const (
	VersionGroupUrl   = "groups/v1/"
	TransferGroupsUrl = "transfer"
)

type TransferGroupOutput struct {
	Status         bool `json:"status"`
	CountGroups    int  `json:"countGroups"`
	CountModerator int  `json:"countModerator"`
}
