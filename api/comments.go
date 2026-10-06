package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/comments"
)

// CreateComment returns created comment
func (us *Uspacy) CreateComment(ctx context.Context, commentsData comments.Comment, opts ...RequestOption) (comment comments.Comment, err error) {
	body, _, err := us.doPost(ctx, us.buildURL(comments.VersionUrl, fmt.Sprintf(comments.CommentsUrl, "")), commentsData, opts...)
	if err != nil {
		return comment, err
	}
	return comment, json.Unmarshal(body, &comment)
}
