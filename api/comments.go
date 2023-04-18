package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/comments"
)

// CreateComment returns created comment
func (us *Uspacy) CreateComment(commentsData comments.Comment) (comment comments.Comment, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(comments.VersionUrl, fmt.Sprintf(comments.CommentsUrl, "")), commentsData)
	if err != nil {
		return comment, err
	}
	return comment, json.Unmarshal(body, &comment)
}
