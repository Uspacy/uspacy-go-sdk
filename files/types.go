package files

const (
	VersionUrl = "files/v1/"
	FilesUrl   = "files"
)

type Files struct {
	Data []struct {
		Id               int         `json:"id"`
		EntityType       string      `json:"entityType"`
		EntityId         interface{} `json:"entityId"`
		UploadId         string      `json:"uploadId"`
		OriginalFilename string      `json:"originalFilename"`
		LastModified     int         `json:"lastModified"`
		Size             int         `json:"size"`
		Url              string      `json:"url"`
	} `json:"data"`
}
