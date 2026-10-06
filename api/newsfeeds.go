package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/v2/newsfeed"
)

// CreateNewsfeedPost creates a newsfeed post.
func (us *Uspacy) CreateNewsfeedPost(ctx context.Context, postData url.Values, opts ...RequestOption) (err error) {
	_, err = us.doPostEncodedForm(ctx, us.buildURL(newsfeed.VersionUrl, newsfeed.DoPostUrl), postData, opts...)
	if err != nil {
		return err
	}
	return nil
}

// GetNewsfeeds returns a page of newsfeed posts of a group.
func (us *Uspacy) GetNewsfeeds(ctx context.Context, page, list, groupId int, opts ...RequestOption) (newsfeed.GetNewsfeed, error) {
	body, err := us.doGet(ctx, us.buildURL(newsfeed.VersionUrl, fmt.Sprintf(newsfeed.GetPostsUrl, page, list, groupId)), opts...)
	return decodeJSON[newsfeed.GetNewsfeed](body, err)
}
