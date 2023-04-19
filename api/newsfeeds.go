package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Uspacy/uspacy-go-sdk/newsfeed"
)

// CreateNewsfeedPost returns created post
func (us *Uspacy) CreateNewsfeedPost(postData url.Values) (err error) {
	_, err = us.doPostFormData(us.buildURL(newsfeed.VersionUrl, newsfeed.DoPostUrl), postData)
	if err != nil {
		return err
	}
	return nil
}

// GetNewsfeeds gets all newsfeeds
func (us *Uspacy) GetNewsfeeds(page, list, groupId int) (posts newsfeed.GetNewsfeed, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(newsfeed.VersionUrl, fmt.Sprintf(newsfeed.GetPostsUrl, page, list, groupId)))
	if err != nil {
		return posts, err
	}
	return posts, json.Unmarshal(body, &posts)
}
