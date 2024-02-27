package emails

import "time"

const (
	VersionUrl = "email/v1"
)

const (
	MailBoxesUrl       = "emails"
	MailFoldersUrl     = "folders"
	LettersByFolderUrl = "letters/by_folder/%s"
	LetterById         = "letters/%d"
)

type (
	MailFolders struct {
		Data []MailFolder `json:"data"`
	}

	MailFolder struct {
		ID            int    `json:"id"`
		EmailID       int    `json:"email_id"`
		FolderName    string `json:"folder_name"`
		Path          string `json:"path"`
		Delimitter    string `json:"delimitter"`
		MessageCount  int    `json:"message_count"`
		HighestModSeq int    `json:"highest_mod_seq"`
		IsTrash       bool   `json:"is_trash"`
		IsSpam        bool   `json:"is_spam"`
		IsDraft       bool   `json:"is_draft"`
		IsJunk        bool   `json:"is_junk"`
		IsSent        bool   `json:"is_sent"`
		HasChildren   bool   `json:"has_children"`
		Pivot         Pivot  `json:"pivot"`
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
		ImapHost      string    `json:"imap_host"`
		ImapPort      string    `json:"imap_port"`
		Email         string    `json:"email"`
		Password      string    `json:"password"`
		Name          string    `json:"name"`
		SenderName    string    `json:"sender_name"`
		AccessLevel   string    `json:"access_level"`
		LastMessageID string    `json:"last_message_id"`
		Tariff        int       `json:"tariff"`
		HasFile       bool      `json:"has_file"`
		LastSyncedAt  int       `json:"last_synced_at"`
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
		MessageID      string             `json:"message_id"`
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
