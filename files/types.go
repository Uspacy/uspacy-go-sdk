package files

const (
	VersionUrl = "files/v1"
	FilesUrl   = "files"
)

type Files struct {
	Data []File `json:"data"`
}

type File struct {
	Id               int    `json:"id"`
	EntityType       string `json:"entityType"`
	EntityId         any    `json:"entityId"`
	UploadId         string `json:"uploadId"`
	CreatorId        int    `json:"creatorId"`
	OriginalFilename string `json:"originalFilename"`
	LastModified     int    `json:"lastModified"`
	Size             int    `json:"size"`
	Url              string `json:"url"`
	Width            any    `json:"width"`
	Height           any    `json:"height"`
}
