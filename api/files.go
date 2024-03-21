package api

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Uspacy/uspacy-go-sdk/files"
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

// DeleteFileById this method delete file by Id and return answer code and error
func (us *Uspacy) DeleteFileById(fileId int) (code int, err error) {
	code, err = us.doDeleteEmptyHeaders(us.buildURL(files.VersionUrl, fmt.Sprintf("%s/%d", files.FilesUrl, fileId)))
	if err != nil {
		return code, err
	}
	return code, err
}

// DeleteFilesByEntityId this method delete all files by EntityId and return answer code and error
func (us *Uspacy) DeleteFilesByEntityId(entityType string, entityId int64) (code int, err error) {
	code, err = us.doDeleteEmptyHeaders(us.buildURL(files.VersionUrl, fmt.Sprintf("%s?%s&%d", files.FilesUrl, entityType, entityId)))
	if err != nil {
		return code, err
	}
	return code, err
}
