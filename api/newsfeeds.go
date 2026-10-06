package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/v2/newsfeed"
)

// CreateNewsfeedPost returns created post
func (us *Uspacy) CreateNewsfeedPost(ctx context.Context, postData url.Values, opts ...RequestOption) (err error) {
	_, err = us.doPostEncodedForm(ctx, us.buildURL(newsfeed.VersionUrl, newsfeed.DoPostUrl), postData, opts...)
	if err != nil {
		return err
	}
	return nil
}

// GetNewsfeeds gets all newsfeeds
func (us *Uspacy) GetNewsfeeds(ctx context.Context, page, list, groupId int, opts ...RequestOption) (newsfeed.GetNewsfeed, error) {
	body, err := us.doGet(ctx, us.buildURL(newsfeed.VersionUrl, fmt.Sprintf(newsfeed.GetPostsUrl, page, list, groupId)), opts...)
	return decodeJSON[newsfeed.GetNewsfeed](body, err)
}
