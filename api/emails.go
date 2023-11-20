package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/emails"
)

// GetMailFolders this method return list of mail folders
func (us *Uspacy) GetMailFolders() (folders emails.MailFolders, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(emails.VersionUrl, emails.MailFoldersUrl))
	if err != nil {
		return folders, err
	}
	return folders, json.Unmarshal(body, &folders)
}

// DoMailFolder this method create mail folder and return created mail folder object or error
func (us *Uspacy) DoMailFolder(folder emails.MailFolder) (createdFolder emails.MailFolder, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(emails.VersionUrl, emails.MailFoldersUrl), folder)
	if err != nil {
		return createdFolder, err
	}
	return createdFolder, json.Unmarshal(body, &createdFolder)
}

// GetMailBoxes this method return list of mail boxes
func (us *Uspacy) GetMailBoxes() (boxes emails.MailBoxes, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(emails.VersionUrl, emails.MailBoxesUrl))
	if err != nil {
		return boxes, err
	}
	return boxes, json.Unmarshal(body, &boxes)
}

// DoLettersByFolder this method crete letter in folder and return created letter object or error
func (us *Uspacy) DoLettersByFolder(folderID string, letter emails.Letter) (createdLetter emails.Letter, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(emails.VersionUrl, fmt.Sprintf(emails.LettersByFolderUrl, folderID)), letter)
	if err != nil {
		return createdLetter, err
	}
	return createdLetter, json.Unmarshal(body, &createdLetter)
}
