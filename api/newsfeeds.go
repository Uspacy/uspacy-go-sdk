package api

import (
	"encoding/json"
	"fmt"
	"uspacy-go-sdk/newsfeed"
)

// GetNewsfeeds gets all newsfeeds
func (us *Uspacy) GetNewsfeeds(page, list, groupId int) (newsfeed.GetNewsfeed, error) {
	var posts newsfeed.GetNewsfeed

	body, err := us.doGetEmptyHeaders(buildURL(mainHost, newsfeed.VersionNewsfeedUrl, fmt.Sprintf(newsfeed.PostsUrl, page, list, groupId)))
	if err != nil {
		return posts, err
	}
	return posts, json.Unmarshal(body, &posts)
}
