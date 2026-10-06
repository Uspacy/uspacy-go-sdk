package api

import (
	"context"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/emails"
)

// GetMailFolders this method return list of mail folders
func (us *Uspacy) GetMailFolders(ctx context.Context) (emails.MailFolders, error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(emails.VersionUrl, emails.MailFoldersUrl))
	return decodeJSON[emails.MailFolders](body, err)
}

// DoMailFolder this method create mail folder and return created mail folder object or error
func (us *Uspacy) DoMailFolder(ctx context.Context, folder emails.MailFolder, opts ...RequestOption) (emails.MailFolder, error) {
	body, _, err := us.doPost(ctx, us.buildURL(emails.VersionUrl, emails.MailFoldersUrl), folder, opts...)
	return decodeJSON[emails.MailFolder](body, err)
}

// GetMailBoxes this method return list of mail boxes
func (us *Uspacy) GetMailBoxes(ctx context.Context) (emails.MailBoxes, error) {
	body, err := us.doGetEmptyHeaders(ctx, us.buildURL(emails.VersionUrl, emails.MailBoxesUrl))
	return decodeJSON[emails.MailBoxes](body, err)
}

// DoLettersByFolder this method crete letter in folder and return created letter object or error
func (us *Uspacy) DoLettersByFolder(ctx context.Context, folderID string, letter map[string]any, opts ...RequestOption) (emails.Letter, int, error) {
	body, code, err := us.doPost(ctx, us.buildURL(emails.VersionUrl, fmt.Sprintf(emails.LettersByFolderUrl, folderID)), letter, opts...)
	v, err := decodeJSON[emails.Letter](body, err)
	return v, code, err
}

// DeleteLetterById this method delete letter by Id and return answer code and error
func (us *Uspacy) DeleteLetterById(ctx context.Context, letterId int) (code int, err error) {
	code, err = us.doDeleteEmptyHeaders(ctx, us.buildURL(emails.VersionUrl, fmt.Sprintf(emails.LetterById, letterId)), nil)
	if err != nil {
		return code, err
	}
	return code, err
}
