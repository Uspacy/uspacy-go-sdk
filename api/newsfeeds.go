package api

import (
	"encoding/json"
	"fmt"

	"github.com/Uspacy/uspacy-go-sdk/newsfeed"
)

// GetNewsfeeds gets all newsfeeds
func (us *Uspacy) GetNewsfeeds(page, list, groupId int) (posts newsfeed.GetNewsfeed, err error) {
	body, err := us.doGetEmptyHeaders(buildURL(us.mainHost, newsfeed.VersionUrl, fmt.Sprintf(newsfeed.PostsUrl, page, list, groupId)))
	if err != nil {
		return posts, err
	}
	return posts, json.Unmarshal(body, &posts)
}
