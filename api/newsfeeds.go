package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/newsfeed"
)

// CreateNewsfeedPost returns created post
func (us *Uspacy) CreateNewsfeedPost(postData newsfeed.Post) (post newsfeed.Post, err error) {
	body, err := us.doPostEmptyHeaders(us.buildURL(newsfeed.VersionUrl, ""), postData)
	if err != nil {
		return post, err
	}
	return post, json.Unmarshal(body, &post)
}

// GetNewsfeeds gets all newsfeeds
func (us *Uspacy) GetNewsfeeds(page, list, groupId int) (posts newsfeed.GetNewsfeed, err error) {
	body, err := us.doGetEmptyHeaders(us.buildURL(newsfeed.VersionUrl, fmt.Sprintf(newsfeed.PostsUrl, page, list, groupId)))
	if err != nil {
		return posts, err
	}
	return posts, json.Unmarshal(body, &posts)
}
