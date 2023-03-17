package newsfeed

import (
	"github.com/Uspacy/uspacy-go-sdk/common"
	"github.com/google/uuid"
)

const (
	VersionUrl = "newsfeed/v1/"
	PostsUrl   = "posts/?page=%d&list=%d&group_id=%d"
)

type (
	GetNewsfeed struct {
		Data  []Post       `json:"data"`
		Links common.Links `json:"links"`
		Meta  common.Meta  `json:"meta"`
	}

	Post struct {
		ID            int         `json:"id"`
		Title         string      `json:"title"`
		Message       string      `json:"message"`
		Date          int         `json:"date"`
		AuthorMood    string      `json:"author_mood"`
		GroupID       int64       `json:"group_id"`
		AuthorID      int         `json:"authorId"`
		Comments      []Comments  `json:"comments"`
		TotalComments int         `json:"total_comments"`
		Files         []Files     `json:"files"`
		Reactions     []Reactions `json:"reactions"`
	}

	Files struct {
		ID               int       `json:"id"`
		EntityType       string    `json:"entityType"`
		EntityID         int       `json:"entityId"`
		UploadID         uuid.UUID `json:"uploadId"`
		OriginalFilename string    `json:"originalFilename"`
		LastModified     int       `json:"lastModified"`
		Size             int       `json:"size"`
		URL              string    `json:"url"`
	}

	Comments struct {
		ID         int         `json:"id"`
		EntityType string      `json:"entityType"`
		EntityID   int         `json:"entityId"`
		AuthorID   int         `json:"authorId"`
		Message    string      `json:"message"`
		Date       int         `json:"date"`
		NextId     int         `json:"nextId"`
		PrevId     int         `json:"prevId"`
		Reactions  []Reactions `json:"reactions"`
	}

	Reactions struct {
		Reaction int64 `json:"reaction"`
		Amount   int64 `json:"amount"`
		EntityId int64 `json:"entityId"`
	}
)
