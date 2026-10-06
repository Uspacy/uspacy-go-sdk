package api

import (
	"context"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/v2/comments"
)

// CreateComment returns created comment
func (us *Uspacy) CreateComment(ctx context.Context, commentsData comments.Comment, opts ...RequestOption) (comments.Comment, error) {
	body, _, err := us.doPost(ctx, us.buildURL(comments.VersionUrl, fmt.Sprintf(comments.CommentsUrl, "")), commentsData, opts...)
	return decodeJSON[comments.Comment](body, err)
}
