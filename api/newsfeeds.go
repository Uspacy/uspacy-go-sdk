package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/v2/newsfeed"
)

// CreateNewsfeedPost returns created post
func (us *Uspacy) CreateNewsfeedPost(ctx context.Context, postData url.Values) (err error) {
	_, err = us.doPostEncodedForm(ctx, us.buildURL(newsfeed.VersionUrl, newsfeed.DoPostUrl), postData)
	if err != nil {
		return err
	}
	return nil
}

// GetNewsfeeds gets all newsfeeds
func (us *Uspacy) GetNewsfeeds(ctx context.Context, page, list, groupId int) (newsfeed.GetNewsfeed, error) {
	body, err := us.doGet(ctx, us.buildURL(newsfeed.VersionUrl, fmt.Sprintf(newsfeed.GetPostsUrl, page, list, groupId)))
	return decodeJSON[newsfeed.GetNewsfeed](body, err)
}
