package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/Uspacy/uspacy-go-sdk/v2/files"
)

// CreateFile creates a file attached to the given entity.
func (us *Uspacy) CreateFile(ctx context.Context, entityType, entityId string, filesMap map[string]io.ReadCloser) (file files.Files, err error) {
	textParams := map[string]string{
		"entityType": entityType,
		"entityId":   entityId,
	}
	body, err := us.doPostFormData(ctx, us.buildURL(files.VersionUrl, files.FilesUrl), textParams, filesMap)
	if err != nil {
		return file, err
	}
	return file, json.Unmarshal(body, &file)

}

// DeleteFileById this method delete file by Id and return answer code and error
func (us *Uspacy) DeleteFileById(ctx context.Context, fileId int) (code int, err error) {
	code, err = us.doDeleteEmptyHeaders(ctx, us.buildURL(files.VersionUrl, fmt.Sprintf("%s/%d", files.FilesUrl, fileId)), nil)
	if err != nil {
		return code, err
	}
	return code, err
}

// DeleteFilesByEntityId this method delete all files by EntityId and return answer code and error
func (us *Uspacy) DeleteFilesByEntityId(ctx context.Context, entityType string, entityId int64) (code int, err error) {
	code, err = us.doDeleteEmptyHeaders(ctx, us.buildURL(files.VersionUrl, fmt.Sprintf("%s?%s&%d", files.FilesUrl, entityType, entityId)), nil)
	if err != nil {
		return code, err
	}
	return code, err
}

// GetFileById this method get file by Id and return file object and error
func (us *Uspacy) GetFileById(ctx context.Context, fileId int) (file files.File, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(files.VersionUrl, files.FilesUrl, strconv.Itoa(fileId)))
	if err != nil {
		return file, err
	}
	return file, json.Unmarshal(body, &file)
}

// UpdateFile this method updates file entityId by fileId and returns updated file object and error
func (us *Uspacy) UpdateFile(ctx context.Context, fileId int, entityId int64) (file files.File, err error) {
	body := map[string]int64{
		"entityId": entityId,
	}
	response, err := us.doPatchEmptyHeaders(ctx, us.buildURL(files.VersionUrl, fmt.Sprintf("%s/%d", files.FilesUrl, fileId)), body)
	if err != nil {
		return file, err
	}
	return file, json.Unmarshal(response, &file)
}
