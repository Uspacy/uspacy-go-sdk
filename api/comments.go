package api

import (
	"context"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/comments"
)

// CreateComment creates a comment and returns it.
func (us *Uspacy) CreateComment(ctx context.Context, commentsData comments.Comment, opts ...RequestOption) (comments.Comment, error) {
	body, _, err := us.doPost(ctx, us.buildURL(comments.VersionUrl, fmt.Sprintf(comments.CommentsUrl, "")), commentsData, opts...)
	return decodeJSON[comments.Comment](body, err)
}
