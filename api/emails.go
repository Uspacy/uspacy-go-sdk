package api

import (
	"context"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/emails"
)

// GetMailFolders returns the mail folders.
func (us *Uspacy) GetMailFolders(ctx context.Context, opts ...RequestOption) (emails.MailFolders, error) {
	body, err := us.doGet(ctx, us.buildURL(emails.VersionUrl, emails.MailFoldersUrl), opts...)
	return decodeJSON[emails.MailFolders](body, err)
}

// DoMailFolder creates a mail folder and returns it.
func (us *Uspacy) DoMailFolder(ctx context.Context, folder emails.MailFolder, opts ...RequestOption) (emails.MailFolder, error) {
	body, _, err := us.doPost(ctx, us.buildURL(emails.VersionUrl, emails.MailFoldersUrl), folder, opts...)
	return decodeJSON[emails.MailFolder](body, err)
}

// GetMailBoxes returns the mailboxes.
func (us *Uspacy) GetMailBoxes(ctx context.Context, opts ...RequestOption) (emails.MailBoxes, error) {
	body, err := us.doGet(ctx, us.buildURL(emails.VersionUrl, emails.MailBoxesUrl), opts...)
	return decodeJSON[emails.MailBoxes](body, err)
}

// DoLettersByFolder creates a letter in a folder and returns it with the HTTP status code.
func (us *Uspacy) DoLettersByFolder(ctx context.Context, folderID string, letter map[string]any, opts ...RequestOption) (emails.Letter, int, error) {
	body, code, err := us.doPost(ctx, us.buildURL(emails.VersionUrl, fmt.Sprintf(emails.LettersByFolderUrl, folderID)), letter, opts...)
	v, err := decodeJSON[emails.Letter](body, err)
	return v, code, err
}

// DeleteLetterById deletes a letter and returns the HTTP status code.
func (us *Uspacy) DeleteLetterById(ctx context.Context, letterId int, opts ...RequestOption) (code int, err error) {
	code, err = us.doDelete(ctx, us.buildURL(emails.VersionUrl, fmt.Sprintf(emails.LetterById, letterId)), nil, opts...)
	if err != nil {
		return code, err
	}
	return code, err
}
