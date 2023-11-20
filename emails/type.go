package emails

import "time"

const (
	VersionUrl = "email/v1"
)

const (
	MailBoxesUrl       = "/emails"
	MailFoldersUrl     = "/folders"
	LettersByFolderUrl = "/letters/by-folder/%s"
)

type (
	MailFolders struct {
		Data []MailFolder `json:"data"`
	}

	MailFolder struct {
		ID           int    `json:"id"`
		EmailID      int    `json:"email_id"`
		FolderName   string `json:"folder_name"`
		Path         string `json:"path"`
		Delimitter   string `json:"delimitter"`
		MessageCount int    `json:"message_count"`
		IsTrash      int    `json:"is_trash"`
		IsSpam       int    `json:"is_spam"`
		IsDraft      int    `json:"is_draft"`
		IsJunk       int    `json:"is_junk"`
		IsSent       int    `json:"is_sent"`
		HasChildren  int    `json:"has_children"`
		Pivot        Pivot  `json:"pivot"`
	}

	Pivot struct {
		LetterID      int `json:"letter_id"`
		EmailFolderID int `json:"email_folder_id"`
	}
)

type (
	MailBoxes struct {
		Data []MailBox `json:"data"`
	}
	MailBox struct {
		ID            int       `json:"id"`
		PortalName    string    `json:"portal_name"`
		AddedBy       int       `json:"added_by"`
		ImapHost      any       `json:"imap_host"`
		ImapPort      any       `json:"imap_port"`
		Email         string    `json:"email"`
		Password      string    `json:"password"`
		Name          string    `json:"name"`
		SenderName    string    `json:"sender_name"`
		AccessLevel   string    `json:"access_level"`
		LastMessageID any       `json:"last_message_id"`
		Tariff        int       `json:"tariff"`
		HasFile       any       `json:"has_file"`
		LastSyncedAt  any       `json:"last_synced_at"`
		Status        string    `json:"status"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}
)

type (
	Letter struct {
		ID             int                `json:"id"`
		Subject        string             `json:"subject"`
		Body           string             `json:"body"`
		BodyHTML       string             `json:"body_html"`
		HasAttachments bool               `json:"has_attachments"`
		MassageID      string             `json:"massage_id"`
		Date           int64              `json:"date"`
		IsRead         bool               `json:"is_read"`
		Attachments    []LetterAttachment `json:"attachments"`
		Contacts       []LetterContact    `json:"contacts"`
		Folders        []MailFolder       `json:"folders"`
	}
	LetterAttachment struct {
		FileID   string `json:"file_id"`
		FileURL  string `json:"file_url"`
		FileName string `json:"file_name"`
	}
	LetterContact struct {
		Email       string `json:"email"`
		Name        string `json:"name"`
		ContactType string `json:"contact_type"`
	}
)
