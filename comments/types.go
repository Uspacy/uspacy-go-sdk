package comments

const (
	VersionUrl  = "comments/v1"
	CommentsUrl = "comments/%v"
)

type Comment struct {
	ID         int    `json:"id"`
	AuthorID   int    `json:"authorId"`
	EntityType string `json:"entityType"`
	EntityID   int    `json:"entityId"`
	Message    string `json:"message"`
	Date       int    `json:"date"`
	FileIds    []int  `json:"file_ids"`
}
