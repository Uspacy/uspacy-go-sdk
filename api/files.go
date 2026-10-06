package api

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/Uspacy/uspacy-go-sdk/v2/files"
)

// CreateFile creates a file attached to the given entity.
func (us *Uspacy) CreateFile(ctx context.Context, entityType, entityId string, filesMap map[string]io.ReadCloser, opts ...RequestOption) (files.Files, error) {
	textParams := map[string]string{
		"entityType": entityType,
		"entityId":   entityId,
	}
	body, err := us.doPostFormData(ctx, us.buildURL(files.VersionUrl, files.FilesUrl), textParams, filesMap, opts...)
	return decodeJSON[files.Files](body, err)
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
	params := url.Values{}
	params.Set("entityType", entityType)
	params.Set("entityId", strconv.FormatInt(entityId, 10))
	code, err = us.doDeleteEmptyHeaders(ctx, us.buildURL(files.VersionUrl, files.FilesUrl)+"?"+params.Encode(), nil)
	if err != nil {
		return code, err
	}
	return code, err
}

// GetFileById this method get file by Id and return file object and error
func (us *Uspacy) GetFileById(ctx context.Context, fileId int) (files.File, error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(files.VersionUrl, files.FilesUrl, strconv.Itoa(fileId)))
	return decodeJSON[files.File](body, err)
}

// UpdateFile this method updates file entityId by fileId and returns updated file object and error
func (us *Uspacy) UpdateFile(ctx context.Context, fileId int, entityId int64) (files.File, error) {
	body := map[string]int64{
		"entityId": entityId,
	}
	response, err := us.doPatchEmptyHeaders(ctx, us.buildURL(files.VersionUrl, fmt.Sprintf("%s/%d", files.FilesUrl, fileId)), body)
	return decodeJSON[files.File](response, err)
}
