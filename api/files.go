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

// DeleteFileById deletes a file and returns the HTTP status code.
func (us *Uspacy) DeleteFileById(ctx context.Context, fileId int, opts ...RequestOption) (code int, err error) {
	code, err = us.doDelete(ctx, us.buildURL(files.VersionUrl, fmt.Sprintf("%s/%d", files.FilesUrl, fileId)), nil, opts...)
	if err != nil {
		return code, err
	}
	return code, err
}

// DeleteFilesByEntityId deletes all files attached to an entity and returns the HTTP status code.
func (us *Uspacy) DeleteFilesByEntityId(ctx context.Context, entityType string, entityId int64, opts ...RequestOption) (code int, err error) {
	params := url.Values{}
	params.Set("entityType", entityType)
	params.Set("entityId", strconv.FormatInt(entityId, 10))
	code, err = us.doDelete(ctx, us.buildURL(files.VersionUrl, files.FilesUrl)+"?"+params.Encode(), nil, opts...)
	if err != nil {
		return code, err
	}
	return code, err
}

// GetFileById returns a file by its ID.
func (us *Uspacy) GetFileById(ctx context.Context, fileId int, opts ...RequestOption) (files.File, error) {
	body, err := us.doGet(ctx, us.buildURL(files.VersionUrl, files.FilesUrl, strconv.Itoa(fileId)), opts...)
	return decodeJSON[files.File](body, err)
}

// UpdateFile attaches a file to another entity and returns the updated file.
func (us *Uspacy) UpdateFile(ctx context.Context, fileId int, entityId int64, opts ...RequestOption) (files.File, error) {
	body := map[string]int64{
		"entityId": entityId,
	}
	response, err := us.doPatch(ctx, us.buildURL(files.VersionUrl, fmt.Sprintf("%s/%d", files.FilesUrl, fileId)), body, opts...)
	return decodeJSON[files.File](response, err)
}
