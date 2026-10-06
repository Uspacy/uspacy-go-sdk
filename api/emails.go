package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/emails"
)

// GetMailFolders this method return list of mail folders
func (us *Uspacy) GetMailFolders(ctx context.Context) (folders emails.MailFolders, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(emails.VersionUrl, emails.MailFoldersUrl))
	if err != nil {
		return folders, err
	}
	return folders, json.Unmarshal(body, &folders)
}

// DoMailFolder this method create mail folder and return created mail folder object or error
func (us *Uspacy) DoMailFolder(ctx context.Context, folder emails.MailFolder, opts ...RequestOption) (createdFolder emails.MailFolder, err error) {
	body, _, err := us.doPost(ctx, us.buildURL(emails.VersionUrl, emails.MailFoldersUrl), folder, opts...)
	if err != nil {
		return createdFolder, err
	}
	return createdFolder, json.Unmarshal(body, &createdFolder)
}

// GetMailBoxes this method return list of mail boxes
func (us *Uspacy) GetMailBoxes(ctx context.Context) (boxes emails.MailBoxes, err error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(emails.VersionUrl, emails.MailBoxesUrl))
	if err != nil {
		return boxes, err
	}
	return boxes, json.Unmarshal(body, &boxes)
}

// DoLettersByFolder this method crete letter in folder and return created letter object or error
func (us *Uspacy) DoLettersByFolder(ctx context.Context, folderID string, letter map[string]any, opts ...RequestOption) (createdLetter emails.Letter, code int, err error) {
	body, code, err := us.doPost(ctx, us.buildURL(emails.VersionUrl, fmt.Sprintf(emails.LettersByFolderUrl, folderID)), letter, opts...)
	if err != nil {
		return createdLetter, code, err
	}
	return createdLetter, code, json.Unmarshal(body, &createdLetter)
}

// DeleteLetterById this method delete letter by Id and return answer code and error
func (us *Uspacy) DeleteLetterById(ctx context.Context, letterId int) (code int, err error) {
	code, err = us.doDeleteEmptyHeaders(ctx, us.buildURL(emails.VersionUrl, fmt.Sprintf(emails.LetterById, letterId)), nil)
	if err != nil {
		return code, err
	}
	return code, err
}
