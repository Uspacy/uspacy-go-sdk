package newsfeed

const (
	VersionNewsfeedUrl = "newsfeed/v1/"
	NewsfeedsUrl       = "posts"
)

type (
	PostsOutput struct {
		ID         int    `json:"id"`
		Title      string `json:"title"`
		Message    string `json:"message"`
		AuthorID   int    `json:"authorId"`
		Date       int    `json:"date"`
		CreatedAt  string `json:"createdAt"`
		UpdatedAt  string `json:"updatedAt"`
		AuthorMood string `json:"authorMood"`
	}

	SearchNewsfeed struct {
		Data []PostsOutput `json:"data"`
	}
)
