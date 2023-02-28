package api

import (
	"encoding/json"
	"uspacy-go-sdk/newsfeed"
)

// GetNewsFeeds gets all news feeds
func (us *Uspacy) GetNewsFeeds() (newsfeed.SearchNewsfeed, error) {
	var posts newsfeed.SearchNewsfeed

	body, err := us.doGetEmptyHeaders(buildURL(mainHost, newsfeed.VersionNewsfeedUrl, newsfeed.NewsfeedsUrl))
	if err != nil {
		return posts, err
	}
	return posts, json.Unmarshal(body, &posts)
}
