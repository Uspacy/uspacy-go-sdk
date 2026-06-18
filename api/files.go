package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

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
	code, err = us.doDeleteEmptyHeaders(us.buildURL(files.VersionUrl, fmt.Sprintf("%s/%d", files.FilesUrl, fileId)), nil)
	if err != nil {
		return code, err
	}
	return code, err
}

// DeleteFilesByEntityId this method delete all files by EntityId and return answer code and error
func (us *Uspacy) DeleteFilesByEntityId(entityType string, entityId int64) (code int, err error) {
	code, err = us.doDeleteEmptyHeaders(us.buildURL(files.VersionUrl, fmt.Sprintf("%s?%s&%d", files.FilesUrl, entityType, entityId)), nil)
	if err != nil {
		return code, err
	}
	return code, err
}

// GetFileById this method get file by Id and return file object and error
func (us *Uspacy) GetFileById(fileId int) (file files.File, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(files.VersionUrl, files.FilesUrl, strconv.Itoa(fileId)))
	if err != nil {
		return file, err
	}
	return file, json.Unmarshal(body, &file)
}

// UpdateFile this method updates file entityId by fileId and returns updated file object and error
func (us *Uspacy) UpdateFile(fileId int, entityId int64) (file files.File, err error) {
	body := map[string]int64{
		"entityId": entityId,
	}
	response, err := us.doPatchEmptyHeaders(us.buildURL(files.VersionUrl, fmt.Sprintf("%s/%d", files.FilesUrl, fileId)), body)
	if err != nil {
		return file, err
	}
	return file, json.Unmarshal(response, &file)
}
