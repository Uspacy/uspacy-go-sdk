package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/comments"
)

// CreateComment returns created comment
func (us *Uspacy) CreateComment(commentsData comments.Comment, headers ...map[string]string) (comment comments.Comment, err error) {
	body, _, err := us.doPost(us.buildURL(comments.VersionUrl, fmt.Sprintf(comments.CommentsUrl, "")), commentsData, headers...)
	if err != nil {
		return comment, err
	}
	return comment, json.Unmarshal(body, &comment)
}
