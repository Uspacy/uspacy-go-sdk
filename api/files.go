package api

import (
	"encoding/json"
	"github.com/Uspacy/uspacy-go-sdk/files"
	"io"
)

func (us *Uspacy) CreateFile(entityType, entityId string, filesMap map[string]io.ReadCloser) (file files.Files, err error) {
	textParams := map[string]string{
		"entityType": entityType,
		"entityId":   entityId,
	}
	body, err := us.doPostFormData(us.buildURL(files.VersionUrl, files.FilesUrl), textParams, filesMap)
	if err != nil {
		return file, err
	}
	return file, json.Unmarshal(body, &file)

}
