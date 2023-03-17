package group

const (
	VersionUrl  = "groups/v1/"
	TransferUrl = "transfer"
)

type TransferGroupOutput struct {
	Status         bool `json:"status"`
	CountGroups    int  `json:"countGroups"`
	CountModerator int  `json:"countModerator"`
}
